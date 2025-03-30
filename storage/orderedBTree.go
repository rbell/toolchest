package storage

import (
	"cmp"
	"sync"

	"github.com/google/btree"
)

type orderedItem[K cmp.Ordered, V any] struct {
	key   K
	value *V
}

func (i orderedItem[K, V]) Less(than btree.Item) bool {
	return i.key < than.(orderedItem[K, V]).key
}

// OrderedBTree is a thread-safe implementation of a BTree that supports ordered types.  It is a wrapper around the google/btree package.
type OrderedBTree[K cmp.Ordered, V any] struct {
	*btree.BTree
	writeMux *sync.Mutex
}

// NewOrderedBTree returns an initialized reference to an OrderedBTree of K and V
func NewOrderedBTree[K cmp.Ordered, V any]() *OrderedBTree[K, V] {
	return &OrderedBTree[K, V]{
		BTree:    btree.New(2),
		writeMux: &sync.Mutex{},
	}
}

// Set sets the value of type V for the key of type K.
func (t *OrderedBTree[K, V]) Set(key K, value *V) {
	t.writeMux.Lock()
	defer t.writeMux.Unlock()
	orderedKey := orderedItem[K, V]{key, value}
	t.ReplaceOrInsert(orderedKey)
}

// Get returns the value of type V for the key of type K.  If the key is not found, ok is returned as false.
func (t *OrderedBTree[K, V]) Get(key K) (value *V, ok bool) {
	orderedKey := orderedItem[K, V]{key, nil}
	item := t.BTree.Get(orderedKey)
	if item == nil {
		return
	}
	return item.(orderedItem[K, V]).value, true
}

// Delete deletes the key of type K.  If the key is not found, ok is returned as false.
func (t *OrderedBTree[K, V]) Delete(key K) (value *V, ok bool) {
	t.writeMux.Lock()
	defer t.writeMux.Unlock()
	orderedKey := orderedItem[K, V]{key, nil}
	deleted := t.BTree.Delete(orderedKey)
	if deleted != nil {
		return deleted.(orderedItem[K, V]).value, true
	}
	return nil, false
}

// Has returns true if the key of type K exists in the BTree.
func (t *OrderedBTree[K, V]) Has(key K) bool {
	orderedKey := orderedItem[K, V]{key, nil}
	return t.BTree.Has(orderedKey)
}

// Len returns the length of the BTree
func (t *OrderedBTree[K, V]) Len() int {
	return t.BTree.Len()
}

// Min returns the minimum key and value in the BTree
func (t *OrderedBTree[K, V]) Min() (key K, value *V) {
	item := t.BTree.Min()
	if item == nil {
		return
	}
	return item.(orderedItem[K, V]).key, item.(orderedItem[K, V]).value
}

// Max returns the maximum key and value in the BTree
func (t *OrderedBTree[K, V]) Max() (key K, value *V) {
	item := t.BTree.Max()
	if item == nil {
		return
	}
	return item.(orderedItem[K, V]).key, item.(orderedItem[K, V]).value
}

// Ascend calls the iter function for every key/value pair in the BTree in ascending order.
func (t *OrderedBTree[K, V]) Ascend(iter func(key K, value *V) bool) {
	t.BTree.Ascend(func(item btree.Item) bool {
		return iter(item.(orderedItem[K, V]).key, item.(orderedItem[K, V]).value)
	})
}

// AscendGreaterOrEqual calls the iter function for every key/value pair in the BTree in ascending order starting with the first key/value pair that is greater than or equal to the pivot key.
func (t *OrderedBTree[K, V]) AscendGreaterOrEqual(pivot K, iter func(key K, value *V) bool) {
	greaterThanEqualKey := orderedItem[K, V]{pivot, nil}
	t.BTree.AscendGreaterOrEqual(greaterThanEqualKey, func(item btree.Item) bool {
		return iter(item.(orderedItem[K, V]).key, item.(orderedItem[K, V]).value)
	})
}

