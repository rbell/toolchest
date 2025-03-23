/*
 * Copyright (c) 2025 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package generic

import (
	"iter"
	"sync"
)

// SyncMap is a generic wrapper around sync.Map that provides type safety
type SyncMap[K comparable, V any] struct {
	*sync.Map
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		Map: &sync.Map{},
	}
}

// Load returns the value stored in the map for a key, or the zero value if no
// value is present. The ok result indicates whether value was found in the map.
func (m *SyncMap[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.Map.Load(key)
	if !ok {
		return
	}

	return v.(V), ok
}

// Store sets the value for a key.
func (m *SyncMap[K, V]) Store(key K, value V) {
	m.Map.Store(key, value)
}

// Swap swaps the value for a key and returns the previous value if any. The loaded result reports whether the key was present.
func (m *SyncMap[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	v, l := m.Map.Swap(key, value)
	if !l {
		return
	}
	return v.(V), l
}

// Delete deletes the value for a key.
func (m *SyncMap[K, V]) Delete(key K) {
	m.Map.Delete(key)
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value. The loaded result is true if the value was loaded, false if stored.
func (m *SyncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, l := m.Map.LoadOrStore(key, value)
	return v.(V), l
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (m *SyncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, l := m.Map.LoadAndDelete(key)
	if !l {
		return
	}
	return v.(V), l
}

// CompareAndDelete deletes the entry for key if its value is equal to old.
// The old value must be of a comparable type.
func (m *SyncMap[K, V]) CompareAndDelete(key K, old V) bool {
	return m.Map.CompareAndDelete(key, old)
}

// CompareAndSwap swaps the old and new values for key if the value stored in the map is equal to old.
func (m *SyncMap[K, V]) CompareAndSwap(key K, old V, new V) bool {
	return m.Map.CompareAndSwap(key, old, new)
}

// Range calls f sequentially for each key and value present in the map. If f returns false, range stops the iteration.
func (m *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	m.Map.Range(func(k any, v any) bool {
		return f(k.(K), v.(V))
	})
}

// Iterate returns an iterator that can be used to iterate over the map.
func (m *SyncMap[K, V]) Iterate() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for anyKey, anyValue := range m.Map.Range {
			if !yield(
				anyKey.(K),
				anyValue.(V),
			) {
				break
			}
		}
	}
}
