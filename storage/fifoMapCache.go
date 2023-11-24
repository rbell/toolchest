/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package storage

import (
	"math"
	"sync/atomic"
	"time"
)

type FifoMapCache[K comparable, V any] struct {
	partitions          *GenericStack[map[K]V]
	valuePartitionIndex map[K]uint64
	stopCh              chan struct{}
	partitionLength     int
	capacity            int
}

type numPartitionCalculator func(capacity int) (int, int)

type fifoMapConfiguration struct {
	numPartitionCalculator numPartitionCalculator
	sweepFrequency         time.Duration
}

type fifoInitializationOption func(configuration fifoMapConfiguration)

func NewFifoMapCache[K comparable, V any](capacity int, options ...fifoInitializationOption) *FifoMapCache[K, V] {
	// default config
	cfg := fifoMapConfiguration{
		numPartitionCalculator: calcBalancedPartitions,
		sweepFrequency:         time.Second * 20,
	}
	for _, opt := range options {
		opt(cfg)
	}

	numPartitions, partitionLength := cfg.numPartitionCalculator(capacity)
	cache := &FifoMapCache[K, V]{
		partitions:      NewGenericStack[map[K]V](numPartitions),
		stopCh:          make(chan struct{}),
		partitionLength: partitionLength,
	}

	go func() {
		sweeping := atomic.Bool{}
		ticker := time.NewTicker(cfg.sweepFrequency)
		for {
			select {
			case <-ticker.C:
				if !sweeping.Load() {
					sweeping.Store(true)

					sweeping.Store(false)
				}
			case <-cache.stopCh:
				return
			}
		}

	}()

}

// default numPartitionCalculator, creates a balance between number of partitions and the size of each partition.
func calcBalancedPartitions(capacity int) (int, int) {
	numPartitions := int(math.Ceil(math.Sqrt(float64(capacity))))
	partitionLenth := int(math.Ceil(float64(capacity) / float64(numPartitions)))
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
	return func(configuration fifoMapConfiguration) {
		configuration.numPartitionCalculator = func(capacity int) (int, int) {
			factor := 1 / nRoot
			numPartitions := int(math.Ceil(math.Pow(float64(capacity), factor)))
			if numPartitions > minimumPartitions {
				numPartitions = minimumPartitions
			}
			partitionLenth := int(math.Ceil(float64(capacity) / float64(numPartitions)))
			return numPartitions, partitionLenth
		}
	}
}

//endregion
