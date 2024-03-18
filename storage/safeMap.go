/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package storage

import "sync"

// SafeMap wraps a map[K]V with a simple RWMutex to facilitate concurrency
type SafeMap[K comparable, V any] struct {
	mux *sync.RWMutex
	m   map[K]V
}

// NewSafeMap returns an initialized reference to a SafeMap of K, V
func NewSafeMap[K comparable, V any](initialCapacity int) *SafeMap[K, V] {
	var m map[K]V
	if initialCapacity > 0 {
		m = make(map[K]V, initialCapacity)
	} else {
		m = make(map[K]V)
	}

	return &SafeMap[K, V]{
		mux: &sync.RWMutex{},
		m:   m,
	}
}

// Contains returns true if the key of type K is in the map
func (s *SafeMap[K, V]) Contains(key K) bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	_, ok := s.m[key]
	return ok
}

// Get returns the value of type V for the key of type K.  If the key is not found, the zero value of V is returned.
func (s *SafeMap[K, V]) Get(key K) (value V) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	value = s.m[key]
	return
}

// GetOrAdd returns the value of type V for the key of type K.  If the key is not found, the value is added to the map and returned.
func (s *SafeMap[K, V]) GetOrAdd(key K, val V) (value V) {
	if s.Has(key) {
		return s.Get(key)
	}
	s.mux.Lock()
	defer s.mux.Unlock()
	if v, ok := s.m[key]; ok {
		return v
	}
	s.m[key] = val
	return val
}

// Set sets the value of type V for the key of type K.
func (s *SafeMap[K, V]) Set(key K, value V) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.m[key] = value
}

// Delete deletes the key of type K from the map
func (s *SafeMap[K, V]) Delete(key K) {
	s.mux.Lock()
	defer s.mux.Unlock()
	delete(s.m, key)
}

// Clear removes all the keys and values from the map
func (s *SafeMap[K, V]) Clear() {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.m = make(map[K]V)
}

// ClearAndResize clears the map and resize it to the new size
func (s *SafeMap[K, V]) ClearAndResize(newSize int) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.m = make(map[K]V, newSize)
}

// Has returns true if the key of type K is in the map
func (s *SafeMap[K, V]) Has(key K) bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	_, ok := s.m[key]
	return ok
}

// Len returns the length of the map
func (s *SafeMap[K, V]) Len() int {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return len(s.m)
}

// Keys returns a slice of all the keys in the map
func (s *SafeMap[K, V]) Keys() []K {
	s.mux.RLock()
	defer s.mux.RUnlock()
	keys := make([]K, 0, len(s.m))
	for k := range s.m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of all the values in the map
func (s *SafeMap[K, V]) Values() []V {
	s.mux.RLock()
	defer s.mux.RUnlock()

	values := make([]V, 0, len(s.m))
	for _, v := range s.m {
		values = append(values, v)
	}
	return values
}

// CopyToMap returns a copy of the map
func (s *SafeMap[K, V]) CopyToMap() map[K]V {
	s.mux.RLock()
	defer s.mux.RUnlock()
	result := make(map[K]V, len(s.m))
	for k, v := range s.m {
		result[k] = v
	}
	return result
}

// TranslateToMapOf returns a map of type D from the map of type V
func TranslateToMapOf[K comparable, V any, D any](s *SafeMap[K, V], translator func(V) D) map[K]D {
	s.mux.RLock()
	defer s.mux.RUnlock()
	result := make(map[K]D, len(s.m))
	for k, v := range s.m {
		result[k] = translator(v)
	}
	return result
}
