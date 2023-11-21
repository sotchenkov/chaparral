package broker

import (
	"fmt"
	"testing"
)

func TestMessages(t *testing.T) {

	testTable := []struct {
		messages string
		expected string
	}{
		{
			messages: "hello, broker",
			expected: "hello, broker",
		},
	}

	// Act
	broker := NewBroker()

	sub := broker.Subscribe("test_topik")

	for _, testCase := range testTable {
		broker.Publish("test_topik", testCase.messages)
		testCase := testCase

		go func() {
			for {
				select {
				case msg, ok := <-sub.Channel:
					if !ok {
						fmt.Println("Subscriber channel closed.")
						return
					}
					// Assert
					if msg != testCase.expected {
						t.Errorf("incorrect result. Expected %s, got %s", testCase.expected, msg)
					}
				case <-sub.Unsubscribe:
					fmt.Println("Unsubscribed.")
					return
				}
			}
		}()
	}

}
