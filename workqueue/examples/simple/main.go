/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package main

import (
	"fmt"
	"github.com/rbell/toolchest/workqueue"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// Create queue with number of workers equal to number of cpus and queue length equal to number cpus * 2
	q := workqueue.NewQueue()

	// wg to prevent app from exiting before all work is done
	wg := &sync.WaitGroup{}

	count := atomic.Int32{}

	// make 100 work functions to perform on the queue - each work function will increment count and print the resulting value, then emulate doing some work.
	work := make([]workqueue.Work, 100)
	for i := 0; i < 100; i++ {
		work[i] = func() error {
			defer wg.Done()
			index := count.Add(1)
			fmt.Printf("Doing some work! %v\n", index) // emulate logging
			time.Sleep(time.Millisecond * 100)         // emulate doing some processing
			return nil
		}
	}

	// Queue the work
	for _, w := range work {
		wg.Add(1)
		q.QueueWork(w)
	}

	wg.Wait()
}
