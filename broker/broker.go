package broker

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

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "quic"
)

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

// func main() {
// 	broker := NewBroker()

// 	subscriber1 := broker.Subscribe("topik_num_1")
// 	go func() {
// 		for {
// 			select {
// 			case msg, ok := <-subscriber1.Channel:
// 				if !ok {
// 					fmt.Println("Subscriber channel closed.")
// 					return
// 				}
// 				fmt.Printf("Received: %v\n", msg)
// 			case <-subscriber1.Unsubscribe:
// 				fmt.Println("Unsubscribed.")
// 				return
// 			}
// 		}
// 	}()

// 	subscriber2 := broker.Subscribe("topik_num_2")
// 	go func() {
// 		for {
// 			select {
// 			case msg, ok := <-subscriber2.Channel:
// 				if !ok {
// 					fmt.Println("Subscriber channel closed.")
// 					return
// 				}
// 				fmt.Printf("Received: %v\n", msg)
// 			case <-subscriber2.Unsubscribe:
// 				fmt.Println("Unsubscribed.")
// 				return
// 			}
// 		}
// 	}()

// 	broker.Publish("topik_num_1", "Hello, World!")
// 	broker.Publish("topik_num_2", "This is a test message.")

// 	time.Sleep(2 * time.Second)
// 	broker.Unsubscribe("topik_num_1", subscriber1)

// 	broker.Publish("topik_num_1", "This message won't be received.")

// 	time.Sleep(time.Second)
// }

func checkErr(err error, args ...interface{}) {
	if err != nil {
		panic(fmt.Sprint(append([]interface{}{err}, args...)...))
	}
}
