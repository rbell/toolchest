/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package workqueue

import (
	"container/heap"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
)

type workState int32

// work states
const (
	IN_QUEUE workState = iota
	IN_PROCESS
)

type Work func() error

type WorkQueueOption func(*Queue)

type QueuedWork struct {
	id       uuid.UUID
	name     string
	priority int
	index    int
	state    *atomic.Int32
}

type workItem struct {
	*QueuedWork
	workToDo       Work
	adjustPriority func() int
}

type workOption func(item *workItem)

// Queue allow work to be queued up and worked on in a set number of go routines
type Queue struct {
	workerCount      int
	queueLength      *atomic.Int32
	workChan         chan *workItem
	workQueue        *workHeap
	errChan          chan error
	stopSignal       chan struct{}
	errSubScriberMux *sync.Mutex
	errorSubscribers []chan error
	stopped          atomic.Bool
	breaked          bool
	workItems        *sync.Map
}

// NewQueue returns a reference to an initialized Queue
func NewQueue(options ...WorkQueueOption) *Queue {
	wq := &Queue{
		workerCount:      runtime.NumCPU(),
		queueLength:      &atomic.Int32{},
		stopSignal:       make(chan struct{}),
		errChan:          make(chan error),
		errSubScriberMux: &sync.Mutex{},
		errorSubscribers: []chan error{},
		stopped:          atomic.Bool{},
		breaked:          false,
		workItems:        &sync.Map{},
	}
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

// QueueWork queues work to do on the workChan to be processed
func (w *Queue) QueueWork(workToDo Work, options ...workOption) {
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
}

// Dequeue removes from the queue the work item with the specified id.  If the work item is in process, then an error is returned.
func (w *Queue) Dequeue(id uuid.UUID) error {
	if i, ok := w.workItems.Load(id); ok {
		wi := i.(*workItem)
		if wi.state.Load() == int32(IN_QUEUE) {
			w.workQueue.Remove(wi.id)
			w.workItems.Delete(wi.id)
		} else if wi.state.Load() == int32(IN_PROCESS) {
			return fmt.Errorf("cannot delete work item %v because it is in process", id.String())
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
	w.stopSignal <- struct{}{}
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
		close(w.stopSignal)
		close(w.workChan)
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
			case <-w.stopSignal:
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
					w.workQueue.AdjustPriorities()
				} else {
					// queue is full, block and wait for worker to finish a task then add work to queue
					<-workerSemaphore
					w.workQueue.AdjustPriorities()
					wtemp := heap.Pop(w.workQueue).(*workItem)
					workerCh <- wtemp
					w.workQueue.Push(work)
				}
			}
		case <-workerSemaphore:
			// worker done, pop and start next work (if anything in queue)
			if w.workQueue.Len() > 0 {
				w.workQueue.AdjustPriorities()
				wtemp := heap.Pop(w.workQueue).(*workItem)
				workerCh <- wtemp
			}
		case <-w.stopSignal:
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
		wi.state.Store(int32(IN_PROCESS))
		err := wi.workToDo()
		if err != nil {
			w.errChan <- err
		}
		w.workItems.Delete(wi.id)
		semaphore <- true
	}
}
