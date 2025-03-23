/*
 * Copyright (c) 2025 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package publisher

import (
	"sync/atomic"
	"time"

	"github.com/rbell/toolchest/generic"
)

// defaultTimeout is the default timeout used when publishing a message.
var defaultTimeout = 10 * time.Second

// Publication is a generic struct that manages the publication of messages to subscribers.
type Publication[T any] struct {
	subscribers     *generic.SyncMap[uint64, *Subscriber[T]]
	subscriberCount atomic.Uint64
}

// Subscriber is a struct that represents a subscriber to a publication.
type Subscriber[T any] struct {
	subscriberID uint64
	filter       func(T) bool
	onFiltered   func(T)
	receiveCh    chan T
	publisher    *Publication[T]
	timeout      time.Duration
	onTimeout    func(T)
}

type SubscriberOption[T any] func(sub *Subscriber[T])

// NewPublication creates a new Publication.
func NewPublication[T any]() *Publication[T] {
	return &Publication[T]{
		subscribers:     generic.NewSyncMap[uint64, *Subscriber[T]](),
		subscriberCount: atomic.Uint64{},
	}
}

// Subscribe creates a new subscriber to the publication.
//
//	buffer: the size of the buffer for the subscriber's receive channel.
//	options: optional options for the subscriber.
//
// Returns a reference to the new subscriber.
func (p *Publication[T]) Subscribe(buffer int, opts ...SubscriberOption[T]) *Subscriber[T] {
	sub := &Subscriber[T]{
		subscriberID: p.subscriberCount.Add(1),
		receiveCh:    make(chan T, buffer),
		publisher:    p,
		timeout:      defaultTimeout,
	}

	for _, opt := range opts {
		opt(sub)
	}

	p.subscribers.Store(sub.subscriberID, sub)
	return sub
}

// Publish publishes a message to all subscribers.
//
//	message: the message to publish.
//	timeout: the timeout for sending the message to each subscriber. If nil, the default timeout is used.
func (p *Publication[T]) Publish(message T) {
	for _, sub := range p.subscribers.Iterate() {
		if sub.filter == nil || sub.filter(message) {
			go func() {
				select {
				case sub.receiveCh <- message:
					// continue
				case <-time.After(sub.timeout):
					// continue
				}
			}()
		}
	}
}

// Close closes the publication and all subscriber channels.
func (p *Publication[T]) Close() {
	for _, listener := range p.subscribers.Iterate() {
		close(listener.receiveCh)
	}
	p.subscribers.Clear()
}

// unsubscribe removes a subscriber from the publication.
func (p *Publication[T]) unsubscribe(subscriberID uint64) {
	if s, ok := p.subscribers.Load(subscriberID); ok {
		close(s.receiveCh)
		p.subscribers.Delete(subscriberID)
	}
}

// Close closes the subscriber's receive channel and unsubscribes them from the publication.
func (s *Subscriber[T]) Close() {
	s.publisher.unsubscribe(s.subscriberID)
}

// Receive returns the subscriber's receive channel.
func (s *Subscriber[T]) Receive() <-chan T {
	return s.receiveCh
}

//region Subscriber Options

func WithFilter[T any](filter func(T) bool) SubscriberOption[T] {
	return func(sub *Subscriber[T]) {
		sub.filter = filter
	}
}

func WithTimeout[T any](timeout time.Duration) SubscriberOption[T] {
	return func(sub *Subscriber[T]) {
		sub.timeout = timeout
	}
}

func OnTimeout[T any](onTimeout func(T)) SubscriberOption[T] {
	return func(sub *Subscriber[T]) {
		sub.onTimeout = onTimeout
	}
}

func OnFiltered[T any](onFiltered func(T)) SubscriberOption[T] {
	return func(sub *Subscriber[T]) {
		sub.onFiltered = onFiltered
	}
}

//endregion
