/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package storage

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type FifoMapCache[K comparable, V any] struct {
	partitions          *GenericStack[*SafeMap[K, V]]
	currentPartitionId  uint64
	currentPartitionMux *sync.RWMutex
	valuePartitionIndex *SafeMap[K, uint64]
	ctx                 context.Context
	partitionCapacity   int
	maxPartitions       int
	capacity            int
	count               atomic.Int64
	sweepingWG          *sync.WaitGroup
}

type numPartitionCalculator func(capacity int) (int, int)

type fifoMapConfiguration struct {
	numPartitionCalculator numPartitionCalculator
	sweepFrequency         time.Duration
}

type fifoInitializationOption func(configuration *fifoMapConfiguration)

func NewFifoMapCache[K comparable, V any](ctx context.Context, capacity int, options ...fifoInitializationOption) *FifoMapCache[K, V] {
	// default config
	cfg := &fifoMapConfiguration{
		numPartitionCalculator: calcBalancedPartitions,
		sweepFrequency:         time.Second * 20,
	}
	for _, opt := range options {
		opt(cfg)
	}

	numPartitions, partitionLength := cfg.numPartitionCalculator(capacity)
	cache := &FifoMapCache[K, V]{
		partitions:          NewGenericStack[*SafeMap[K, V]](numPartitions),
		maxPartitions:       numPartitions,
		ctx:                 ctx,
		partitionCapacity:   partitionLength,
		valuePartitionIndex: NewSafeMap[K, uint64](0),
		currentPartitionMux: &sync.RWMutex{},
		sweepingWG:          &sync.WaitGroup{},
	}

	go func() {

		ticker := time.NewTicker(cfg.sweepFrequency)
		for {
			select {
			case <-ticker.C:
				cache.sweep()
			case <-cache.ctx.Done():
				return
			}
		}

	}()

	return cache
}

// Contains returns true if the key of type K is in the map
func (f *FifoMapCache[K, V]) Contains(key K) bool {
	if partitionId := f.valuePartitionIndex.Get(key); partitionId > 0 {
		partition, _ := f.partitions.Peek(partitionId)
		if partition != nil {
			return partition.Contains(key)
		}
	}
	return false
}

// Get returns the value of type V for the key of type K.  If the key is not found, the zero value of V is returned.
func (f *FifoMapCache[K, V]) Get(key K) (value V) {
	if partitionId := f.valuePartitionIndex.Get(key); partitionId > 0 {
		partition, _ := f.partitions.Peek(partitionId)
		if partition != nil {
			return partition.Get(key)
		}
	}
	return
}

// Set sets the value of type V for the key of type K.
func (f *FifoMapCache[K, V]) Set(key K, value V) {
	var partitionId uint64
	// if key exists, update value
	if partitionId = f.valuePartitionIndex.Get(key); partitionId > 0 {
		partition, _ := f.partitions.Peek(partitionId)
		if partition != nil {
			partition.Set(key, value)
			return
		}
	}

	partition, partitionId := f.getCurrentPartition()
	partition.Set(key, value)
	f.count.Add(1)
	f.valuePartitionIndex.Set(key, partitionId)
}

// getCurrentPartition returns reference to the currentPartition which new key/values should be added to
func (f *FifoMapCache[K, V]) getCurrentPartition() (*SafeMap[K, V], uint64) {
	f.currentPartitionMux.RLock()
	if currentPartition, _ := f.partitions.Peek(f.currentPartitionId); currentPartition != nil && currentPartition.Len() < f.partitionCapacity {
		defer f.currentPartitionMux.RUnlock()
		return currentPartition, f.currentPartitionId
	}
	f.currentPartitionMux.RUnlock()
	f.currentPartitionMux.Lock()
	defer f.currentPartitionMux.Unlock()
	newPartition := NewSafeMap[K, V](f.partitionCapacity)
	f.currentPartitionId = f.partitions.Push(newPartition)
	f.sweep()
	return newPartition, f.currentPartitionId
}

// sweep removes partitions from the stack if the number of partitions exceeds the maxPartitions
func (f *FifoMapCache[K, V]) sweep() {
	// restrict to single sweep at a time
	f.sweepingWG.Wait()
	f.sweepingWG.Add(1)
	defer f.sweepingWG.Done()

	if f.partitions.Len() > f.maxPartitions {
		numToPop := f.partitions.Len() - f.maxPartitions
		for i := 0; i < numToPop; i++ {
			popped := f.partitions.Pop()
			f.count.Add(-1 * int64(popped.Len()))
		}
	}
}

// default numPartitionCalculator, creates a balance between number of partitions and the size of each partition.
func calcBalancedPartitions(capacity int) (int, int) {
	numPartitions := int(math.Floor(math.Sqrt(float64(capacity))))
	partitionLenth := int(math.Floor(float64(capacity) / float64(numPartitions)))
	return numPartitions, partitionLenth
}

//region fifoMapCacheOptions

// WithBalancedPartitions balances the partitions based upon the nRoot parameter, calculating the number of partitions equal to the Nth root of capacity
// nRoot should be greater than 1.
// nRoot of 2 is same as the default number of partitions, which balances the number of values in each partition close to the number of partitions.
// nRoot less than 2 (greater than 1) will reduce number of partitions, making each partition containing more values.
// nRoot greater than 2 will increase number of partitions, making each partition contains fewer values.
// If the calculated number of partitions is less than minimumPartitions, minimumPartitions is used.
func WithBalancedPartitions(nRoot float64, minimumPartitions int) fifoInitializationOption {
	return func(configuration *fifoMapConfiguration) {
		configuration.numPartitionCalculator = func(capacity int) (int, int) {
			factor := 1 / nRoot

			rawNumPartitions := math.Pow(float64(capacity), factor)
			numPartitions := int(math.Floor(rawNumPartitions))
			if numPartitions < minimumPartitions {
				numPartitions = minimumPartitions
			}
			partitionLenth := int(math.Floor(float64(capacity) / float64(numPartitions)))
			return numPartitions, partitionLenth
		}
	}
}

//endregion
