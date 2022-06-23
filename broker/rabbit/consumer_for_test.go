package rabbit

import (
	"fmt"

	"github.com/streadway/amqp"
)

//ConsumeTest just for testing. Counts amount of consumed massages
func (x RMQConsumer) ConsumeTest(counter *int) {
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
		for range imgs {
			*counter++
		}
	}()

	fmt.Printf("Started listening for messages on '%s' queue", x.Queue)
	<-forever
}
