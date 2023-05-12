/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package workqueue

import (
	"container/heap"
	"sync"
)

// workHeap implements heap interface for prioritizing work to be worked on when queued
type workHeap struct {
	items []*workItem
	mux   *sync.RWMutex
}

func newWorkHeap(length int) *workHeap {
	return &workHeap{
		items: make([]*workItem, 0, length),
		mux:   &sync.RWMutex{},
	}
}

func (wh workHeap) AdjustPriorities() {
	for _, workItem := range wh.items {
		wi := workItem
		if wi.adjustPriority != nil {
			newPriority := wi.adjustPriority()
			if newPriority != wi.priority {
				wi.priority = newPriority
				heap.Fix(&wh, wi.position)
			}
		}
	}
}

// Len returns the length of the workHeap
func (wh *workHeap) Len() int {
	wh.mux.RLock()
	defer wh.mux.RUnlock()
	return len(wh.items)
}

// Less returns true if the work item at index i is less than the work item at index j
func (wh *workHeap) Less(i, j int) bool {
	wh.mux.RLock()
	defer wh.mux.RUnlock()
	return wh.items[i].priority < wh.items[j].priority
}

// Swap swaps the work items at index i and j
func (wh *workHeap) Swap(i, j int) {
	wh.mux.Lock()
	defer wh.mux.Unlock()
	wh.items[i], wh.items[j] = wh.items[j], wh.items[i]
	wh.items[i].position = i
	wh.items[j].position = j
}

// Push puts a work item on the heap
func (wh *workHeap) Push(x any) {
	wh.mux.Lock()
	defer wh.mux.Unlock()
	n := len(wh.items)
	item := x.(*workItem)
	item.state.Store(int32(IN_QUEUE))
	item.position = n
	wh.items = append(wh.items, item)
}

// Pop pops the next highest priority work item off the heap and returns it
func (wh *workHeap) Pop() any {
	wh.mux.Lock()
	defer wh.mux.Unlock()
	old := wh.items
	n := len(old)
	item := old[n-1]
	old[n-1] = nil     // avoid memory leak
	item.position = -1 // for safety
	wh.items = old[0 : n-1]
	return item
}

// Remove removes the work item with the id
func (wh *workHeap) Remove(position int) {
	if position >= 0 {
		heap.Remove(wh, position)
	}
}
