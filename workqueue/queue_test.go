/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package workqueue

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestQueue_PerformsWork(t *testing.T) {
	// setup
	q := NewQueue()
	wg := &sync.WaitGroup{}

	work := make([]Work, 100)
	for i := 0; i < 100; i++ {
		count := i
		work[i] = func() error {
			defer wg.Done()
			fmt.Println(count)
			time.Sleep(time.Millisecond * 100)
			return nil
		}
	}

	// test
	for i := 0; i < 100; i++ {
		wg.Add(1)
		priority := i % 2
		q.QueueWork(work[i], WithPriority(priority))
	}

	// assert
	wg.Wait()
}
