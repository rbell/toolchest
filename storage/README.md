# Storage Package

The `storage` package provides concurrency-safe map and cache implementations in Go.

## SafeMap

`SafeMap` is a struct that wraps a map with a simple `RWMutex` to facilitate concurrency. It supports generic types for keys and values.

## FifoMapCache

`FifoMapCache` is a struct that implements a First-In-First-Out (FIFO) cache with a maximum size. When the cache is full, the oldest entries are evicted. It supports generic types for keys and values.

## GenericStack

`GenericStack` is a struct that implements a generic stack data structure. It supports any type of values.

## Testing

Unit tests for the `SafeMap` and `FifoMapCache` are located in the `storage/safeMap_test.go` and `storage/fifoMapCache_test.go` files respectively. They cover all methods and some edge cases, including concurrent operations.