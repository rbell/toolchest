package MapOps

import (
	"cmp"
	"sort"
)

type pair[K comparable, V cmp.Ordered] struct {
	key   K
	value V
}

type pairList[K comparable, V cmp.Ordered] []pair[K, V]

func (p pairList[K, V]) Len() int           { return len(p) }
func (p pairList[K, V]) Less(i, j int) bool { return p[i].value < p[j].value }
func (p pairList[K, V]) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// SortedAscKeys returns a sorted slice of keys for the given map
func SortAscKeys[K comparable, V cmp.Ordered](m map[K]V) []K {
	pairs := make(pairList[K, V], 0, len(m))
	for k, v := range m {
		pairs = append(pairs, pair[K, V]{k, v})
	}
	sort.Sort(pairs)
	sortedKeys := make([]K, 0, len(pairs))
	for _, p := range pairs {
		sortedKeys = append(sortedKeys, p.key)
	}
	return sortedKeys
}

// SortDescKeys returns a sorted slice of keys for the given map in descending order
func SortDescKeys[K comparable, V cmp.Ordered](m map[K]V) []K {
	pairs := make(pairList[K, V], 0, len(m))
	for k, v := range m {
		pairs = append(pairs, pair[K, V]{k, v})
	}
	sort.Sort(sort.Reverse(pairs))
	sortedKeys := make([]K, 0, len(pairs))
	for _, p := range pairs {
		sortedKeys = append(sortedKeys, p.key)
	}
	return sortedKeys
}
