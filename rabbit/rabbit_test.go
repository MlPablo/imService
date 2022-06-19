package rabbit_test

import (
	"github.com/go-playground/assert/v2"
	"imService/rabbit"
	"testing"
	"time"
)

// TestQueue checks if Producer and Consumer works properly
func TestQueue(t *testing.T) {
	queue := rabbit.NewProducer("TestImageQue", "amqp://guest:guest@localhost:5672")
	consumer := rabbit.NewConsumer("TestImageQue", "amqp://guest:guest@localhost:5672")
	var counter int
	go consumer.ConsumeTest(&counter)
	for i := 0; i < 100; i++ {
		go queue.PublishMessage("1", []byte{0})
	}
	time.Sleep(time.Second)
	assert.Equal(t, counter, 100)
}
