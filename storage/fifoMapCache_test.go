/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package storage

import (
	"context"
	"github.com/rbell/toolchest/propositions"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFifoMapCache_ReturnsBalancedFifoMapCache(t *testing.T) {
	// setup
	ctx := context.Background()

	// test
	m := NewFifoMapCache[int, int](ctx, 100)

	// assert
	assert.NotNil(t, m, "Expected map to be initialized")
	assert.Equalf(t, 10, m.maxPartitions, "Expected max size to be 10")
	assert.Equalf(t, 10, m.partitionCapacity, "Expected partition capacity to be 10")
	assert.NotNilf(t, m.partitions, "Expected partitions to be initialized")
	assert.NotNilf(t, m.valuePartitionIndex, "Expected value partition index to be initialized")
	assert.NotNilf(t, m.ctx, "Expected context to be initialized")
	assert.NotNilf(t, m.currentPartitionMux, "Expected current partition mutex to be initialized")
}

func TestNewFifoMapCache_WithBalancedPartitions_ReturnsCustomBalancedFifoMapCache(t *testing.T) {
	// setup
	ctx := context.Background()

	// test
	m := NewFifoMapCache[int, int](ctx, 100, WithBalancedPartitions(float64(1.5), 1))

	// assert
	assert.NotNil(t, m, "Expected map to be initialized")
	assert.Equalf(t, 21, m.maxPartitions, "Expected max size to be 21")
	assert.Equalf(t, 4, m.partitionCapacity, "Expected partition capacity to be 4")
	assert.NotNilf(t, m.partitions, "Expected partitions to be initialized")
	assert.NotNilf(t, m.valuePartitionIndex, "Expected value partition index to be initialized")
	assert.NotNilf(t, m.ctx, "Expected context to be initialized")
	assert.NotNilf(t, m.currentPartitionMux, "Expected current partition mutex to be initialized")
}

func TestNewFifoMapCache_WithBalancedPartitions_WithMinimumPartitions_ReturnsCustomBalancedFifoMapCacheWithMinimumPartiitons(t *testing.T) {
	// setup
	ctx := context.Background()

	// test
	m := NewFifoMapCache[int, int](ctx, 100, WithBalancedPartitions(float64(1.5), 50))

	// assert
	assert.NotNil(t, m, "Expected map to be initialized")
	assert.Equalf(t, 50, m.maxPartitions, "Expected max size to be 50")
	assert.Equalf(t, 2, m.partitionCapacity, "Expected partition capacity to be 2")
	assert.NotNilf(t, m.partitions, "Expected partitions to be initialized")
	assert.NotNilf(t, m.valuePartitionIndex, "Expected value partition index to be initialized")
	assert.NotNilf(t, m.ctx, "Expected context to be initialized")
	assert.NotNilf(t, m.currentPartitionMux, "Expected current partition mutex to be initialized")
}

func TestFifoMapCache_Set_AddsValueToMap(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 100)

	// test
	m.Set(1, 1)

	// assert
	currPartition, _ := m.partitions.Peek(m.currentPartitionId)
	assert.Equal(t, int64(1), m.count.Load(), "Expected count to be 1")
	assert.Equal(t, 1, currPartition.Get(1), "Expected value to be 1")
}

func TestFifoMapCache_Set_UpdatesExistingEntryInMap(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 100)

	// test
	m.Set(1, 1)
	m.Set(1, 2)

	// assert
	currPartition, _ := m.partitions.Peek(m.currentPartitionId)
	assert.Equal(t, int64(1), m.count.Load(), "Expected count to be 1")
	assert.Equal(t, 2, currPartition.Get(1), "Expected value to be 2")
}

func TestFifoMapCache_Get_KeyNotFound_ReturnsZeroValue(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 100)

	// test
	value := m.Get(1)

	// assert
	assert.Equal(t, 0, value, "Expected value to be 0")
}

func TestFifoMapCache_Get_KeyFound_ReturnsValue(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 100)
	m.Set(1, 1)

	// test
	value := m.Get(1)

	// assert
	assert.Equal(t, 1, value, "Expected value to be 1")
}

func TestFifoMapCache_SetAndSweep_AddsValueToMapAndEvictsRemovesOldest(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)

	// test
	for i := 0; i < 15; i++ {
		m.Set(i, i)
	}

	m.sweep()

	// assert
	assert.Equal(t, int64(9), m.count.Load(), "Expected count to be 9")
	assert.False(t, m.Contains(0), "Expected value 0 to be removed")
	assert.False(t, m.Contains(1), "Expected value 1 to be removed")
	assert.False(t, m.Contains(2), "Expected value 2 to be removed")
	assert.False(t, m.Contains(3), "Expected value 3 to be removed")
	assert.False(t, m.Contains(4), "Expected value 4 to be removed")
	assert.False(t, m.Contains(5), "Expected value 5 to be removed")
	assert.True(t, m.Contains(6), "Expected value 6 to be removed")
	assert.True(t, m.Contains(7), "Expected value 7 to be removed")
	assert.True(t, m.Contains(8), "Expected value 8 to be removed")
	assert.True(t, m.Contains(9), "Expected value 9 to be removed")
	assert.True(t, m.Contains(10), "Expected value 10 to be present")
	assert.True(t, m.Contains(11), "Expected value 10 to be present")
	assert.True(t, m.Contains(12), "Expected value 10 to be present")
	assert.True(t, m.Contains(13), "Expected value 10 to be present")
	assert.True(t, m.Contains(14), "Expected value 10 to be present")
}

