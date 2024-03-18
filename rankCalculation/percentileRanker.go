/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package rankCalculation

import "github.com/rbell/toolchest/mapOps"

type PercentileRanker[T comparable] struct {
	positionalRanking bool
}

// NewPercentileRanker returns an initialized reference to a PercentileRanker
// positionalRank is true if the rank is based on the position of the entry in the sorted list, otherwise the rank is based on the value of the entry
func NewPercentileRanker[T comparable](positionalRank bool) *PercentileRanker[T] {
	return &PercentileRanker[T]{positionalRanking: positionalRank}
}

func (r *PercentileRanker[T]) Rank(entries map[T]int64) (map[T]float64, error) {
	percentiles := make(map[T]float64)
	if len(entries) == 0 {
		return percentiles, nil
	}

	// sort entries ascending
	sortedKeys := mapOps.SortAscKeys(entries)

	if r.positionalRanking {
		// calculate the percentile for each entry based on position
		for i, k := range sortedKeys {
			if i == 0 {
				percentiles[k] = 0
				continue
			}
			if i == len(sortedKeys)-1 {
				percentiles[k] = 100
				continue
			}
			percentiles[k] = float64(i+1) / float64(len(sortedKeys)+1) * 100
		}
		return percentiles, nil
	}

	// calculate the percentile for the value of each entry
	maxV := int64(0)
	for _, v := range entries {
		if v > maxV {
			maxV = v
		}
	}
	for _, k := range sortedKeys {
		percentiles[k] = float64(entries[k]) / float64(maxV) * 100
	}
	return percentiles, nil
}
