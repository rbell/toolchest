/*
 * Copyright (c) 2024  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package sliceOps

import (
	"github.com/rbell/toolchest/propositions"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCutPreserveOrder(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	expectedResult := []int{2, 3}
	expectedRemaining := []int{1, 4, 5}
	result := Cut(&s, 1, 3)
	assert.True(t, propositions.SliceContainsAll(s, expectedRemaining))
	assert.True(t, propositions.SliceContainsAll(result, expectedResult))
}

func TestRemove(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	expectedS := []int{1, 4, 5}
	Remove(&s, 1, 3)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestFilterInPlace(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	expectedS := []int{2, 4}
	FilterInPlace(&s, func(e int) bool {
		return e%2 == 0
	})
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestFilterInPlace_WorksWithReferences(t *testing.T) {
	t1 := 1
	t2 := 2
	t3 := 3
	t4 := 4
	t5 := 5
	s := []*int{&t1, &t2, &t3, &t4, &t5}
	expectedS := []*int{&t2, &t4}
	FilterInPlace(&s, func(e *int) bool {
		return *e%2 == 0
	})
	if !reflect.DeepEqual(s, expectedS) {
		t.Errorf("Expected %v, got %v", expectedS, s)
	}
}

func TestInsert(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	expectedS := []int{1, 2, 6, 7, 3, 4, 5}
	s = Insert(s, 2, 6, 7)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestPush(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	expectedS := []int{6, 7, 1, 2, 3, 4, 5}
	Push(&s, 6, 7)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestPop(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	expectedS := []int{2, 3, 4, 5}
	result := Pop(&s)
	if result != 1 {
		t.Errorf("Expected %v, got %v", 1, result)
	}
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestPop_EmptySlice(t *testing.T) {
	s := []int{}
	result := Pop(&s)
	if result != 0 {
		t.Errorf("Expected %v, got %v", 0, result)
	}
	assert.True(t, propositions.SliceContainsAll(s, []int{}))
}

func TestPop_OneElementSlice(t *testing.T) {
	s := []int{1}
	result := Pop(&s)
	if result != 1 {
		t.Errorf("Expected %v, got %v", 1, result)
	}
	assert.True(t, propositions.SliceContainsAll(s, []int{}))
}

func TestPop_TwoElementSlice(t *testing.T) {
	s := []int{1, 2}
	result := Pop(&s)
	if result != 1 {
		t.Errorf("Expected %v, got %v", 1, result)
	}
	assert.True(t, propositions.SliceContainsAll(s, []int{2}))
}

func TestDistinct(t *testing.T) {
	s := []int{1, 1, 2, 3, 3, 4, 5, 5}
	expectedS := []int{1, 2, 3, 4, 5}
	s = Distinct(s)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestDistinct_EmptySlice(t *testing.T) {
	s := []int{}
	expectedS := []int{}
	s = Distinct(s)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestDistinct_OneElementSlice(t *testing.T) {
	s := []int{1}
	expectedS := []int{1}
	s = Distinct(s)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestUnion(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := []int{3, 4, 5, 6, 7}
	expectedS := []int{1, 2, 3, 4, 5, 6, 7}
	s := Union(s1, s2)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestUnion_EmptyS1(t *testing.T) {
	s1 := []int{}
	s2 := []int{3, 4, 5, 6, 7}
	expectedS := []int{3, 4, 5, 6, 7}
	s := Union(s1, s2)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestUnion_EmptyS2(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := []int{}
	expectedS := []int{1, 2, 3, 4, 5}
	s := Union(s1, s2)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestUnion_EmptyS1S2(t *testing.T) {
	s1 := []int{}
	s2 := []int{}
	expectedS := []int{}
	s := Union(s1, s2)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestUnion_ThreeSlices(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := []int{3, 4, 5, 6, 7}
	s3 := []int{5, 6, 7, 8, 9}
	expectedS := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	s := Union(s1, s2, s3)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestIntersection(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := []int{3, 4, 5, 6, 7}
	expectedS := []int{3, 4, 5}
	s := Intersection(s1, s2)

	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestIntersection_EmptyS1(t *testing.T) {
	s1 := []int{}
	s2 := []int{3, 4, 5, 6, 7}
	expectedS := []int{}
	s := Intersection(s1, s2)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestIntersection_EmptyS2(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := []int{}
	expectedS := []int{}
	s := Intersection(s1, s2)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestIntersection_EmptyS1S2(t *testing.T) {
	s1 := []int{}
	s2 := []int{}
	expectedS := []int{}
	s := Intersection(s1, s2)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestInstersection_ThreeSlices(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := []int{3, 4, 5, 6, 7}
	s3 := []int{5, 6, 7, 8, 9}
	expectedS := []int{5}
	s := Intersection(s1, s2, s3)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestDifference(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := []int{3, 4, 5, 6, 7}
	expectedS := []int{1, 2}
	s := Difference(s1, s2)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestDifference_EmptyS1(t *testing.T) {
	s1 := []int{}
	s2 := []int{3, 4, 5, 6, 7}
	expectedS := []int{}
	s := Difference(s1, s2)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestDifference_EmptyS2(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := []int{}
	expectedS := []int{1, 2, 3, 4, 5}
	s := Difference(s1, s2)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestDifference_EmptyS1S2(t *testing.T) {
	s1 := []int{}
	s2 := []int{}
	expectedS := []int{}
	s := Difference(s1, s2)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestDisjoin(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := []int{3, 4, 5, 6, 7}
	expectedS := []int{1, 2, 6, 7}
	s := Disjoin(s1, s2)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestDisjoin_EmptyS1(t *testing.T) {
	s1 := []int{}
	s2 := []int{3, 4, 5, 6, 7}
	expectedS := []int{3, 4, 5, 6, 7}
	s := Disjoin(s1, s2)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestDisjoin_EmptyS2(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := []int{}
	expectedS := []int{1, 2, 3, 4, 5}
	s := Disjoin(s1, s2)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestDisjoin_EmptyS1S2(t *testing.T) {
	s1 := []int{}
	s2 := []int{}
	expectedS := []int{}
	s := Disjoin(s1, s2)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}

func TestDisjoin_ThreeSlices(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5}
	s2 := []int{3, 4, 5, 6, 7}
	s3 := []int{5, 6, 7, 8, 9}
	expectedS := []int{1, 2, 8, 9}
	s := Disjoin(s1, s2, s3)
	assert.True(t, propositions.SliceContainsAll(s, expectedS))
}