func TestFifoMapCache_SetAndResize_AddsValueToMap(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 25)
	m.Set(1, 1)

	// test
	m.Resize(100)
	m.Set(2, 2)

	// assert
	assert.Equal(t, int64(2), m.count.Load(), "Expected count to be 2")
	assert.True(t, m.Contains(1), "Expected value 1 to be present")
	assert.True(t, m.Contains(2), "Expected value 2 to be present")
}

func TestFifoMapCache_GetAfterSweptKey_ReturnsZeroValue(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)

	// test
	for i := 0; i < 15; i++ {
		m.Set(i, i)
	}
	m.sweep()
	value := m.Get(1)

	// assert
	assert.Equal(t, 0, value, "Expected value to be 0")
}

func TestFifoMapCache_GetAfterSweep_ReturnsNonZeroValue(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)

	// test
	for i := 0; i < 15; i++ {
		m.Set(i, i)
	}
	m.sweep()
	value := m.Get(12)

	// assert
	assert.Equal(t, 12, value, "Expected value to be 12")
}

func TestFifoMapCache_GetAfterResize_ReturnsNonZeroValue(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 100)

	// test
	for i := 0; i < 110; i++ {
		m.Set(i, i)
	}
	m.Resize(25) // rsize to 25 keeping last 25 items added
	value := m.Get(105)

	// assert
	assert.Equal(t, 105, value, "Expected value to be 105")
}

func TestFifoMapCache_ContainsAfterSweep_ReturnsFalse(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)

	// test
	for i := 0; i < 15; i++ {
		m.Set(i, i)
	}
	m.sweep()
	contains := m.Contains(1)

	// assert
	assert.False(t, contains, "Expected value to be false")
}

func TestFifoMapCache_Delete_Deletes(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)
	m.Set(1, 1)

	// test
	m.Delete(1)

	// assert
	assert.Equal(t, int64(0), m.count.Load(), "Expected count to be 0")
	assert.False(t, m.Contains(1), "Expected value to be deleted")
}

func TestFifoMapCache_DeleteAfterSweep_Deletes(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)
	for i := 0; i < 15; i++ {
		m.Set(i, i)
	}
	m.sweep()

	// test
	m.Delete(12)

	// assert
	assert.Equal(t, int64(8), m.count.Load(), "Expected count to be 8")
	assert.False(t, m.Contains(12), "Expected value to be deleted")
}

func TestFifoMapCache_ContainsAfterSweep_ReturnsTrue(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)

	// test
	for i := 0; i < 15; i++ {
		m.Set(i, i)
	}
	m.sweep()
	contains := m.Contains(12)

	// assert
	assert.True(t, contains, "Expected value to be true")
}

func TestFifoMapCache_Len_ReturnsZero(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)

	// test
	length := m.Len()

	// assert
	assert.Equal(t, 0, length, "Expected length to be 0")
}

func TestFifoMapCache_Len_ReturnsNonZero(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)
	m.Set(1, 1)

	// test
	length := m.Len()

	// assert
	assert.Equal(t, 1, length, "Expected length to be 1")
}

func TestFifoMapCache_LenAfterSweep_ReturnsNonZero(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)
	for i := 0; i < 15; i++ {
		m.Set(i, i)
	}
	m.sweep()

	// test
	length := m.Len()

	// assert
	assert.Equal(t, 9, length, "Expected length to be 9")
}

func TestFifoMapCache_LenAfterDelete_ReturnsNonZero(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)
	for i := 0; i < 5; i++ {
		m.Set(i, i)
	}
	m.Delete(4)

	// test
	length := m.Len()

	// assert
	assert.Equal(t, 4, length, "Expected length to be 4")
}

func TestFifoMapCache_Keys_ReturnsEmptySlice(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)

	// test
	keys := m.Keys()

	// assert
	assert.Len(t, keys, 0, "Expected keys to be empty")
}

func TestFifoMapCache_Keys_ReturnsKeysInOrder(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)
	m.Set(2, 2)
	m.Set(1, 1)

	// test
	keys := m.Keys()

	// assert
	assert.Len(t, keys, 2, "Expected keys to have 2 entries")
	assert.Equal(t, 2, keys[0], "Expected first key to be 2")
	assert.Equal(t, 1, keys[1], "Expected second key to be 1")
}

func TestFifoMapCache_KeysAfterSweep_ReturnsKeys(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)
	for i := 0; i < 15; i++ {
		m.Set(i, i)
	}
	m.sweep()

	// test
	keys := m.Keys()

	// assert
	assert.Len(t, keys, 9, "Expected keys to have 9 entries")
	assert.True(t, propositions.SliceContainsAll(keys, []int{6, 7, 8, 9, 10, 11, 12, 13, 14}), "Expected keys to contain all values")
}

