/*
 * Copyright (c) 2024  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package propositions

// MapContainsKey returns true if the map contains the key.
func MapContainsKey[T comparable, U any](m map[T]U, k T) bool {
	_, ok := m[k]
	return ok
}

// MapContainsValue returns true if the map contains the value.
func MapContainsValue[T comparable, U comparable](m map[T]U, v U) bool {
	for _, val := range m {
		if val == v {
			return true
		}
	}
	return false
}

// MapKeyPropositionAny returns true if any of the keys in the map satisfy the proposition.
func MapKeyPropositionAny[T comparable, U any](m map[T]U, p func(T) bool) bool {
	for k := range m {
		if p(k) {
			return true
		}
	}
	return false
}

// MapKeyPropositionAll returns true if all of the keys in the map satisfy the proposition.
func MapKeyPropositionAll[T comparable, U any](m map[T]U, p func(T) bool) bool {
	for k := range m {
		if !p(k) {
			return false
		}
	}
	return true
}

// MapKeyPropositionNone returns true if none of the keys in the map satisfy the proposition.
func MapKeyPropositionNone[T comparable, U any](m map[T]U, p func(T) bool) bool {
	for k := range m {
		if p(k) {
			return false
		}
	}
	return true
}

// MapValuePropositionAny returns true if any of the values in the map satisfy the proposition.
func MapValuePropositionAny[T comparable, U any](m map[T]U, p func(U) bool) bool {
	for _, v := range m {
		if p(v) {
			return true
		}
	}
	return false
}

// MapValuePropositionAll returns true if all of the values in the map satisfy the proposition.
func MapValuePropositionAll[T comparable, U any](m map[T]U, p func(U) bool) bool {
	for _, v := range m {
		if !p(v) {
			return false
		}
	}
	return true
}

// MapvaluePropositionNone returns true if none of the values in the map satisfy the proposition.
func MapValuePropositionNone[T comparable, U any](m map[T]U, p func(U) bool) bool {
	for _, v := range m {
		if p(v) {
			return false
		}
	}
	return true
}
