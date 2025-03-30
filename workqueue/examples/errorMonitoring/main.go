/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rbell/toolchest/workqueue"
)

func main() {
	// Create queue with number of workers equal to number of cpus and queue length equal to number cpus * 2
	q := workqueue.NewQueue()

	// launch go routine to monitor errors
	go func() {
		errCh := q.Errors()
		for {
			err := <-errCh
			if err != nil {
				fmt.Println(err)
			} else {
				break
			}
		}
	}()

	// wg to prevent app from exiting before all work is done
	wg := &sync.WaitGroup{}

	count := atomic.Int32{}
	// make 100 work functions to perform on the queue
	work := make([]workqueue.Work, 100)
	for i := 0; i < 100; i++ {
		work[i] = func() error {
			defer wg.Done()
			index := count.Add(1)
			time.Sleep(time.Millisecond * 100) // emulate doing some processing
			// emulate error thrown every 20th index
			if index%20 == 0 {
				return fmt.Errorf("error on index %v", index)
			} else {
				fmt.Printf("Doing some work! %v\n", index) // emulate logging
			}
			return nil
		}
	}

	for _, w := range work {
		wg.Add(1)
		q.Enqueue(w)
	}

	wg.Wait()
}
