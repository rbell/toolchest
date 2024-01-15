/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package storage

import (
	"github.com/rbell/toolchest/propositions"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSafeMap_ZeroInitialCapacity_CreatesSafeMap(t *testing.T) {
	// setup

	// test
	m := NewSafeMap[int, string](0)

	// assert
	assert.Lenf(t, m.m, 0, "Expected map to be empty")
	assert.NotNil(t, m.mux, "Expected mutex to be initialized")
}

func TestNewSafeMap_NonZeroInitialCapacity_CreatesSafeMap(t *testing.T) {
	// setup

	// test
	m := NewSafeMap[int, string](10)

	// assert
	assert.Lenf(t, m.m, 0, "Expected map to be empty")
	assert.NotNil(t, m.mux, "Expected mutex to be initialized")
}

func TestSafeMap_Get_KeyNotFound_ReturnsZeroValue(t *testing.T) {
	// setup
	m := NewSafeMap[int, string](0)

	// test
	value := m.Get(1)

	// assert
	assert.Empty(t, value, "Expected value to be empty")
}

func TestSafeMap_Get_KeyFound_ReturnsValue(t *testing.T) {
	// setup
	m := NewSafeMap[int, string](0)
	m.Set(1, "test")

	// test
	value := m.Get(1)

	// assert
	assert.Equal(t, "test", value, "Expected value to be 'test'")
}

func TestSafeMap_Contains_KeyNotFound_ReturnsFalse(t *testing.T) {
	// setup
	m := NewSafeMap[int, string](0)

	// test
	contains := m.Contains(1)

	// assert
	assert.False(t, contains, "Expected contains to be false")
}

func TestSafeMap_Contains_KeyFound_ReturnsTrue(t *testing.T) {
	// setup
	m := NewSafeMap[int, string](0)
	m.Set(1, "test")

	// test
	contains := m.Contains(1)

	// assert
	assert.True(t, contains, "Expected contains to be true")
}

func TestSafeMap_Set_KeyNotFound_SetsValue(t *testing.T) {
	// setup
	m := NewSafeMap[int, string](0)

	// test
	m.Set(1, "test")

	// assert
	assert.Equal(t, "test", m.m[1], "Expected value to be 'test'")
}

func TestSafeMap_Set_KeyFound_UpdatesValue(t *testing.T) {
	// setup
	m := NewSafeMap[int, string](0)
	m.Set(1, "test")

	// test
	m.Set(1, "test2")

	// assert
	assert.Equal(t, "test2", m.m[1], "Expected value to be 'test2'")
}

func TestSafeMap_Delete_KeyNotFound_DoesNothing(t *testing.T) {
	// setup
	m := NewSafeMap[int, string](0)

	// test
	m.Delete(1)

	// assert
	assert.Lenf(t, m.m, 0, "Expected map to be empty")
}

func TestSafeMap_Delete_KeyFound_DeletesKey(t *testing.T) {
	// setup
	m := NewSafeMap[int, string](0)
	m.Set(1, "test")

	// test
	m.Delete(1)

	// assert
	assert.Lenf(t, m.m, 0, "Expected map to be empty")
}

func TestSafeMap_Delete_KeyFound_DeletesCorrectKey(t *testing.T) {
	// setup
	m := NewSafeMap[int, string](0)
	m.Set(1, "test")
	m.Set(2, "test2")

	// test
	m.Delete(1)

	// assert
	assert.Lenf(t, m.m, 1, "Expected map to have 1 entry")
	assert.Equal(t, "test2", m.m[2], "Expected value to be 'test2'")
}

func TestSafeMap_Len_EmptyMap_ReturnsZero(t *testing.T) {
	// setup
	m := NewSafeMap[int, string](0)

	// test
	length := m.Len()

	// assert
	assert.Equal(t, 0, length, "Expected length to be 0")
}

func TestSafeMap_Len_NonEmptyMap_ReturnsLength(t *testing.T) {
	// setup
	m := NewSafeMap[int, string](0)
	m.Set(1, "test")
	m.Set(2, "test2")

	// test
	length := m.Len()

	// assert
	assert.Equal(t, 2, length, "Expected length to be 2")
}

func TestSafeMap_Keys_EmptyMap_ReturnsEmptySlice(t *testing.T) {
	// setup
	m := NewSafeMap[int, string](0)

	// test
	keys := m.Keys()

	// assert
	assert.Lenf(t, keys, 0, "Expected keys to be empty")
}

func TestSafeMap_Keys_NonEmptyMap_ReturnsKeys(t *testing.T) {
	// setup
	m := NewSafeMap[int, string](0)
	m.Set(2, "test2")
	m.Set(1, "test")

	// test
	keys := m.Keys()

	// assert
	assert.Lenf(t, keys, 2, "Expected keys to have 2 entries")
	assert.True(t, propositions.SliceContainsAll(keys, []int{1, 2}), "Expected keys to contain 1 and 2")
}

func TestSafeMap_Keys_NonEmptyMap_ReturnsKeysAfterDelete(t *testing.T) {
	// setup
	m := NewSafeMap[int, string](0)
	m.Set(2, "test2")
	m.Set(1, "test")
	m.Delete(2)

	// test
	keys := m.Keys()

	// assert
	assert.Lenf(t, keys, 1, "Expected keys to have 1 entry")
	assert.Equal(t, 1, keys[0], "Expected first key to be 1")
}

func TestSafeMap_Values_EmptyMap_ReturnsEmptuySlice(t *testing.T) {
	// setup
	m := NewSafeMap[int, string](0)

	// test
	values := m.Values()

	// assert
	assert.Lenf(t, values, 0, "Expected values to be empty")
}

func TestSafeMap_Values_NonEmptyMap_ReturnsValuesInOrder(t *testing.T) {
	// setup
	m := NewSafeMap[int, string](0)
	m.Set(2, "test2")
	m.Set(1, "test")

	// test
	values := m.Values()

	// assert
	assert.Lenf(t, values, 2, "Expected values to have 2 entries")
	assert.True(t, propositions.SliceContainsAll(values, []string{"test", "test2"}), "Expected values to contain 'test' and 'test2'")
}

func TestSafeMap_Values_NonEmptyMap_ReturnsValuesInOrderAfterDelete(t *testing.T) {
	// setup
	m := NewSafeMap[int, string](0)
	m.Set(2, "test2")
	m.Set(1, "test")
	m.Delete(2)

	// test
	values := m.Values()

	// assert
	assert.Lenf(t, values, 1, "Expected values to have 1 entry")
	assert.Equal(t, "test", values[0], "Expected first value to be 'test'")
}
