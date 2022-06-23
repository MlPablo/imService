package broker

import (
	"os"

	"imService/broker/rabbit"
	"imService/storage"
)

type Broker interface {
	Consume(db storage.Storage)
	PublishMessage(string, string, []byte)
}

func NewQue() Broker {
	return rabbit.NewRabbit(os.Getenv("RABBIT_QUE"), os.Getenv("RABBIT_PATH"))
}
