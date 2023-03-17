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
	"time"
)

func main() {
	// Create queue with number of workers equal to number of cpus and queue length equal to number cpus * 2
	q := workqueue.NewQueue()

	// wg to prevent app from exiting before all work is done
	wg := &sync.WaitGroup{}

	type workTask struct {
		priority int
		work     workqueue.Work
	}

	// make 100 work functions to perform on the queue
	work := make([]*workTask, 100)
	for i := 0; i < 100; i++ {
		index := i + 1
		priority := index % 3
		work[i] = &workTask{
			priority: priority,
			work: func() error {
				defer wg.Done()
				fmt.Printf("Doing some work with priority %v! %v\n", priority, index) // emulate logging
				time.Sleep(time.Millisecond * 100)                                    // emulate doing some processing
				return nil
			},
		}
	}

	for _, w := range work {
		wg.Add(1)
		q.QueueWork(w.work, workqueue.WithPriority(w.priority))
	}

	wg.Wait()
}
