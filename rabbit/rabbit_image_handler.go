package rabbit

import (
	"bytes"
	"image"
	"log"

	"github.com/streadway/amqp"

	"imService/storage"
)

// imageHandler consume an image from amqp.Delivery and send it to storage as image.Image
func (x RMQConsumer) imageHandler(img amqp.Delivery, db storage.Storage) error {
	im, _, err := image.Decode(bytes.NewReader(img.Body))
	if err != nil {
		return err
	}
	if err := db.Add(im, img.ContentType); err != nil {
		log.Fatal(err)
	}

	return nil
}
