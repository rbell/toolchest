# WorkQueue
Thw workqueue package is inspired by the taskQueue at https://github.com/richardwilkes/toolbox/tree/master/taskqueue.  The differences being:
- Work queued can be prioritized
- Work's priority can be adjusted
- Errors can be monitored by observing a channel

The workqueue.queue allows the processing of work in a queued fashion.  By default, the queue's length is set to two times the number of CPUs, and the number of go routines working on the queue equal to the number of CPUs.

Work submitted to the queue takes the form of a function signature of `func() error`.  When work is performed off the queue any errors are communicated back via a channel accessible via a call to the queue's `Errors()` function.

Work submitted to the queue, by default, all have equal priority and are processed on the queue's go routines in the order they were placed.  Of course since the queue's go routines may operate in parallel, the work on the queue may not be finished in the same sequence they were put on the queue depending on go's scheduler.  

However, there is an option to set a priority when submitting work to the queue.  This will influence the order which work is processed.  Work with a lower priority number is taken first (i.e. Work X is priority 1, while work Y is priority 2, resulting in work X taking precedence over work Y).

Sometimes the priority of queued work needs adjusted.  For instance if a low priority work is queued for quite a while due to higher priority work being added to the queue.  The work's priority can be dynamically adjuested with an option when the work is submitted (see example below).

## License
This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.


## Simple Example
```go
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
```

## Configuration
Configure the queue's length:
```go
q := workqueue.NewQueue(workqueue.WithQueueLength(10))
```
Configure the number of go routines performing work on the queue:
```go
q := workqueue.NewQueue(workqueue.WithWorkers(2))
```
Configuration options can be combined:
```go
q := workqueue.NewQueue(workqueue.WithWorkers(2), workqueue.WithQueueLength(10))
```

## Prioritize Work
Work can be prioritized when submitting the work to the queue:
```go
q.QueueWork(workX, workqueue.WithPriority(1))
q.QueueWork(workY, workqueue.WithPriority(2))
```
(workX will take precidence over workY)

## Dynamically Re Prioritizing Work
Work can be re prioritized dynamically:
```go
adjustAt := time.Now().Add(time.Minute)
q.QueueWork(workX, 
	workqueue.WithPriority(100),
	workqueue.WithAdjustPriority(
		func() int {
			if time.Now() > adjustAt {
				return 1
			}
		}
    ))
```
In the above example, if workX is still queued after one minute, its priority will be set to priority one and, assuming no other work is prioritized above it, workX will be performed ahead of all other work.
## Error Monitoring
Errors returned by work placed on the queue can be monitored via a call to the queue's `Errors()` function:
```go
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
```
Each call to `Errors()` returns a unique channel allowing multiple routines to "subscribe" to errors being reported by the queue.

## Examples
Runnable examples can be found in the `toolchest/workqueue/examples` folder.