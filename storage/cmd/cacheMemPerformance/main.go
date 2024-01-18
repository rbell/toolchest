/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package main

import (
	"context"
	"fmt"
	"github.com/rbell/toolchest/storage"
	"runtime"
	"time"
)

const (
	capacity = 1000000
	size     = 128
)

func main() {
	cache := storage.NewFifoMapCache[int, *[size]byte](context.Background(), capacity/2, storage.WithSweepFrequency(time.Minute*10))

	actualCapacity := cache.Capacity()
	fmt.Printf("Actual Capacity %v\n", actualCapacity)
	partialStart := int(.75 * float64(capacity))
	partialEnd := capacity
	fmt.Printf("Partial Start %v, Partial End %v\n", partialStart, partialEnd)

	fmt.Println("Start measuring ...")

	for i := 0; i < capacity; i++ {
		cache.Set(i, &[size]byte{})
	}
	fmt.Printf("Cache Length After Adding Capacity %v\n", cache.Len())
	printAlloc("After Adding Capacity")
	cache.Sweep()
	fmt.Printf("Cache Length After Forcing Sweep %v\n", cache.Len())
	printAlloc("After Forcing Sweep")

	for i := 0; i < capacity; i++ {
		cache.Set(i, &[size]byte{})
	}
	cache.Sweep()
	fmt.Printf("Cache Length Adding Everything Again %v\n", cache.Len())
	printAlloc("After Forcing Sweep")

	fmt.Println("Deleting part of the cache ...")
	for i := partialStart; i < partialEnd; i++ {
		cache.Delete(i)
	}
	cache.Sweep()
	fmt.Printf("Cache Length After Deleting Partial %v\n", cache.Len())
	printAlloc("After Deleting Partial")

	fmt.Println("Adding part of the cache back ...")
	for i := partialStart; i < partialEnd; i++ {
		cache.Set(i, &[size]byte{})
	}
	cache.Sweep()
	fmt.Printf("Cache Length After Adding Parital Back %v\n", cache.Len())
	printAlloc("After Adding Partial Back")

	fmt.Println("Deleting part of the cache again ...")
	for i := partialStart; i < partialEnd; i++ {
		cache.Delete(i)
	}
	cache.Sweep()
	fmt.Printf("Cache Length After Deleting Again %v\n", cache.Len())
	printAlloc("After Deleting Again")

	for i := 0; i < actualCapacity; i++ {
		cache.Set(i, &[size]byte{})
	}
	fmt.Printf("Cache Length After Adding All Back In %v\n", cache.Len())
	printAlloc("After Adding Capacity")
	cache.Sweep()
	fmt.Printf("Cache Length After Adding All Back In %v\n", cache.Len())
	printAlloc("After Forcing Sweep")

	runtime.KeepAlive(cache)
}

func printAlloc(msg string) {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc (%v) = %v MB\n", msg, m.Alloc/(1024*1024))
}
