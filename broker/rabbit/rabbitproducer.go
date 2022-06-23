package rabbit

import (
	"log"

	"github.com/streadway/amqp"
)

// Error helper func to log the errors
func (x RMQProducer) Error(err error) {
	if err != nil {
		log.Printf("Error occurred while publishing message on '%s' queue. Error message: %s", x.Queue, err)
	}
}

// PublishMessage send a message to queue, if queue not exists, creates it
func (x BrokerRabbit) PublishMessage(contentType string, body []byte) {
	conn, err := amqp.Dial(x.Producer.ConnectionString)
	x.Producer.Error(err)
	defer conn.Close()

	ch, err := conn.Channel()
	x.Producer.Error(err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		x.Producer.Queue,
		false,
		false,
		false,
		false,
		nil,
	)
	x.Producer.Error(err)

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        body,
		})
	x.Producer.Error(err)
}
