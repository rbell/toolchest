package rankCalculation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPercentileRanker_CreatesRanker(t *testing.T) {
	// test
	ranker := NewPercentileRanker[int](true)

	// assert
	assert.NotNil(t, ranker)
	assert.Equal(t, true, ranker.positionalRanking)
}

func TestRank_PositionalRankTrue_RanksEntries(t *testing.T) {
	// setup
	ranker := NewPercentileRanker[int](true)
	entries := map[int]int64{100: 10, 200: 20, 300: 30}
	expected := map[int]float64{100: 0, 200: 50, 300: 100}

	// test
	result, err := ranker.Rank(entries)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestRank_PositionalRankFalse_RanksEntries(t *testing.T) {
	// setup
	ranker := NewPercentileRanker[int](false)
	entries := map[int]int64{1: 10, 2: 20, 3: 30}
	expected := map[int]float64{1: 33.33333333333333, 2: 66.66666666666666, 3: 100}

	// test
	result, err := ranker.Rank(entries)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestRank_NoEntries_ReturnsEmptyMap(t *testing.T) {
	// setup
	ranker := NewPercentileRanker[int](false)
	entries := map[int]int64{}
	expected := map[int]float64{}

	// test
	result, err := ranker.Rank(entries)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestRank_OneEntry_ReturnsRankOneHundred(t *testing.T) {
	// setup
	ranker := NewPercentileRanker[int](false)
	entries := map[int]int64{1: 10}
	expected := map[int]float64{1: 100}

	// test
	result, err := ranker.Rank(entries)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestRank_TwoEntries_ReturnsRanksFiftyAndOneHundred(t *testing.T) {
	// setup
	ranker := NewPercentileRanker[int](false)
	entries := map[int]int64{1: 10, 2: 20}
	expected := map[int]float64{1: 50, 2: 100}

	// test
	result, err := ranker.Rank(entries)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestRank_ThreeEntries_ReturnsRanksThirtyThreeSixtySixAndOneHundred(t *testing.T) {
	// setup
	ranker := NewPercentileRanker[int](false)
	entries := map[int]int64{1: 10, 2: 20, 3: 30}
	expected := map[int]float64{1: 33.33333333333333, 2: 66.66666666666666, 3: 100}

	// test
	result, err := ranker.Rank(entries)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestRank_ThreeEntriesWithZero_ReturnsRanksZeroFiftyAndOneHundred(t *testing.T) {
	// setup
	ranker := NewPercentileRanker[int](false)
	entries := map[int]int64{0: 0, 1: 10, 2: 20}
	expected := map[int]float64{0: 0, 1: 50, 2: 100}

	// test
	result, err := ranker.Rank(entries)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}
