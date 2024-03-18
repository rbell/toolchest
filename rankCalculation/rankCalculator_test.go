/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package rankCalculation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRankCalculator(t *testing.T) {
	// setup

	// test
	calculator := NewRankCalculator[int]()

	// assert
	assert.NotNil(t, calculator)
	assert.IsType(t, &PercentileRanker[int]{}, calculator.ranker)
}

func TestAccumulate(t *testing.T) {
	// setup
	calculator := NewRankCalculator[int]()
	entry := 1

	// test
	calculator.Accumulate(entry)

	// assert
	assert.Equal(t, int64(1), calculator.entries.Get(entry).Load())
}

func TestReset(t *testing.T) {
	// setup
	calculator := NewRankCalculator[int]()
	entry := 1
	calculator.Accumulate(entry)

	// test
	calculator.Reset()

	// assert
	assert.Nil(t, calculator.entries.Get(entry))
}

func TestCalculate(t *testing.T) {
	// setup
	calculator := NewRankCalculator[int]()
	entries := []int{1, 2, 3}
	for _, entry := range entries {
		for i := 0; i < entry; i++ {
			calculator.Accumulate(entry)
		}
	}

	// test
	ranks, err := calculator.Calculate()

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, ranks)
	assert.Equal(t, float64(33.33333333333333), ranks[1])
	assert.Equal(t, float64(66.66666666666666), ranks[2])
	assert.Equal(t, float64(100), ranks[3])
}
