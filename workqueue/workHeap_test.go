/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package workqueue

import (
	"container/heap"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWorkHeap_Push(t *testing.T) {
	// setup
	wh := make(workHeap, 0, 10)
	heap.Init(&wh)

	work := &workItem{
		workToDo: func() error {
			return nil
		},
		priority: 1,
		index:    1,
	}

	// test
	heap.Push(&wh, work)

	// assert
	assert.Len(t, wh, 1)
}

func TestWorkHeap_Pop(t *testing.T) {
	// setup
	wh := make(workHeap, 0, 10)
	heap.Init(&wh)

	work := &workItem{
		workToDo: func() error {
			return nil
		},
		priority: 1,
		index:    1,
	}

	heap.Push(&wh, work)

	// test
	result := heap.Pop(&wh)

	// assert
	assert.Equal(t, work, result)
	assert.Len(t, wh, 0)
}

func TestWorkHeap_Priority_Pop(t *testing.T) {
	// setup
	wh := make(workHeap, 0, 10)
	heap.Init(&wh)

	work1 := &workItem{
		workToDo: func() error {
			return nil
		},
		priority: 2,
		index:    1,
	}
	heap.Push(&wh, work1)
	work2 := &workItem{
		workToDo: func() error {
			return nil
		},
		priority: 1,
		index:    1,
	}
	heap.Push(&wh, work2)
	work3 := &workItem{
		workToDo: func() error {
			return nil
		},
		priority: 2,
		index:    1,
	}
	heap.Push(&wh, work3)

	// test
	result := heap.Pop(&wh)

	// assert
	// Pop should pop priority 1 over work with a greater priority number
	assert.Equal(t, work2, result)
	assert.Len(t, wh, 2)
}

func TestWorkHeap_AdjustPriorities_ChangesPriorities(t *testing.T) {
	// setup
	wh := make(workHeap, 0, 10)
	heap.Init(&wh)

	work1 := &workItem{
		workToDo: func() error {
			return nil
		},
		adjustPriority: func() int {
			return 1
		},
		priority: 2,
		index:    1,
	}
	heap.Push(&wh, work1)
	work2 := &workItem{
		workToDo: func() error {
			return nil
		},
		adjustPriority: func() int {
			return 2
		},
		priority: 1,
		index:    1,
	}
	heap.Push(&wh, work2)
	work3 := &workItem{
		workToDo: func() error {
			return nil
		},
		priority: 2,
		index:    1,
	}
	heap.Push(&wh, work3)

	// test
	wh.AdjustPriorities()
	result := heap.Pop(&wh)

	// assert
	// Pop should pop priority 1 over work with a greater priority number
	assert.Equal(t, work1, result)
	assert.Len(t, wh, 2)
}
