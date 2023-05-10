/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package workqueue

import (
	"container/heap"
	"runtime"
	"sync"
	"sync/atomic"
)

type Work func() error

type WorkQueueOption func(*Queue)

type workOption func(item *workItem)

// Queue allow work to be queued up and worked on in a set number of go routines
type Queue struct {
	workerCount      int
	queueLength      *atomic.Int32
	workChan         chan *workItem
	workQueue        workHeap
	errChan          chan error
	stopSignal       chan struct{}
	errSubScriberMux *sync.Mutex
	errorSubscribers []chan error
	stopped          atomic.Bool
	breaked          bool
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
	}
	wq.queueLength.Store(int32(runtime.NumCPU() * 2))
	wq.stopped.Store(false)
	for _, o := range options {
		o(wq)
	}

	wq.workChan = make(chan *workItem)
	wq.workQueue = make(workHeap, 0, wq.queueLength.Load())

	go wq.start()

	return wq
}

// QueueWork queues work to do on the workChan to be processed
func (w *Queue) QueueWork(workToDo Work, options ...workOption) {
	wi := &workItem{workToDo: workToDo, priority: 1}
	for _, option := range options {
		option(wi)
	}
	if !w.stopped.Load() {
		w.workChan <- wi
	}
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

// ResizeQueueLength adjusts the size of the queue
func (w *Queue) ResizeQueueLength(length int) {
	w.queueLength.Store(int32(length))
}

// Break stops the queue form accepting any work and any work in queue is skipped
func (w *Queue) Break() {
	w.breaked = true
	w.Stop()
}

func (w *Queue) start() {
	defer func() {
		close(w.errChan)
		close(w.stopSignal)
		close(w.workChan)
	}()

	heap.Init(&w.workQueue)

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
	workerCh := make(chan Work, w.workerCount)
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
					case workerCh <- work.workToDo:
						break insideFor
					default:
					}
				}
				// Workers are busy and placing the workToDo on the channel failed - queue it on prioritized queue
				if w.workQueue.Len() < int(w.queueLength.Load()) {
					heap.Push(&w.workQueue, work)
					w.workQueue.AdjustPriorities()
				} else {
					// queue is full, block and wait for worker to finish a task then add work to queue
					<-workerSemaphore
					w.workQueue.AdjustPriorities()
					wtemp := heap.Pop(&w.workQueue).(*workItem).workToDo
					workerCh <- wtemp
					w.workQueue.Push(work)
				}
			}
		case <-workerSemaphore:
			// worker done, pop and start next work (if anything in queue)
			if w.workQueue.Len() > 0 {
				w.workQueue.AdjustPriorities()
				wtemp := heap.Pop(&w.workQueue).(*workItem).workToDo
				workerCh <- wtemp
			}
		case <-w.stopSignal:
			break outsideFor
		}
	}

	// Finish any work left on queue
	for _, work := range w.workQueue {
		if !w.breaked {
			workerCh <- work.workToDo
		}
	}
}

func (w *Queue) doWork(workCh chan Work, semaphore chan bool) {
	for workToDo := range workCh {
		err := workToDo()
		if err != nil {
			w.errChan <- err
		}
		semaphore <- true
	}
}
