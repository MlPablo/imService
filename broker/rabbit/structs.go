package rabbit

type BrokerRabbit struct {
	Consumer *RMQConsumer
	Producer *RMQProducer
}

// RMQProducer struct for producer with Queue name and ConnectionString
type RMQProducer struct {
	Queue            string
	ConnectionString string
}

// RMQConsumer struct for Consumer with fields Queue name and ConnectionString
type RMQConsumer struct {
	Queue            string
	ConnectionString string
}

// NewProducer creates new producer
func NewProducer(queue, ConnectionString string) *RMQProducer {
	return &RMQProducer{Queue: queue, ConnectionString: ConnectionString}
}

// NewConsumer creates new consumer
func NewConsumer(queue, ConnectionString string) *RMQConsumer {
	return &RMQConsumer{
		Queue:            queue,
		ConnectionString: ConnectionString,
	}
}

func NewRabbit(queue, path string) *BrokerRabbit {
	return &BrokerRabbit{
		NewConsumer(queue, path),
		NewProducer(queue, path),
	}
}
