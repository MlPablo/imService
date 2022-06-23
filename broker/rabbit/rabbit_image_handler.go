package rabbit

import (
	"bytes"
	"image"
	"log"
	"strconv"

	"github.com/streadway/amqp"

	"imService/imageProccesor"
	"imService/storage"
)

// imageHandler consume an image from amqp.Delivery and send it to storage as image.Image
func (x RMQConsumer) imageHandler(message amqp.Delivery, db storage.Storage) error {
	image, _, err := image.Decode(bytes.NewReader(message.Body))
	if err != nil {
		return err
	}
	for quality := 25; quality <= 100; quality += 25 {
		if err := db.Add(imageProccesor.Lancoz(image, quality), message.ContentType, message.MessageId, strconv.Itoa(quality)); err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
