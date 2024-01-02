/*
 * Copyright (c) 2023-2024  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package propositions

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSliceContainsValue(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	result := SliceContains(s, 3)
	assert.True(t, result)

	result = SliceContains(s, 6)
	assert.False(t, result)
}

func TestSliceContainsAllValues(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	v := []int{2, 3}
	result := SliceContainsAll(s, v)
	assert.True(t, result)

	v = []int{2, 6}
	result = SliceContainsAll(s, v)
	assert.False(t, result)
}

func TestSliceContainsAnyValues(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	v := []int{2, 6}
	result := SliceContainsAny(s, v)
	assert.True(t, result)

	v = []int{6, 7}
	result = SliceContainsAny(s, v)
	assert.False(t, result)
}

func TestSliceContainsNoValues(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	v := []int{6, 7}
	result := SliceContainsNone(s, v)
	assert.True(t, result)

	v = []int{2, 6}
	result = SliceContainsNone(s, v)
	assert.False(t, result)
}

func TestSliceAllValuesLessThan(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	result := SliceAllLessThan(s, 6)
	assert.True(t, result)

	result = SliceAllLessThan(s, 4)
	assert.False(t, result)
}

func TestSliceAllValuesLessThanOrEqualTo(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	result := SliceAllLessThanOrEqualTo(s, 5)
	assert.True(t, result)

	result = SliceAllLessThanOrEqualTo(s, 4)
	assert.False(t, result)
}

func TestSliceAllValuesGreaterThan(t *testing.T) {
	s := []int{2, 3, 4, 5, 6}
	result := SliceAllGreaterThan(s, 1)
	assert.True(t, result)

	result = SliceAllGreaterThan(s, 2)
	assert.False(t, result)
}

func TestSliceAllValuesGreaterThanOrEqualTo(t *testing.T) {
	s := []int{2, 3, 4, 5, 6}
	result := SliceAllGreaterThanOrEqualTo(s, 2)
	assert.True(t, result)

	result = SliceAllGreaterThanOrEqualTo(s, 3)
	assert.False(t, result)
}

func TestSliceAnyValuesLessThan(t *testing.T) {
	s := []int{2, 3, 4, 5, 6}
	result := SliceAnyLessThan(s, 3)
	assert.True(t, result)

	result = SliceAnyLessThan(s, 2)
	assert.False(t, result)
}

func TestSliceAnyValuesLessThanOrEqualTo(t *testing.T) {
	s := []int{2, 3, 4, 5, 6}
	result := SliceAnyLessThanOrEqualTo(s, 2)
	assert.True(t, result)

	result = SliceAnyLessThanOrEqualTo(s, 1)
	assert.False(t, result)
}

func TestSliceAnyValuesGreaterThan(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	result := SliceAnyGreaterThan(s, 4)
	assert.True(t, result)

	result = SliceAnyGreaterThan(s, 5)
	assert.False(t, result)
}

func TestSliceAnyValuesGreaterThanOrEqualTo(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	result := SliceAnyGreaterThanOrEqualTo(s, 5)
	assert.True(t, result)

	result = SliceAnyGreaterThanOrEqualTo(s, 6)
	assert.False(t, result)
}

func TestSlicePropositionNone(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	result := SlicePropositionNone(s, func(e int) bool {
		return e > 5
	})
	assert.True(t, result)

	result = SlicePropositionNone(s, func(e int) bool {
		return e < 5
	})
	assert.False(t, result)
}

func TestSlicePropositionAny(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	result := SlicePropositionAny(s, func(e int) bool {
		return e > 5
	})
	assert.False(t, result)

	result = SlicePropositionAny(s, func(e int) bool {
		return e < 5
	})
	assert.True(t, result)
}
func TestSlicePropositionAll(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	result := SlicePropositionAll(s, func(e int) bool {
		return e < 6
	})
	assert.True(t, result)

	result = SlicePropositionAll(s, func(e int) bool {
		return e < 5
	})
	assert.False(t, result)
}
