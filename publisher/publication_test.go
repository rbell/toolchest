/*
 * Copyright (c) 2025 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package publisher

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPublication(t *testing.T) {
	p := NewPublication[int]()
	assert.NotNil(t, p)
	assert.NotNil(t, p.subscribers)
	assert.Equal(t, uint64(0), p.subscriberCount.Load())
}

func TestPublication_Subscribe(t *testing.T) {
	p := NewPublication[int]()
	sub := p.Subscribe(10)
	assert.NotNil(t, sub)
	assert.Equal(t, uint64(1), sub.subscriberID)
	assert.Nil(t, sub.filter)
	assert.NotNil(t, sub.receiveCh)
	assert.Equal(t, p, sub.publisher)

	sub2 := p.Subscribe(5, WithFilter(func(i int) bool { return i > 5 }))
	assert.NotNil(t, sub2)
	assert.Equal(t, uint64(2), sub2.subscriberID)
	assert.NotNil(t, sub2.filter)
	assert.NotNil(t, sub2.receiveCh)
	assert.Equal(t, p, sub2.publisher)
}

func TestPublication_Publish(t *testing.T) {
	p := NewPublication[int]()
	sub := p.Subscribe(10)
	sub2 := p.Subscribe(5, WithFilter(func(i int) bool { return i > 5 }))

	p.Publish(10)

	assert.Eventually(t, func() bool {
		select {
		case msg := <-sub.receiveCh:
			assert.Equal(t, 10, msg)
			return true
		default:
			return false
		}
	}, time.Second, 10*time.Millisecond)

	assert.Eventually(t, func() bool {
		select {
		case msg := <-sub2.receiveCh:
			assert.Equal(t, 10, msg)
			return true
		default:
			return false
		}
	}, time.Second, 10*time.Millisecond)

}

func TestPublication_Close(t *testing.T) {
	p := NewPublication[int]()
	sub := p.Subscribe(10)
	sub2 := p.Subscribe(5, WithFilter(func(i int) bool { return i > 5 }))

	p.Close()

	_, ok := <-sub.receiveCh
	assert.False(t, ok)
	_, ok = <-sub2.receiveCh
	assert.False(t, ok)

	// verify subscribers are removed
	var count int
	p.subscribers.Range(func(key uint64, value *Subscriber[int]) bool {
		count++
		return true
	})
	assert.Equal(t, 0, count)
}

func TestPublication_unsubscribe(t *testing.T) {
	p := NewPublication[int]()
	sub := p.Subscribe(10)
	sub2 := p.Subscribe(5, WithFilter(func(i int) bool { return i > 5 }))

	p.unsubscribe(sub.subscriberID)

	// verify subscriber is removed
	var count int
	p.subscribers.Range(func(key uint64, value *Subscriber[int]) bool {
		count++
		return true
	})
	assert.Equal(t, 1, count)

	// verify channel is closed
	_, ok := <-sub.receiveCh
	assert.False(t, ok)

	// verify other channel is still open
	p.Publish(10)
	assert.Eventually(t, func() bool {
		select {
		case msg := <-sub2.receiveCh:
			assert.Equal(t, 10, msg)
			return true
		default:
			return false
		}
	}, time.Second, 10*time.Millisecond)
}

func TestSubscriber_Close(t *testing.T) {
	p := NewPublication[int]()
	sub := p.Subscribe(1)
	sub2 := p.Subscribe(5, WithFilter(func(i int) bool { return i > 5 }))

	sub.Close()

	// verify subscriber is removed
	var count int
	p.subscribers.Range(func(key uint64, value *Subscriber[int]) bool {
		count++
		return true
	})
	assert.Equal(t, 1, count)

	// verify channel is closed
	_, ok := <-sub.receiveCh
	assert.False(t, ok)

	// verify other channel is still open
	p.Publish(10)
	assert.Eventually(t, func() bool {
		select {
		case msg := <-sub2.receiveCh:
			assert.Equal(t, 10, msg)
			return true
		default:
			return false
		}
	}, time.Second, 10*time.Millisecond)
}

func TestPublication_Publish_Concurrent(t *testing.T) {
	p := NewPublication[int]()
	numSubscribers := 100
	numMessages := 1000
	var wg sync.WaitGroup
	wg.Add(numSubscribers)

	for i := 0; i < numSubscribers; i++ {
		sub := p.Subscribe(10)
		go func() {
			defer wg.Done()
			for j := 0; j < numMessages; j++ {
				select {
				case <-sub.Receive():
					// Received message
				case <-time.After(time.Second):
					assert.Fail(t, "timeout waiting for message")
					return
				}
			}
		}()
	}

	for i := 0; i < numMessages; i++ {
		p.Publish(i)
	}

	wg.Wait()
	p.Close()
}

func TestPublication_Subscribe_Concurrent(t *testing.T) {
	p := NewPublication[int]()
	numSubscribers := 100
	var wg sync.WaitGroup
	wg.Add(numSubscribers)

	for i := 0; i < numSubscribers; i++ {
		go func() {
			defer wg.Done()
			p.Subscribe(10)
		}()
	}

	wg.Wait()
	assert.Equal(t, uint64(numSubscribers), p.subscriberCount.Load())
	p.Close()
}

func TestPublication_Publish_WithFilter(t *testing.T) {
	p := NewPublication[int]()
	sub1 := p.Subscribe(10, WithFilter(func(i int) bool { return i%2 == 0 }))
	sub2 := p.Subscribe(10, WithFilter(func(i int) bool { return i%2 != 0 }))

	p.Publish(2)
	p.Publish(3)

	assert.Eventually(t, func() bool {
		select {
		case msg := <-sub1.Receive():
			assert.Equal(t, 2, msg)
			return true
		default:
			return false
		}
	}, time.Second, 10*time.Millisecond)

	select {
	case <-sub1.Receive():
		assert.Fail(t, "should not have received message")
	default:
		// continue
	}

	assert.Eventually(t, func() bool {
		select {
		case msg := <-sub2.Receive():
			assert.Equal(t, 3, msg)
			return true
		default:
			return false
		}
	}, time.Second, 10*time.Millisecond)

	select {
	case <-sub2.Receive():
		assert.Fail(t, "should not have received message")
	default:
		// continue
	}
	p.Close()
}

func TestPublication_Publish_WithFilter_Concurrent(t *testing.T) {
	p := NewPublication[int]()
	numSubscribers := 100
	numMessages := 1000
	var wg sync.WaitGroup
	wg.Add(numSubscribers)

	for i := 0; i < numSubscribers; i++ {
		sub := p.Subscribe(10, WithFilter(func(i int) bool { return i%2 == 0 }))
		go func() {
			defer wg.Done()
			for j := 0; j < numMessages; j++ {
				select {
				case msg := <-sub.Receive():
					assert.Equal(t, 0, msg%2)
				case <-time.After(time.Second):
					assert.Fail(t, "timeout waiting for message")
					return
				}
			}
		}()
	}

	for i := 0; i < numMessages; i++ {
		p.Publish(i * 2)
	}

	wg.Wait()
	p.Close()
}

func TestPublication_Publish_WithFilter_NoMatch(t *testing.T) {
	p := NewPublication[int]()
	sub := p.Subscribe(10, WithFilter(func(i int) bool { return i > 100 }))

	p.Publish(1)

	select {
	case <-sub.Receive():
		assert.Fail(t, "should not have received message")
	default:
		// continue
	}
	p.Close()
}

func TestPublication_Publish_WithFilter_NoMatch_Concurrent(t *testing.T) {
	p := NewPublication[int]()
	numSubscribers := 100
	numMessages := 1000
	var wg sync.WaitGroup
	wg.Add(numSubscribers)

	for i := 0; i < numSubscribers; i++ {
		sub := p.Subscribe(10, WithFilter(func(i int) bool { return i < 0 }))
		go func() {
			defer wg.Done()
			for j := 0; j < numMessages; j++ {
				select {
				case <-sub.Receive():
					assert.Fail(t, "should not have received message")
				case <-time.After(time.Millisecond * 10):
					// continue
				}
			}
		}()
	}

	for i := 0; i < numMessages; i++ {
		p.Publish(i)
	}

	wg.Wait()
	p.Close()
}