func TestFifoMapCache_KeysAfterDelete_ReturnsKeysInOrder(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)
	for i := 0; i < 5; i++ {
		m.Set(i, i)
	}
	m.Delete(4)

	// test
	keys := m.Keys()

	// assert
	assert.Len(t, keys, 4, "Expected keys to have 4 entries")
	assert.True(t, propositions.SliceContainsAll(keys, []int{0, 1, 2, 3}), "Expected keys to contain all values")
}

func TestFifoMapCache_Clear_ClearsMap(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)
	m.Set(1, 1)

	// test
	m.Clear()

	// assert
	assert.Equal(t, int64(0), m.count.Load(), "Expected count to be 0")
	assert.False(t, m.Contains(1), "Expected value to be deleted")
}

func TestFifoMapCache_ClearAfterSweep_ClearsMap(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)
	for i := 0; i < 15; i++ {
		m.Set(i, i)
	}
	m.sweep()

	// test
	m.Clear()

	// assert
	assert.Equal(t, int64(0), m.count.Load(), "Expected count to be 0")
	assert.False(t, m.Contains(10), "Expected value to be deleted")
	assert.False(t, m.Contains(11), "Expected value to be deleted")
	assert.False(t, m.Contains(12), "Expected value to be deleted")
	assert.False(t, m.Contains(13), "Expected value to be deleted")
	assert.False(t, m.Contains(14), "Expected value to be deleted")
	assert.False(t, m.Contains(5), "Expected value to be deleted")
	assert.False(t, m.Contains(6), "Expected value to be deleted")
	assert.False(t, m.Contains(7), "Expected value to be deleted")
	assert.False(t, m.Contains(8), "Expected value to be deleted")
	assert.False(t, m.Contains(9), "Expected value to be deleted")
}

func TestFifoMapCache_ClearAfterDelete_ClearsMap(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)
	for i := 0; i < 5; i++ {
		m.Set(i, i)
	}
	m.Delete(4)

	// test
	m.Clear()

	// assert
	assert.Equal(t, int64(0), m.count.Load(), "Expected count to be 0")
	assert.False(t, m.Contains(0), "Expected value to be deleted")
	assert.False(t, m.Contains(1), "Expected value to be deleted")
	assert.False(t, m.Contains(2), "Expected value to be deleted")
	assert.False(t, m.Contains(3), "Expected value to be deleted")
}

func TestFifoMapCache_Values_ReturnsEmptySlice(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)

	// test
	values := m.Values()

	// assert
	assert.Len(t, values, 0, "Expected values to be empty")
}

func TestFifoMapCache_Values_ReturnsValuesInOrder(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)
	m.Set(2, 2)
	m.Set(1, 1)

	// test
	values := m.Values()

	// assert
	assert.Len(t, values, 2, "Expected values to have 2 entries")
	assert.Equal(t, 2, values[0], "Expected first value to be 2")
	assert.Equal(t, 1, values[1], "Expected second value to be 1")
}

func TestFifoMapCache_ValuesAfterSweep_ReturnsValuesInOrder(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)
	for i := 0; i < 15; i++ {
		m.Set(i, i)
	}
	m.sweep()

	// test
	values := m.Values()

	// assert
	assert.Len(t, values, 9, "Expected values to have 9 entries")
	assert.True(t, propositions.SliceContainsAll(values, []int{6, 7, 8, 9, 10, 11, 12, 13, 14}), "Expected values to contain all values")
}

func TestFifoMapCache_ValuesAfterDelete_ReturnsValuesInOrder(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 10)
	for i := 0; i < 5; i++ {
		m.Set(i, i)
	}
	m.Delete(4)

	// test
	values := m.Values()

	// assert
	assert.Len(t, values, 4, "Expected values to have 4 entries")
	assert.Equal(t, 0, values[0], "Expected first value to be 0")
	assert.Equal(t, 1, values[1], "Expected second value to be 1")
	assert.Equal(t, 2, values[2], "Expected third value to be 2")
	assert.Equal(t, 3, values[3], "Expected fourth value to be 3")
}

func TestFifoMapCache_Resize(t *testing.T) {
	// setup
	ctx := context.Background()
	m := NewFifoMapCache[int, int](ctx, 25)

	// test
	m.Resize(100)

	// assert
	assert.Equal(t, 10, m.maxPartitions, "Expected max size to be 100")
	assert.Equal(t, 10, m.partitionCapacity, "Expected partition capacity to be 10")
	assert.NotNilf(t, m.partitions, "Expected partitions to be initialized")
	assert.NotNilf(t, m.valuePartitionIndex, "Expected value partition index to be initialized")
	assert.NotNilf(t, m.ctx, "Expected context to be initialized")
	assert.NotNilf(t, m.currentPartitionMux, "Expected current partition mutex to be initialized")
}
