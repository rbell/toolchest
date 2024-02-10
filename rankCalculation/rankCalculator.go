package rankCalculation

import (
	"github.com/rbell/toolchest/storage"
	"sync"
	"sync/atomic"
)

type RankCalculatorOption[T comparable] func(calculator *RankCalculator[T])

type Ranker[T comparable] interface {
	Rank(entries map[T]uint64) map[T]float64 // calculates and returns the ranking of the entries
}

// RankCalculator is a thread-safe implementation of a rank calculator
type RankCalculator[T comparable] struct {
	entries *storage.SafeMap[T, *atomic.Uint64] // map of entries to their number of hits
	ranker  Ranker[T]                           // the ranker to use
	mux     *sync.RWMutex
}

// NewRankCalculator returns an initialized reference to a RankCalculator of T
func NewRankCalculator[T comparable](options ...RankCalculatorOption[T]) *RankCalculator[T] {
	return &RankCalculator[T]{
		entries: storage.NewSafeMap[T, *atomic.Uint64](0),
		ranker:  NewPercentileRanker[T](false),
		mux:     &sync.RWMutex{},
	}
}

// Accumulate adds the value of type T to the rank calculator
func (r *RankCalculator[T]) Accumulate(entry T) {
	r.entries.GetOrAdd(entry, &atomic.Uint64{}).Add(1)
}

// Reset clears the rank calculator
func (r *RankCalculator[T]) Reset() {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.entries = storage.NewSafeMap[T, *atomic.Uint64](0)
}

// Calculate returns the ranking of the entries
func (r *RankCalculator[T]) Calculate() map[T]float64 {
	r.mux.RLock()
	entryCpy := storage.TranslateToMapOf[T, *atomic.Uint64, uint64](r.entries, func(v *atomic.Uint64) uint64 {
		return v.Load()
	})
	return r.ranker.Rank(entryCpy)
}