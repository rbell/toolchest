/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package workqueue

import (
	"container/heap"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkHeap_Push(t *testing.T) {
	// setup
	wh := newWorkHeap(10)
	heap.Init(wh)

	work := &workItem{
		workToDo: func() error {
			return nil
		},
		QueuedWork: &QueuedWork{
			priority: 1,
			index:    -1,
			state:    &atomic.Int32{},
		},
	}

	// test
	heap.Push(wh, work)

	// assert
	assert.Len(t, wh.items, 1)
}

func TestWorkHeap_Pop(t *testing.T) {
	// setup
	wh := newWorkHeap(10)
	heap.Init(wh)

	work := &workItem{
		workToDo: func() error {
			return nil
		},
		QueuedWork: &QueuedWork{
			priority: 1,
			index:    -1,
			state:    &atomic.Int32{},
		},
	}

	heap.Push(wh, work)

	// test
	result := heap.Pop(wh)

	// assert
	assert.Equal(t, work, result)
	assert.Len(t, wh.items, 0)
}

func TestWorkHeap_Priority_Pop(t *testing.T) {
	// setup
	wh := newWorkHeap(10)
	heap.Init(wh)

	work1 := &workItem{
		workToDo: func() error {
			return nil
		},
		QueuedWork: &QueuedWork{
			priority: 1,
			index:    -1,
			state:    &atomic.Int32{},
		},
	}
	heap.Push(wh, work1)
	work2 := &workItem{
		workToDo: func() error {
			return nil
		},
		QueuedWork: &QueuedWork{
			priority: 1,
			index:    -1,
			state:    &atomic.Int32{},
		},
	}
	heap.Push(wh, work2)
	work3 := &workItem{
		workToDo: func() error {
			return nil
		},
		QueuedWork: &QueuedWork{
			priority: 2,
			index:    -1,
			state:    &atomic.Int32{},
		},
	}
	heap.Push(wh, work3)

	// test
	result := heap.Pop(wh)

	// assert
	// Pop should pop priority 1 over work with a greater priority number
	expected := work2.QueuedWork
	expected.index = -1
	assert.Equal(t, expected, result.(*workItem).QueuedWork)
	assert.Len(t, wh.items, 2)
}

func TestWorkHeap_AdjustPriorities_ChangesPriorities(t *testing.T) {
	// setup
	wh := newWorkHeap(10)
	heap.Init(wh)

	work1 := &workItem{
		workToDo: func() error {
			return nil
		},
		adjustPriority: func() int {
			return 1
		},
		QueuedWork: &QueuedWork{
			priority: 2,
			index:    -1,
			state:    &atomic.Int32{},
		},
	}
	heap.Push(wh, work1)
	work2 := &workItem{
		workToDo: func() error {
			return nil
		},
		adjustPriority: func() int {
			return 2
		},
		QueuedWork: &QueuedWork{
			priority: 1,
			index:    -1,
			state:    &atomic.Int32{},
		},
	}
	heap.Push(wh, work2)
	work3 := &workItem{
		workToDo: func() error {
			return nil
		},
		QueuedWork: &QueuedWork{
			priority: 2,
			index:    -1,
			state:    &atomic.Int32{},
		},
	}
	heap.Push(wh, work3)

	// test
	wh.AdjustPriorities()
	result := heap.Pop(wh)

	// assert
	// Pop should pop priority 1 over work with a greater priority number
	assert.Equal(t, work1, result)
	assert.Len(t, wh.items, 2)
}
