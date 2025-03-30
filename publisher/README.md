# Publication Package

The `publication` package provides a generic publish-subscribe (pub/sub) mechanism for Go applications. It allows you to create publications to which multiple subscribers can listen. When a message is published, it's distributed to all relevant subscribers.

## Key Features

*   **Generic:** Supports publishing and subscribing to messages of any type.
*   **Filtering:** Subscribers can define filters to receive only messages that meet specific criteria.
*   **Buffered Channels:** Subscribers receive messages through buffered channels, allowing for asynchronous message handling.
*   **Concurrency-Safe:** Designed for concurrent use, ensuring safe message delivery in multi-threaded environments.
*   **Timeout Support:** Control the maximum time spent attempting to deliver a message to a subscriber.
*   **Clean Shutdown:** Gracefully close publications and subscriber channels.
* **Unsubscribe:** Subscribers can unsubscribe from a publication.

## Core Components

### `Publication[T]`

The `Publication` struct is the central component of the pub/sub system. It manages the list of subscribers and handles message distribution.

*   **`NewPublicationT *Publication[T]`:** Creates a new publication for messages of type `T`.
*   **`Subscribe(buffer int, filter func(T) bool) *Subscriber[T]`:** Adds a new subscriber to the publication.
    *   `buffer`: The size of the buffer for the subscriber's receive channel.
    *   `filter`: An optional function that determines whether a subscriber should receive a message. If `nil`, all messages are received.
*   **`Publish(message T, timeout *time.Duration)`:** Publishes a message to all subscribers.
    *   `message`: The message to publish.
    *   `timeout`: An optional timeout for sending the message to each subscriber. If `nil`, a default timeout of 10 seconds is used.
*   **`Close()`:** Closes the publication and all subscriber channels.

### `Subscriber[T]`

The `Subscriber` struct represents a single subscriber to a publication.

*   **`Close()`:** Closes the subscriber's receive channel and unsubscribes them from the publication.
*   **`Receive() <-chan T`:** Returns the subscriber's receive channel, from which messages can be read.

## Usage Examples

### Basic Publication and Subscription
```go
package main

import (
	"fmt"
	"time"
    "github.com/rbell/toolchest/publisher" // Replace with your actual import path
)

func main() {
	// Create a new publication for integers. 
	pub := publisher.NewPublication[int]()
	
    // Subscribe to the publication with a buffer of 10.
    sub := pub.Subscribe(10)
	
    // Publish some messages.
    pub.Publish(1)
    pub.Publish(2)
    pub.Publish(3)
	
    // Receive messages from the subscriber.
    for i := 0; i < 3; i++ {
        select {
        case msg := <-sub.Receive():
            fmt.Println("Received:", msg)
        case <-time.After(time.Second):
            fmt.Println("Timeout waiting for message")
            return
        }
    }
    // Close the publication.
    pub.Close()
}
```

### Filtering Messages
```go
package main
import (
	"fmt"
	"time"
    "github.com/rbell/toolchest/publisher" // Replace with your actual import path
)

func main() {
	// Create a new publication for integers. 
	pub := publisher.NewPublication[int]()
	
    // Subscribe to even numbers only.
    evenSub := pub.Subscribe(10, publisher.WithFilter(func(i int) bool { return i%2 == 0 }))

	// Subscribe to odd numbers only.
    oddSub := pub.Subscribe(10, publisher.WithFilter(func(i int) bool { return i%2 != 0 }))

	pub.Publish(1)
    pub.Publish(2)
    pub.Publish(3)
    pub.Publish(4)
	
    // Receive even messages.
    for i := 0; i < 2; i++ {
        select {
        case msg := <-evenSub.Receive():
            fmt.Println("Even Received:", msg)
        case <-time.After(time.Second):
            fmt.Println("Timeout waiting for even message")
            return
        }
    }
    // Receive odd messages.
    for i := 0; i < 2; i++ {
        select {
        case msg := <-oddSub.Receive():
            fmt.Println("Odd Received:", msg)
        case <-time.After(time.Second):
            fmt.Println("Timeout waiting for odd message")
            return
        }
    }
    pub.Close()
}
```