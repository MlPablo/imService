package rabbit

import (
	"fmt"
	"github.com/streadway/amqp"
	"imService/storage"
	"log"
)

// RMQConsumer struct for Consumer with fields Queue name and ConnectionString
type RMQConsumer struct {
	Queue            string
	ConnectionString string
}

// NewConsumer creates new consumer
func NewConsumer(queue, ConnectionString string) RMQConsumer {
	return RMQConsumer{
		Queue:            queue,
		ConnectionString: ConnectionString,
	}
}

// Error helper func to log the errors
func (x RMQConsumer) Error(err error) {
	if err != nil {
		log.Printf("Error occurred while publishing message on '%s' queue. Error message: %s", x.Queue, err)
	}
}

// Consume starts consume messages on RMQConsumer. If queue not exists, create new
func (x RMQConsumer) Consume(db storage.Storage) {
	conn, err := amqp.Dial(x.ConnectionString)
	x.Error(err)
	defer conn.Close()

	ch, err := conn.Channel()
	x.Error(err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		x.Queue,
		false,
		false,
		false,
		false,
		nil,
	)
	x.Error(err)

	imgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	x.Error(err)

	forever := make(chan bool)

	go func() {
		for d := range imgs {
			if err := x.imageHandler(d, db); err != nil {
				log.Println(err)
			}
		}
	}()

	fmt.Printf("Started listening for messages on '%s' queue", x.Queue)
	<-forever
}
