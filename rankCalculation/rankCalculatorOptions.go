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
