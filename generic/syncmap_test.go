/*
 * Copyright (c) 2025 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package generic

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyncMap_Load(t *testing.T) {
	m := SyncMap[string, int]{Map: &sync.Map{}}
	m.Store("key1", 1)
	value, ok := m.Load("key1")
	assert.True(t, ok)
	assert.Equal(t, 1, value)

	value, ok = m.Load("key2")
	assert.False(t, ok)
	assert.Equal(t, 0, value)
}

func TestSyncMap_Store(t *testing.T) {
	m := SyncMap[string, int]{Map: &sync.Map{}}
	m.Store("key1", 1)
	value, ok := m.Load("key1")
	assert.True(t, ok)
	assert.Equal(t, 1, value)

	m.Store("key1", 2)
	value, ok = m.Load("key1")
	assert.True(t, ok)
	assert.Equal(t, 2, value)
}

func TestSyncMap_Swap(t *testing.T) {
	m := SyncMap[string, int]{Map: &sync.Map{}}
	m.Store("key1", 1)
	previous, loaded := m.Swap("key1", 2)
	assert.True(t, loaded)
	assert.Equal(t, 1, previous)

	value, ok := m.Load("key1")
	assert.True(t, ok)
	assert.Equal(t, 2, value)

	previous, loaded = m.Swap("key2", 3)
	assert.False(t, loaded)
	assert.Equal(t, 0, previous)

	value, ok = m.Load("key2")
	assert.True(t, ok)
	assert.Equal(t, 3, value)
}

func TestSyncMap_Delete(t *testing.T) {
	m := SyncMap[string, int]{Map: &sync.Map{}}
	m.Store("key1", 1)
	m.Delete("key1")
	_, ok := m.Load("key1")
	assert.False(t, ok)
}

func TestSyncMap_LoadOrStore(t *testing.T) {
	m := SyncMap[string, int]{Map: &sync.Map{}}
	actual, loaded := m.LoadOrStore("key1", 1)
	assert.False(t, loaded)
	assert.Equal(t, 1, actual)

	actual, loaded = m.LoadOrStore("key1", 2)
	assert.True(t, loaded)
	assert.Equal(t, 1, actual)
}

func TestSyncMap_LoadAndDelete(t *testing.T) {
	m := SyncMap[string, int]{Map: &sync.Map{}}
	m.Store("key1", 1)
	value, loaded := m.LoadAndDelete("key1")
	assert.True(t, loaded)
	assert.Equal(t, 1, value)

	_, ok := m.Load("key1")
	assert.False(t, ok)

	value, loaded = m.LoadAndDelete("key2")
	assert.False(t, loaded)
	assert.Equal(t, 0, value)
}

func TestSyncMap_CompareAndDelete(t *testing.T) {
	m := SyncMap[string, int]{Map: &sync.Map{}}
	m.Store("key1", 1)
	deleted := m.CompareAndDelete("key1", 1)
	assert.True(t, deleted)
	_, ok := m.Load("key1")
	assert.False(t, ok)

	m.Store("key2", 2)
	deleted = m.CompareAndDelete("key2", 3)
	assert.False(t, deleted)
	value, ok := m.Load("key2")
	assert.True(t, ok)
	assert.Equal(t, 2, value)
}

func TestSyncMap_CompareAndSwap(t *testing.T) {
	m := SyncMap[string, int]{Map: &sync.Map{}}
	m.Store("key1", 1)
	swapped := m.CompareAndSwap("key1", 1, 2)
	assert.True(t, swapped)
	value, ok := m.Load("key1")
	assert.True(t, ok)
	assert.Equal(t, 2, value)

	swapped = m.CompareAndSwap("key1", 2, 2)
	assert.True(t, swapped)
	value, ok = m.Load("key1")
	assert.True(t, ok)
	assert.Equal(t, 2, value)
}

func TestSyncMap_Range(t *testing.T) {
	m := SyncMap[string, int]{Map: &sync.Map{}}
	m.Store("key1", 1)
	m.Store("key2", 2)
	m.Store("key3", 3)

	count := 0
	m.Range(func(key string, value int) bool {
		count++
		assert.Contains(t, []string{"key1", "key2", "key3"}, key)
		assert.Contains(t, []int{1, 2, 3}, value)
		return true
	})
	assert.Equal(t, 3, count)

	count = 0
	m.Range(func(key string, value int) bool {
		count++
		return key != "key2"
	})
	assert.Equal(t, 2, count)
}

func TestSyncMap_Iterate(t *testing.T) {
	m := SyncMap[string, int]{Map: &sync.Map{}}
	m.Store("key1", 1)
	m.Store("key2", 2)
	m.Store("key3", 3)

	count := 0
	it := m.Iterate()
	it(func(key string, value int) bool {
		count++
		assert.Contains(t, []string{"key1", "key2", "key3"}, key)
		assert.Contains(t, []int{1, 2, 3}, value)
		return true
	})
	assert.Equal(t, 3, count)

	brokeEarly := false
	fmt.Println("test")
	it = m.Iterate()
	it(func(key string, value int) bool {
		count++
		if key == "key2" {
			brokeEarly = true
			return false
		}
		return true
	})
	assert.True(t, brokeEarly)
}
