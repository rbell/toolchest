/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package rankCalculation

// WithRanker allows setting the ranker to a ranker that implements the Ranker interface, allowing ranking algorithms outside of those supported by this package to be used.
func WithRanker[T comparable](ranker Ranker[T]) RankCalculatorOption[T] {
	return func(calculator *RankCalculator[T]) {
		calculator.ranker = ranker
	}
}

// WithRankPositionally allows setting the ranker to a ranker that ranks entries based on their position in the sorted list of entries.
func WithRankPositionally[T comparable]() RankCalculatorOption[T] {
	return func(calculator *RankCalculator[T]) {
		calculator.ranker = NewPercentileRanker[T](true)
	}
}
