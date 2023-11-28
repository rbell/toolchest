/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package storage

import (
	"context"
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
