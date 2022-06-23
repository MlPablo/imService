package rabbit_test

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"

	"imService/broker/rabbit"
)

// TestQueue checks if Producer and Consumer works properly
func TestQueue(t *testing.T) {
	broker := rabbit.NewRabbit("TestQueue", "amqp://guest:guest@localhost:5672/")
	var counter int
	go broker.Consumer.ConsumeTest(&counter)
	for i := 0; i < 100; i++ {
		go broker.PublishMessage("1", "abcd", []byte{0})
	}
	time.Sleep(time.Second)
	assert.Equal(t, counter, 100)
}
