/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package mapOps

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortAscKeys_sortsAscending(t *testing.T) {
	// setup
	m := map[int]int{3: 3, 2: 2, 1: 1}
	expected := []int{1, 2, 3}

	// test
	result := SortAscKeys(m)

	// assert
	assert.Equal(t, expected, result, "Expected keys to be sorted in ascending order")
}

func TestSortDescKeys_sortsDescending(t *testing.T) {
	// setup
	m := map[int]int{1: 1, 2: 2, 3: 3}
	expected := []int{3, 2, 1}

	// test
	result := SortDescKeys(m)

	// assert
	assert.Equal(t, expected, result, "Expected keys to be sorted in descending order")
}

func TestSortAscKeys_emptyMap_returnsEmptySlice(t *testing.T) {
	// setup
	m := map[int]int{}

	// test
	result := SortAscKeys(m)

	// assert
	assert.Empty(t, result, "Expected empty slice")
}

func TestSortDescKeys_emptyMap_returnsEmptySlice(t *testing.T) {
	// setup
	m := map[int]int{}

	// test
	result := SortDescKeys(m)

	// assert
	assert.Empty(t, result, "Expected empty slice")
}

func TestSortAscKeys_singleElementMap_returnsSingleElementSlice(t *testing.T) {
	// setup
	m := map[int]int{1: 1}
	expected := []int{1}

	// test
	result := SortAscKeys(m)

	// assert
	assert.Equal(t, expected, result, "Expected single element slice")
}

func TestSortDescKeys_singleElementMap_returnsSingleElementSlice(t *testing.T) {
	// setup
	m := map[int]int{1: 1}
	expected := []int{1}

	// test
	result := SortDescKeys(m)

	// assert
	assert.Equal(t, expected, result, "Expected single element slice")
}
