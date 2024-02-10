package rankCalculation

import "github.com/rbell/toolchest/MapOps"

type PercentileRanker[T comparable] struct {
	positionalRanking bool
}

// NewPercentileRanker returns an initialized reference to a PercentileRanker
// positionalRank is true if the rank is based on the position of the entry in the sorted list, otherwise the rank is based on the value of the entry
func NewPercentileRanker[T comparable](positionalRank bool) *PercentileRanker[T] {
	return &PercentileRanker[T]{positionalRanking: positionalRank}
}

func (r *PercentileRanker[T]) Rank(entries map[T]uint64) map[T]float64 {
	percentiles := make(map[T]float64)
	total := uint64(0)
	for _, v := range entries {
		total += v
	}
	// sort entries ascending
	sortedKeys := MapOps.SortAscKeys(entries)

	if r.positionalRanking {
		// calculate the percentile for each entry based on position
		for i, k := range sortedKeys {
			percentiles[k] = float64(i) / float64(len(sortedKeys)) * 100
		}
		return percentiles
	}

	// calculate the percentile for the value of each entry
	for _, k := range sortedKeys {
		percentiles[k] = float64(entries[k]) / float64(total) * 100
	}
	return percentiles
}
