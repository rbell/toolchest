/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package storage

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type testNode struct {
	value    string
	children []childAdderGetter[string]
}

func (t *testNode) Get() string {
	return t.value
}

func (t *testNode) AddChild(child childAdderGetter[string]) {
	t.children = append(t.children, child)
}

func (t *testNode) GetChildren() []childAdderGetter[string] {
	return t.children
}

var testFactory = func(s string) childAdderGetter[string] {
	return &testNode{value: s}
}

var testRoot = &testNode{
	value: "top",
	children: []childAdderGetter[string]{
		&testNode{
			value: "child1",
			children: []childAdderGetter[string]{
				&testNode{
					value: "child1.1",
				},
				&testNode{
					value: "child1.2",
				},
			},
		},
		&testNode{
			value: "child2",
			children: []childAdderGetter[string]{
				&testNode{
					value: "child2.1",
				},
				&testNode{
					value: "child2.2",
				},
			},
		},
	},
}

var testTree = &Tree[string]{
	root:        testRoot,
	nodeFactory: testFactory,
}

func TestTree_GetLowestMatchingLeaf_MatchingTopAncestry(t *testing.T) {
	// setup
	ancestry := []string{"top"}

	// test
	lowest, missingAncestry := getLowestMatchingLeaf(testRoot, ancestry...)

	// assert
	assert.Equal(t, "top", lowest.Get())
	assert.Empty(t, missingAncestry)
}

func TestTree_GetLowestMatchingLeaf_MatchingChild1(t *testing.T) {
	// setup
	ancestry := []string{"top", "child1"}

	// test
	lowest, missingAncestry := getLowestMatchingLeaf(testRoot, ancestry...)

	// assert
	assert.Equal(t, "child1", lowest.Get())
	assert.Empty(t, missingAncestry)
}

func TestTree_GetLowestMatchingLeaf_MatchingChild2(t *testing.T) {
	// setup
	ancestry := []string{"top", "child2"}

	// test
	lowest, missingAncestry := getLowestMatchingLeaf(testRoot, ancestry...)

	// assert
	assert.Equal(t, "child2", lowest.Get())
	assert.Empty(t, missingAncestry)
}

func TestTree_GetLowestMatchingLeaf_MatchingChild1_2(t *testing.T) {
	// setup
	ancestry := []string{"top", "child1", "child1.2"}

	// test
	lowest, missingAncestry := getLowestMatchingLeaf(testRoot, ancestry...)

	// assert
	assert.Equal(t, "child1.2", lowest.Get())
	assert.Empty(t, missingAncestry)
}

func TestTree_GetLowestMatchingLeaf_TailofAncestryNotInTreee_ReturnsLowestMatchingAncestor(t *testing.T) {
	// setup
	ancestry := []string{"top", "child1", "child1.3"} // child1.3 does not exist as child of child1

	// test
	lowest, missingAncestry := getLowestMatchingLeaf(testRoot, ancestry...)

	// assert
	assert.Equal(t, "child1", lowest.Get())
	assert.Contains(t, missingAncestry, "child1.3")
}

func TestTree_GetLowestMatchingLeaf_AncestryChainDoesNotExist_ReturnsNil(t *testing.T) {
	// setup
	ancestry := []string{"NotTop", "child1", "child1.3"} // NotTop is not the root of the testRoot

	// test
	lowest, missingAncestry := getLowestMatchingLeaf(testRoot, ancestry...)

	// assert
	assert.Nil(t, lowest)
	assert.Contains(t, missingAncestry, "NotTop")
	assert.Contains(t, missingAncestry, "child1")
	assert.Contains(t, missingAncestry, "child1.3")
}

func TestTree_GetLowestMatchingLeaf_EmptyAncestry_ReturnsNil(t *testing.T) {
	// setup
	ancestry := []string{}

	// test
	lowest, missingAncestry := getLowestMatchingLeaf(testRoot, ancestry...)

	// assert
	assert.Nil(t, lowest)
	assert.Nil(t, missingAncestry)
}

func TestTree_AddAncestryChain_EmptyTree_ReturnsNil(t *testing.T) {
	// setup
	ancestry := []string{"top", "child1", "child1.1"}
	tree := NewTree[string]()

	// test
	err := AddAncestryChain(tree, ancestry...)

	// assert
	assert.Nil(t, err)

}

func TestTree_AddAncestryChain_AncestryChainDoesNotExist_ReturnsError(t *testing.T) {
	// setup
	ancestry := []string{"NotTop", "child1", "child1.3"}

	// test
	err := AddAncestryChain(testTree, ancestry...)

	// assert
	assert.Error(t, err)
}

func TestTree_AddAncestryChain_AncestryChainExists_ReturnsNil(t *testing.T) {
	// setup
	ancestry := []string{"top", "child1", "child1.1"}
	tree := NewTree[string]()

	// test
	err := AddAncestryChain(tree, ancestry...)

	// assert
	assert.Nil(t, err)
}

func TestTree_AddAncestryChain_MiltipleAncestryChains_ReturnsNil(t *testing.T) {
	// setup
	ancestries := [][]string{
		{"top", "child1", "child1.1"},
		{"top", "child1", "child1.2"},
		{"top", "child2", "child2.1"},
		{"top", "child2", "child2.2"},
	}
	tree := NewTree[string]()

	// test
	for _, ancestry := range ancestries {
		err := AddAncestryChain(tree, ancestry...)
		assert.Nil(t, err)
	}
}

func TestTree_Walk(t *testing.T) {
	// setup
	tree := NewTree[string]()
	ancestries := [][]string{
		{"top", "child1", "child1.1"},
		{"top", "child1", "child1.2"},
		{"top", "child2", "child2.1"},
		{"top", "child2", "child2.2"},
	}
	for _, ancestry := range ancestries {
		//nolint:errcheck // ignore errorlint error for test
		AddAncestryChain(tree, ancestry...)
	}

	// test
	var values []string
	tree.Walk(func(s string, level int) {
		values = append(values, s)
	})

	// assert
	assert.ElementsMatch(t, []string{"top", "child1", "child1.1", "child1.2", "child2", "child2.1", "child2.2"}, values)
}

func TestTree_Walk_UsingLevels(t *testing.T) {
	// setup
	expected := `top
 -> child1
 ->  -> child1.1
 ->  -> child1.2
 -> child2
 ->  -> child2.1
 ->  -> child2.2
`

	tree := NewTree[string]()
	ancestries := [][]string{
		{"top", "child1", "child1.1"},
		{"top", "child1", "child1.2"},
		{"top", "child2", "child2.1"},
		{"top", "child2", "child2.2"},
	}
	for _, ancestry := range ancestries {
		//nolint:errcheck // ignore errorlint error for test
		AddAncestryChain(tree, ancestry...)
	}

	// test
	resultSB := strings.Builder{}
	seperator := " -> "
	tree.Walk(func(s string, level int) {
		indent := strings.Repeat(seperator, level)
		resultSB.WriteString(fmt.Sprintf("%s%s\n", indent, s))
	})

	// assert
	assert.Equal(t, expected, resultSB.String())
}
