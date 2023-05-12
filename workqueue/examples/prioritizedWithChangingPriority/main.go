/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package main

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rbell/toolchest/workqueue"
)

func main() {
	// Create queue
	q := workqueue.NewQueue(workqueue.WithWorkers(2), workqueue.WithQueueLength(200))

	// wg to prevent app from exiting before all work is done
	wg := &sync.WaitGroup{}

	type workTask struct {
		priority int
		work     workqueue.Work
	}

	// make 100 work functions to perform on the queue
	work := make([]*workTask, 100)
	for i := 0; i < 100; i++ {
		work[i] = &workTask{
			priority: 10,
			work: func() error {
				defer wg.Done()
				time.Sleep(time.Second) // emulate doing some processing
				return nil
			},
		}
	}

	for i, w := range work {
		wg.Add(1)
		q.Enqueue(w.work, workqueue.WithPriority(w.priority), workqueue.WithName(fmt.Sprintf("Work Task %v", i)))
	}

	fmt.Println("Before Adding Prioritized Work")
	printItems(q)

	// Emulate situation where priority changes dynamically
	priority := &atomic.Int32{}
	priority.Store(99)
	timer := time.NewTimer(time.Second * 3)
	go func() {
		<-timer.C
		priority.Store(1)
		time.Sleep(time.Second)
		fmt.Println("After Priority Change")
		printItems(q)
	}()

	// Queue prioritized task
	wg.Add(1)
	q.Enqueue(func() error {
		defer wg.Done()
		time.Sleep(time.Second)
		fmt.Println("Done with prioritized work!")
		return nil
	}, workqueue.WithPriority(99), workqueue.WithName("Prioritized"), workqueue.WithAdjustPriority(func() int {
		return int(priority.Load())
	}))

	// q.QueueWork(func() error {
	// 	defer wg.Done()
	// 	fmt.Println("Done with prioritized work!")
	// 	return nil
	// }, workqueue.WithPriority(1), workqueue.WithName("Prioritized"))

	fmt.Println("After Queing Prioritized Work")
	printItems(q)

	wg.Wait()
}

func printItems(q *workqueue.Queue) {
	fmt.Printf("ID\t\t\t\t\tName\t\t\tPriority\tState\n")

	items := q.WorkItems()
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Priority() < items[j].Priority()
	})
	for _, wi := range items {
		fmt.Printf("%v\t%v\t\t%v\t\t%v\n", wi.Id(), wi.Name(), wi.Priority(), wi.State())
	}
}
