/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package workqueue

import "container/heap"

type workItem struct {
	workToDo       Work
	adjustPriority func() int
	priority       int
	index          int
}

// workHeap implements heap interface for prioritizing work to be worked on when queued
type workHeap []*workItem

func (wh workHeap) AdjustPriorities() {
	for _, workItem := range wh {
		wi := workItem
		if wi.adjustPriority != nil {
			priority := wi.adjustPriority()
			if priority != wi.priority {
				wi.priority = priority
				heap.Fix(&wh, wi.index)
			}
		}
	}
}

// Len returns the length of the workHeap
func (wh workHeap) Len() int {
	return len(wh)
}

// Less returns true if the work item at index i is less than the work item at index j
func (wh workHeap) Less(i, j int) bool {
	return wh[i].priority < wh[j].priority
}

// Swap swaps the work items at index i and j
func (wh workHeap) Swap(i, j int) {
	wh[i], wh[j] = wh[j], wh[i]
	wh[i].index = i
	wh[j].index = j
}

// Push puts a work item on the heap
func (wh *workHeap) Push(x any) {
	n := len(*wh)
	item := x.(*workItem)
	item.index = n
	*wh = append(*wh, item)
}

// Pop pops the next highest priority work item off the heap and returns it
func (wh *workHeap) Pop() any {
	old := *wh
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*wh = old[0 : n-1]
	return item
}
