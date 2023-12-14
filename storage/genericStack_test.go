/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGenericStack_ZeroInitialSize_ReturnsGenericStack(t *testing.T) {
	// setup

	// test
	s := NewGenericStack[int](0)

	// assert
	assert.Lenf(t, s.stack.entries, 0, "Expected stack to be empty")
	assert.NotNil(t, s.stack.mux, "Expected mutex to be initialized")
}

func TestNewGenericStack_NonZeroInitialSize_ReturnsGenericStack(t *testing.T) {
	// setup

	// test
	s := NewGenericStack[int](10)

	// assert
	assert.Lenf(t, s.stack.entries, 0, "Expected stack to be empty")
	assert.NotNil(t, s.stack.mux, "Expected mutex to be initialized")
}

func TestGenericStack_Push_AddsValueToStack(t *testing.T) {
	// setup
	s := NewGenericStack[int](0)

	// test
	s.Push(1)

	// assert
	assert.Lenf(t, s.stack.entries, 1, "Expected stack to have 1 entry")
}

func TestGenericStack_Push_ReturnsId(t *testing.T) {
	// setup
	s := NewGenericStack[int](0)

	// test
	id := s.Push(1)

	// assert
	assert.Equal(t, uint64(1), id, "Expected id to be 1")
}

func TestGenericStack_Pop_RemovesValueFromStack(t *testing.T) {
	// setup
	s := NewGenericStack[int](0)
	s.Push(1)

	// test
	s.Pop()

	// assert
	assert.Lenf(t, s.stack.entries, 0, "Expected stack to be empty")
}

func TestGenericStack_Pop_ReturnsFirstValuePushed(t *testing.T) {
	// setup
	s := NewGenericStack[int](0)
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Push(4)
	s.Push(5)
	s.Push(6)

	// test
	value := s.Pop()

	// assert
	assert.Equal(t, 1, value, "Expected value to be 1")
}

func TestGenericStack_Peek_IdNotFound_ReturnsError(t *testing.T) {
	// setup
	s := NewGenericStack[int](0)

	// test
	_, err := s.Peek(1)

	// assert
	assert.Error(t, err, "Expected error to be returned")
}

func TestGenericStack_Peek_IdFound_ReturnsValue(t *testing.T) {
	// setup
	s := NewGenericStack[int](0)
	s.Push(1)

	// test
	value, _ := s.Peek(1)

	// assert
	assert.Equal(t, 1, value, "Expected value to be 1")
}

func TestGenericStack_Len_ReturnsNumberOfEntries(t *testing.T) {
	// setup
	s := NewGenericStack[int](0)
	s.Push(1)

	// test
	length := s.Len()

	// assert
	assert.Equal(t, 1, length, "Expected length to be 1")
}

func TestGenericStack_Len_ReturnsZero(t *testing.T) {
	// setup
	s := NewGenericStack[int](0)

	// test
	length := s.Len()

	// assert
	assert.Equal(t, 0, length, "Expected length to be 0")
}

func TestGenericStack_Len_ReturnsNumberOfEntriesAfterPop(t *testing.T) {
	// setup
	s := NewGenericStack[int](0)
	s.Push(1)
	s.Pop()

	// test
	length := s.Len()

	// assert
	assert.Equal(t, 0, length, "Expected length to be 0")
}

func TestGenericStack_Len_ReturnsNumberOfEntriesAfterPush(t *testing.T) {
	// setup
	s := NewGenericStack[int](0)
	s.Push(1)
	s.Push(2)

	// test
	length := s.Len()

	// assert
	assert.Equal(t, 2, length, "Expected length to be 2")
}

func TestGenericStack_Len_ReturnsNumberOfEntriesAfterPushAndPop(t *testing.T) {
	// setup
	s := NewGenericStack[int](0)
	s.Push(1)
	s.Push(2)
	s.Pop()

	// test
	length := s.Len()

	// assert
	assert.Equal(t, 1, length, "Expected length to be 1")
}
