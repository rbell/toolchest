/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package storage

import "fmt"

type childAdderGetter[T comparable] interface {
	Get() T
	AddChild(child childAdderGetter[T])
	GetChildren() []childAdderGetter[T]
}

type nodeFactory[T comparable] func(T) childAdderGetter[T]

type treeOption[T comparable] func(*Tree[T])

type simpleNode[T comparable] struct {
	value    T
	children []childAdderGetter[T]
}

func (t *simpleNode[T]) Get() T {
	return t.value
}

func (t *simpleNode[T]) AddChild(child childAdderGetter[T]) {
	t.children = append(t.children, child)
}

func (t *simpleNode[T]) GetChildren() []childAdderGetter[T] {
	return t.children
}

type Tree[T comparable] struct {
	root        childAdderGetter[T]
	nodeFactory nodeFactory[T]
}

func NewTree[T comparable](opts ...treeOption[T]) *Tree[T] {
	t := &Tree[T]{
		root: nil,
		nodeFactory: func(v T) childAdderGetter[T] {
			return &simpleNode[T]{value: v}
		},
	}

	for _, o := range opts {
		o(t)
	}

	return t
}

func WithNodeFactory[T comparable](factory nodeFactory[T]) treeOption[T] {
	return func(t *Tree[T]) {
		t.nodeFactory = factory
	}
}

func (t *Tree[T]) Walk(f func(T, int)) {
	t.walk(t.root, 0, f)
}

func (t *Tree[T]) walk(node childAdderGetter[T], ancestryLevel int, f func(T, int)) {
	f(node.Get(), ancestryLevel)
	for _, child := range node.GetChildren() {
		t.walk(child, ancestryLevel+1, f)
	}
}

func AddAncestryChain[T comparable](tree *Tree[T], ancestry ...T) error {
	lowestParent, missingAncestry := getLowestMatchingLeaf(tree.root, ancestry...)

	if lowestParent == nil && tree.root != nil {
		return fmt.Errorf("ancestry cannot be added because the root node does not match the first ancestor")
	}

	if len(missingAncestry) > 0 {
		currentNode := lowestParent
		for i, ancestor := range missingAncestry {
			currentNode = addChild(tree.nodeFactory, currentNode, ancestor)
			// if this is the first ancestor in the list and tree is empty, set it as the root
			if tree.root == nil && i == 0 {
				tree.root = currentNode
			}
		}
	}

	return nil
}

func getLowestMatchingLeaf[T comparable](node childAdderGetter[T], ancestry ...T) (existingAncestor childAdderGetter[T], missingAncestry []T) {
	if len(ancestry) == 0 {
		return nil, nil // no ancestry remaining to match
	}

	currentAncestor := ancestry[0]
	var descendents []T
	if len(ancestry) > 1 {
		descendents = ancestry[1:]
	}

	if node != nil && node.Get() == currentAncestor {
		if len(descendents) == 0 {
			return node, descendents
		}
		for _, child := range node.GetChildren() {
			// if the child matches the next ancestor in the list, descend into it
			if leaf, missingDescendents := getLowestMatchingLeaf(child, descendents...); leaf != nil {
				return leaf, missingDescendents
			}
		}
		return node, descendents
	}

	return nil, ancestry
}

func addChild[T comparable](factory nodeFactory[T], parent childAdderGetter[T], value T) childAdderGetter[T] {
	newNode := factory(value)
	if parent != nil {
		parent.AddChild(newNode)
	}
	return newNode
}
