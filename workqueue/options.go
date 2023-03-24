/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package workqueue

// WithWorkers sets the number of go routines working on the workChan
func WithWorkers(workerCount int) WorkQueueOption {
	return func(queue *Queue) {
		queue.workerCount = workerCount
	}
}

// WithQueueLength sets the number of functions that can be queued up before routines queueing are blocked
func WithQueueLength(queueLength int) WorkQueueOption {
	return func(queue *Queue) {
		queue.queueLength = queueLength
	}
}

// WithPriority sets the priority of the work to be done.  Lower number is higher priority.
func WithPriority(priority int) workOption {
	return func(w *workItem) {
		w.priority = priority
	}
}

// WithAdjustPriority ads an adjustment priorty function to the workItem such that it's priority can be dynamically adjusted in queue
func WithAdjustPriority(adjustment func() int) workOption {
	return func(w *workItem) {
		w.adjustPriority = adjustment
	}
}
