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

func TestGenericStack_Values_ReturnsEmptySlice(t *testing.T) {
	// setup
	s := NewGenericStack[int](0)

	// test
	values := s.Values()

	// assert
	assert.Lenf(t, values, 0, "Expected values to be empty")
}

func TestGenericStack_Values_ReturnsValuesInOrder(t *testing.T) {
	// setup
	s := NewGenericStack[int](0)
	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Push(4)
	s.Push(5)
	s.Push(6)

	// test
	values := s.Values()

	// assert
	assert.Lenf(t, values, 6, "Expected values to have 6 entries")
	assert.Equal(t, 1, values[0], "Expected first value to be 1")
	assert.Equal(t, 2, values[1], "Expected second value to be 2")
	assert.Equal(t, 3, values[2], "Expected third value to be 3")
	assert.Equal(t, 4, values[3], "Expected fourth value to be 4")
	assert.Equal(t, 5, values[4], "Expected fifth value to be 5")
	assert.Equal(t, 6, values[5], "Expected sixth value to be 6")
}
