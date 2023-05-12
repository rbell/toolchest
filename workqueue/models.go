/*
 * Copyright (c) 2023  by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package workqueue

import (
	"sync/atomic"

	"github.com/google/uuid"
)

type workState int32

// work states
const (
	IN_QUEUE workState = iota
	IN_PROGRESS
)

func (ws workState) String() string {
	switch ws {
	case IN_PROGRESS:
		return "In Progress"
	case IN_QUEUE:
		return "Queued"
	}
	return "unknown"
}

type Work func() error

type WorkQueueOption func(*Queue)

type QueuedWork struct {
	id       uuid.UUID
	name     string
	priority int
	position int
	state    *atomic.Int32
}

func (w *QueuedWork) Id() string {
	return w.id.String()
}

func (w *QueuedWork) Name() string {
	return w.name
}

func (w *QueuedWork) Priority() int {
	return w.priority
}

func (w *QueuedWork) State() string {
	st := workState(w.state.Load())
	return st.String()
}

type workItem struct {
	*QueuedWork
	workToDo       Work
	adjustPriority func() int
}
