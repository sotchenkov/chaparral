package main

import (
	"fmt"
	"sync"
	"time"
)

type Message struct {
	Topic   string
	Payload interface{}
}

type Subscriber struct {
	Channel     chan interface{}
	Unsubscribe chan bool
}

type Broker struct {
	subscribers map[string][]*Subscriber
	mutex       sync.Mutex
}

func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[string][]*Subscriber),
	}
}

/*
	The Subscribe method allows clients to subscribe to
	specific topics and receive messages related to those topics
*/
func (b *Broker) Subscribe(topic string) *Subscriber {
	b.mutex.Lock() // to ensure exclusive access to the subscribers map
	defer b.mutex.Unlock()

	subscriber := &Subscriber{
		Channel:     make(chan interface{}, 1),
		Unsubscribe: make(chan bool),
	}

	b.subscribers[topic] = append(b.subscribers[topic], subscriber)

	return subscriber
}

/*
The Unsubscribe method lets subscribers unsubscribe
from topics they're no longer interested in
*/
func (b *Broker) Unsubscribe(topic string, subscriber *Subscriber) {
	b.mutex.Lock() // acquire the mutex lock to safely access the subscribers map
	defer b.mutex.Unlock()

	if subscribers, found := b.subscribers[topic]; found {
		for i, sub := range subscribers {
			if sub == subscriber {
				close(sub.Channel)
				b.subscribers[topic] = append(subscribers[:i], subscribers[i+1:]...)
				return
			}
		}
	}
}

/*
The Publish method sends messages to subscribers of a specific topic
*/
func (b *Broker) Publish(topic string, payload interface{}) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if subscribers, found := b.subscribers[topic]; found {
		for _, sub := range subscribers {
			select {
			case sub.Channel <- payload:
			case <-time.After(time.Second):
				fmt.Printf("Subscriber slow. Unsubscribing from topic: %s\n", topic)
				b.Unsubscribe(topic, sub)
			}
		}
	}
}

func main() {
	broker := NewBroker()

	subscriber := broker.Subscribe("example_topic")
	go func() {
		for {
			select {
			case msg, ok := <-subscriber.Channel:
				if !ok {
					fmt.Println("Subscriber channel closed.")
					return
				}
				fmt.Printf("Received: %v\n", msg)
			case <-subscriber.Unsubscribe:
				fmt.Println("Unsubscribed.")
				return
			}
		}
	}()

	broker.Publish("example_topic", "Hello, World!")
	broker.Publish("example_topic", "This is a test message.")

	time.Sleep(2 * time.Second)
	broker.Unsubscribe("example_topic", subscriber)

	broker.Publish("example_topic", "This message won't be received.")

	time.Sleep(time.Second)
}
