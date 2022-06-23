package rabbit

import (
	"log"

	"github.com/streadway/amqp"
)

// RMQProducer struct for producer with Queue name and ConnectionString
type RMQProducer struct {
	Queue            string
	ConnectionString string
}

// NewProducer creates new producer
func NewProducer(queue, ConnectionString string) RMQProducer {
	return RMQProducer{Queue: queue, ConnectionString: ConnectionString}
}

// Error helper func to log the errors
func (x RMQProducer) Error(err error) {
	if err != nil {
		log.Printf("Error occurred while publishing message on '%s' queue. Error message: %s", x.Queue, err)
	}
}

// PublishMessage send a message to queue, if queue not exists, creates it
func (x RMQProducer) PublishMessage(contentType string, body []byte) {
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

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        body,
		})
	x.Error(err)
}
