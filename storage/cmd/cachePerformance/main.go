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
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	iterations = 1000000
	threads    = 10
)

type (
	data struct {
		key int
		val int
	}
)

var (
	cache     = make(map[int]int, 100)
	syncCache sync.Map
	mutex     sync.Mutex
	mc        = storage.NewFifoMapCache[string, int](context.Background(), iterations)
	safeMap   = storage.NewSafeMap[string, int](iterations)
	ch        = make(chan data)
)

func measure(name string, f func()) {
	fmt.Println("Start measuring", name, "...")
	start := time.Now()
	f()
	taken := time.Since(start)
	fmt.Printf("Finished measuring %s, time taken: %v\n\n", name, taken)
}

func exec(meth0d func(i int)) {
	wg := new(sync.WaitGroup)
	wg.Add(threads)
	for t := 0; t < threads; t++ {
		go func() {
			for i := 0; i < iterations; i++ {
				meth0d(i)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
func main() {
	measure("Mutex", func() {
		exec(func(i int) {
			mutex.Lock()
			cache[i%100000] += 1
			mutex.Unlock()
		})
	})
	measure("FifoCache", func() {
		exec(func(i int) {
			mc.Set(strconv.Itoa(i), 1)
		})
	})
	measure("SafeMap", func() {
		exec(func(i int) {
			safeMap.Set(strconv.Itoa(i), 1)
		})
	})
	measure("sync.Map", func() {
		exec(func(i int) {
			elem, _ := syncCache.LoadOrStore(i, new(int32))
			atomic.AddInt32(elem.(*int32), 1)
		})
	})
	measure("Channels", func() {
		go func() {
			for x := range ch {
				cache[x.key] += x.val
			}
		}()
		exec(func(i int) {
			ch <- data{i, 1}
		})
	})
}
