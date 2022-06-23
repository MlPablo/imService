package rabbit

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"

	"imService/storage"
)

// Error helper func to log the errors
func (x RMQConsumer) Error(err error) {
	if err != nil {
		log.Printf("Error occurred while publishing message on '%s' queue. Error message: %s", x.Queue, err)
	}
}

// Consume starts consume messages on RMQConsumer. If queue not exists, create new
func (x BrokerRabbit) Consume(db storage.Storage) {
	conn, err := amqp.Dial(x.Consumer.ConnectionString)
	x.Consumer.Error(err)
	defer conn.Close()

	ch, err := conn.Channel()
	x.Consumer.Error(err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		x.Consumer.Queue,
		false,
		false,
		false,
		false,
		nil,
	)
	x.Consumer.Error(err)

	imgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	x.Consumer.Error(err)

	forever := make(chan bool)

	go func() {
		for d := range imgs {
			if err := x.Consumer.imageHandler(d, db); err != nil {
				log.Println(err)
			}
		}
	}()

	fmt.Printf("Started listening for messages on '%s' queue", x.Consumer.Queue)
	<-forever
}
