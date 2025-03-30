/*
 * Copyright (c) 2024  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package sliceOps

// Cut removes the elements s[i:j] from the slice and returns them, while preserving the order of s.
func Cut[T any](s *[]T, i, j int) []T {
	result := make([]T, j-i)
	copy(result, (*s)[i:j])
	Remove(s, i, j)
	return result
}

// Remove removes the elements s[i:j] from the slice, while preserving the order of s and zeroing out the remaining elements to prevent memory leak.
func Remove[T any](s *[]T, i, j int) {
	var zeroVal T
	original := *s
	copy(original[i:], original[j:])
	// Zero out the remaining elements to prevent a memory leak.
	for k, n := len(original)-j+i, len(original); k < n; k++ {
		original[k] = zeroVal
	}
	*s = original[:len(original)-j+i]
}

// Insert inserts the elements v into the slice s at index i.
func Insert[T any](s []T, i int, v ...T) []T {
	return append(s[:i], append(v, s[i:]...)...)
}

// FilterInPlace removes the elements of s that do not satisfy the keep predicate.
func FilterInPlace[T any](s *[]T, keep func(T) bool) {
	i := 0
	for _, e := range *s {
		if keep(e) {
			(*s)[i] = e
			i++
		}
	}
	// Zero out the remaining elements to prevent a memory leak.
	var zeroVal T
	for k, n := i, len(*s); k < n; k++ {
		(*s)[k] = zeroVal
	}
	*s = (*s)[:i]
}

// Push appends the elements v to the slice s.
func Push[T any](s *[]T, v ...T) {
	n := make([]T, 0, len(*s)+len(v))
	n = append(n, v...)
	*s = append(n, *s...)
}

// Pop removes the first element of the slice s and returns it.
func Pop[T any](s *[]T) T {
	if len(*s) == 0 {
		var zeroVal T
		return zeroVal
	}
	result := (*s)[0]
	Remove(s, 0, 1)
	return result
}

// Distinct returns a new slice containing only the distinct elements of s.
func Distinct[T comparable](s []T) []T {
	distinctMap := make(map[T]struct{})
	result := []T{}
	for _, e := range s {
		if _, exists := distinctMap[e]; !exists {
			distinctMap[e] = struct{}{}
			result = append(result, e)
		}
	}
	return result
}

// Union returns a new slice containing the distinct union of set of slices.
func Union[T comparable](slices ...[]T) []T {
	distinctMap := make(map[T]struct{})
	result := []T{}
	for _, s := range slices {
		for _, e := range s {
			if _, exists := distinctMap[e]; !exists {
				distinctMap[e] = struct{}{}
				result = append(result, e)
			}
		}
	}
	return result
}

// Intersection returns a new slice containing the distinct intersection of set of slices.
func Intersection[T comparable](slices ...[]T) []T {
	intersectionMap := make(map[T]int)
	result := []T{}
	for _, s := range slices {
		for _, e := range s {
			//nolint:gosimple,staticcheck // This is more readable than the suggested alternative
			if _, exists := intersectionMap[e]; exists {
				intersectionMap[e]++
			} else {
				intersectionMap[e] = 1
			}
		}
	}
	for k, v := range intersectionMap {
		if v == len(slices) {
			result = append(result, k)
		}
	}
	return result
}

// Difference returns a new slice containing distinct elements in s1 that are not in s2.
func Difference[T comparable](s1, s2 []T) []T {
	distinctMap := make(map[T]struct{})
	result := []T{}
	for _, e := range s2 {
		distinctMap[e] = struct{}{}
	}
	for _, e := range s1 {
		if _, exists := distinctMap[e]; !exists {
			result = append(result, e)
		}
	}
	return result
}

// Disjoin returns a new slice containing distinct elements not in common between the slices provided.
func Disjoin[T comparable](slices ...[]T) []T {
	if len(slices) == 0 {
		return []T{}
	}
	result := slices[0]
	removed := []T{}
	for i, s := range slices {
		if i == 0 {
			continue
		}
		r1 := Difference(result, s)
		r2 := Difference(s, result)
		result = Union(r1, r2)
		// Remove any elements that were in the original result but are not in the new result.
		removed = append(removed, Distinct(Difference(s, result))...)
		result = Difference(result, removed)
	}
	return result
}
