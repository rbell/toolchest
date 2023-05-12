/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
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
	// Configure worker queue to have 2 workers operating on queue of 10 work items.  Once queue has 10 items in queue, queuing will be blocked until work is completed off of queue.
	q := workqueue.NewQueue(workqueue.WithWorkers(2), workqueue.WithQueueLength(10))

	// wg to prevent app from exiting before all work is done
	wg := &sync.WaitGroup{}

	count := atomic.Int32{}
	// make 100 work functions to perform on the queue
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

	for _, w := range work {
		wg.Add(1)
		q.Enqueue(w)
	}

	wg.Wait()
}
