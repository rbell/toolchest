/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package workqueue

import (
	"container/heap"
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
)

type workOption func(item *workItem)

// Queue allow work to be queued up and worked on in a set number of go routines
type Queue struct {
	workerCount      int
	queueLength      *atomic.Int32
	workChan         chan *workItem
	workQueue        *workHeap
	errChan          chan error
	errSubScriberMux *sync.Mutex
	errorSubscribers []chan error
	stopped          atomic.Bool
	breaked          bool
	workItems        *sync.Map
	queueContext     context.Context
	queueCancel      context.CancelFunc
}

// NewQueue returns a reference to an initialized Queue
func NewQueue(options ...WorkQueueOption) *Queue {
	wq := &Queue{
		workerCount:      runtime.NumCPU(),
		queueLength:      &atomic.Int32{},
		errChan:          make(chan error),
		errSubScriberMux: &sync.Mutex{},
		errorSubscribers: []chan error{},
		stopped:          atomic.Bool{},
		breaked:          false,
		workItems:        &sync.Map{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	wq.queueContext = ctx
	wq.queueCancel = cancel
	wq.queueLength.Store(int32(runtime.NumCPU() * 2))
	wq.stopped.Store(false)
	for _, o := range options {
		o(wq)
	}

	wq.workChan = make(chan *workItem)
	wq.workQueue = newWorkHeap(int(wq.queueLength.Load()))

	go wq.start()

	return wq
}

// Enqueue queues work to do on the workChan to be processed
func (w *Queue) Enqueue(workToDo Work, options ...workOption) uuid.UUID {
	wi := &workItem{
		QueuedWork: &QueuedWork{
			id:       uuid.New(),
			priority: 1,
			state:    &atomic.Int32{},
		},

		workToDo: workToDo,
	}
	for _, option := range options {
		option(wi)
	}

	w.workItems.Store(wi.id, wi)

	if !w.stopped.Load() {
		w.workChan <- wi
	}

	return wi.id
}

// Dequeue removes from the queue the work item with the specified id.  If the work item is in process, then an error is returned.
func (w *Queue) Dequeue(id uuid.UUID) error {
	if i, ok := w.workItems.Load(id); ok {
		wi := i.(*workItem)
		if wi.state.Load() == int32(IN_QUEUE) {
			w.workQueue.Remove(wi.position)
			w.workQueue.AdjustPriorities()
			w.workItems.Delete(wi.id)
		} else if wi.state.Load() == int32(IN_PROGRESS) {
			return fmt.Errorf("cannot delete work item %v because it is in process", id.String())
		}
	}
	return nil
}

// SetPriority changes the priority of the queued work item with the uuid.
func (w *Queue) SetPriority(id uuid.UUID, priority int) error {
	if i, ok := w.workItems.Load(id); ok {
		wi := i.(*workItem)
		if wi.state.Load() == int32(IN_QUEUE) {
			wi.priority = priority
			w.workQueue.AdjustPriorities()
		} else if wi.state.Load() == int32(IN_PROGRESS) {
			return fmt.Errorf("cannot adjust prioroty on work item %v because it is in process", id.String())
		}
	}
	return nil
}

// WorkItems returns the current state of all queued work items
func (w *Queue) WorkItems() []*QueuedWork {
	result := []*QueuedWork{}
	w.workItems.Range(func(key, value any) bool {
		result = append(result, value.(*workItem).QueuedWork)
		return true
	})
	return result
}

// Errors allows monitoring errors that occur on work submitted to queue
func (w *Queue) Errors() chan error {
	w.errSubScriberMux.Lock()
	defer w.errSubScriberMux.Unlock()

	ch := make(chan error)
	w.errorSubscribers = append(w.errorSubscribers, ch)
	return ch
}

// Stop stops the queue from accepting work
func (w *Queue) Stop() {
	w.stopped.Store(true)
	w.queueCancel()
}

// Break stops the queue form accepting any work and any work in queue is skipped
func (w *Queue) Break() {
	w.breaked = true
	w.Stop()
}

// ResizeQueueLength adjusts the size of the queue
func (w *Queue) ResizeQueueLength(length int) {
	w.queueLength.Store(int32(length))
}

func (w *Queue) start() {
	defer func() {
		close(w.errChan)
		close(w.workChan)
		w.queueCancel()
	}()

	heap.Init(w.workQueue)

	// monitor for errors on a go routine
	go func() {
		stop := false
		for {
			select {
			case e := <-w.errChan:
				for _, sub := range w.errorSubscribers {
					sub <- e
				}
			case <-w.queueContext.Done():
				stop = true
			}
			if stop {
				break
			}
		}
	}()

	workerSemaphore := make(chan bool, w.workerCount)
	workerCh := make(chan *workItem, w.workerCount)
	defer close(workerCh)
	defer close(workerSemaphore)
	for i := 0; i < w.workerCount; i++ {
		go w.doWork(workerCh, workerSemaphore)
	}

	// Process work
outsideFor:
	for {
	insideFor:
		select {
		case work := <-w.workChan:
			if work != nil {
				// If queue is empty try to send directly to workers via workerCh
				if w.workQueue.Len() == 0 {
					select {
					case workerCh <- work:
						break insideFor
					default:
					}
				}
				// Workers are busy and placing the workToDo on the channel failed - queue it on prioritized queue
				if w.workQueue.Len() < int(w.queueLength.Load()) {
					heap.Push(w.workQueue, work)
					//w.workQueue.AdjustPriorities()
				} else {
					// queue is full, block and wait for worker to finish a task then add work to queue
					fmt.Println("Waiting for free worker")
					<-workerSemaphore
					w.workQueue.AdjustPriorities()
					wtemp := heap.Pop(w.workQueue).(*workItem)
					workerCh <- wtemp
					w.workQueue.Push(work)
				}
				fmt.Printf("Queue Length %v\n", w.workQueue.Len())
			}
		case <-workerSemaphore:
			// worker done, pop and start next work (if anything in queue)
			if w.workQueue.Len() > 0 {
				w.workQueue.AdjustPriorities()
				wtemp := heap.Pop(w.workQueue).(*workItem)
				workerCh <- wtemp
			}
		case <-w.queueContext.Done():
			break outsideFor
		}
	}

	// Finish any work left on queue
	for _, work := range w.workQueue.items {
		if !w.breaked {
			workerCh <- work
		}
	}
}

func (w *Queue) doWork(workCh chan *workItem, semaphore chan bool) {
	for wi := range workCh {
		wi.state.Store(int32(IN_PROGRESS))
		err := wi.workToDo()
		if err != nil {
			w.errChan <- err
		}
		w.workItems.Delete(wi.id)
		semaphore <- true
	}
}