// AscendGreaterThan calls the iter function for every key/value pair in the BTree, starting with the Min value up to the pivot key..
func (t *OrderedBTree[K, V]) AscendLessThan(pivot K, iter func(key K, value *V) bool) {
	lessThanKey := orderedItem[K, V]{pivot, nil}
	t.BTree.AscendLessThan(lessThanKey, func(item btree.Item) bool {
		return iter(item.(orderedItem[K, V]).key, item.(orderedItem[K, V]).value)
	})
}

// AscendRange calls the iter function for every key/value pair in the BTree, starting with the first key/value pair that is greater than or equal to the greaterThanEqual key up to the first key/value pair that is less than the lessThan key.
func (t *OrderedBTree[K, V]) AscendRange(greaterThanEqual, lessThan K, iter func(key K, value *V) bool) {
	greaterThanEqualKey := orderedItem[K, V]{greaterThanEqual, nil}
	lessThanKey := orderedItem[K, V]{lessThan, nil}
	t.BTree.AscendRange(greaterThanEqualKey, lessThanKey, func(item btree.Item) bool {
		return iter(item.(orderedItem[K, V]).key, item.(orderedItem[K, V]).value)
	})
}

// Descend calls the iter function for every key/value pair in the BTree in descending order.
func (t *OrderedBTree[K, V]) Descend(iter func(key K, value *V) bool) {
	t.BTree.Descend(func(item btree.Item) bool {
		return iter(item.(orderedItem[K, V]).key, item.(orderedItem[K, V]).value)
	})
}

// DDescendGreaterThan calls the iterator for every value in the tree within the range [last, pivot), until iterator returns false
func (t *OrderedBTree[K, V]) DescendLessOrEqual(pivot K, iter func(key K, value *V) bool) {
	lessThanEqualKey := orderedItem[K, V]{pivot, nil}
	t.BTree.DescendLessOrEqual(lessThanEqualKey, func(item btree.Item) bool {
		return iter(item.(orderedItem[K, V]).key, item.(orderedItem[K, V]).value)
	})
}

// DescendLessOrEqual calls the iterator for every value in the tree within the range [pivot, first], until iterator returns false
func (t *OrderedBTree[K, V]) DescendGreaterThan(pivot K, iter func(key K, value *V) bool) {
	greaterThanKey := orderedItem[K, V]{pivot, nil}
	t.BTree.DescendGreaterThan(greaterThanKey, func(item btree.Item) bool {
		return iter(item.(orderedItem[K, V]).key, item.(orderedItem[K, V]).value)
	})
}

// DescendRange calls the iterator for every value in the tree within the range [lessOrEqual, greaterThan), until iterator returns false
func (t *OrderedBTree[K, V]) DescendRange(greaterThan, lessThanEqual K, iter func(key K, value *V) bool) {
	greaterThanKey := orderedItem[K, V]{greaterThan, nil}
	lessThanEqualKey := orderedItem[K, V]{lessThanEqual, nil}
	t.BTree.DescendRange(greaterThanKey, lessThanEqualKey, func(item btree.Item) bool {
		return iter(item.(orderedItem[K, V]).key, item.(orderedItem[K, V]).value)
	})
}

// DeleteMin deletes the minimum key and value in the BTree
func (t *OrderedBTree[K, V]) DeleteMin() (key K, value *V) {
	item := t.BTree.DeleteMin()
	if item == nil {
		return
	}
	return item.(orderedItem[K, V]).key, item.(orderedItem[K, V]).value
}

// DeleteMax deletes the maximum key and value in the BTree
func (t *OrderedBTree[K, V]) DeleteMax() (key K, value *V) {
	item := t.BTree.DeleteMax()
	if item == nil {
		return
	}
	return item.(orderedItem[K, V]).key, item.(orderedItem[K, V]).value
}

// Clone returns a copy of the BTree
func (t *OrderedBTree[K, V]) Clone() *OrderedBTree[K, V] {
	return &OrderedBTree[K, V]{
		BTree:    t.BTree.Clone(),
		writeMux: &sync.Mutex{},
	}
}
