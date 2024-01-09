/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package storage

import (
	"container/heap"
	"github.com/rbell/toolchest/errors"
	"sync"
	"sync/atomic"
)

type stackEntry[T any] struct {
	id    uint64
	entry T
}

type GenericStack[T any] struct {
	stack      *stack[T]
	currentKey atomic.Uint64
}

// NewGenericStack returns an initialized reference to a GenericStack of T
func NewGenericStack[T any](initialSize int) *GenericStack[T] {
	return &GenericStack[T]{
		stack:      newStack[T](initialSize),
		currentKey: atomic.Uint64{},
	}
}

// Push pushes value of type T on the stack, returning the assigned id in the stack
func (s *GenericStack[T]) Push(value T) (id uint64) {
	v := &stackEntry[T]{
		id:    s.currentKey.Add(1),
		entry: value,
	}
	heap.Push(s.stack, v)
	return v.id
}

// Pop removes the next T from the stack and returns it
func (s *GenericStack[T]) Pop() T {
	return heap.Pop(s.stack).(*stackEntry[T]).entry
}

// Peek returns the value on the stack that was assigned the id requested.  IDNotFoundError returned if id not found.
func (s *GenericStack[T]) Peek(id uint64) (value T, err error) {
	for _, v := range s.stack.entries {
		if v.id == id {
			value = v.entry
			return
		}
	}
	err = &errors.NotFound{}
	return
}

// Len returns the number of elements on the stack
func (s *GenericStack[T]) Len() int {
	return s.stack.Len()
}

// Values returns a slice of all the values on the stack
func (s *GenericStack[T]) Values() []T {
	values := make([]T, 0, s.stack.Len())
	for _, v := range s.stack.entries {
		values = append(values, v.entry)
	}
	return values
}

// Implements container/heap, with push / pop acting in a FIFO order, where each element is a *stackEntry[T]
type stack[T any] struct {
	entries []*stackEntry[T]
	mux     *sync.Mutex
}

func (s *stack[T]) Len() int {
	return len(s.entries)
}
func (s *stack[T]) Less(i, j int) bool {
	return s.entries[i].id < s.entries[j].id
}
func (s *stack[T]) Swap(i, j int) {
	s.entries[i], s.entries[j] = s.entries[j], s.entries[i]
}

func newStack[T any](initialSize int) *stack[T] {
	result := &stack[T]{
		entries: make([]*stackEntry[T], 0, initialSize),
		mux:     &sync.Mutex{},
	}
	heap.Init(result)
	return result
}

// Push pushes x, which must be a *stackEntry[T], to the stack
func (s *stack[T]) Push(x any) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.entries = append(s.entries, x.(*stackEntry[T]))
}

// Pop pops and returns the next *stackEntry[T] from the stack
func (s *stack[T]) Pop() any {
	s.mux.Lock()
	defer s.mux.Unlock()

	old := s.entries
	var result *stackEntry[T]
	result, s.entries = old[len(old)-1], old[:len(old)-1]
	cpy := *result // dereference and return another reference to the value
	result = nil   // nil out the reference to the popped stackEntry in the backing array of the entries to protect memory
	return &cpy
}
