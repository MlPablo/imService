package rabbit

import (
	"bytes"
	"github.com/streadway/amqp"
	"imService/storage"
	"image"
)

// imageHandler consume an image from amqp.Delivery and send it to storage as image.Image
func (x RMQConsumer) imageHandler(img amqp.Delivery, db storage.Storage) error {
	im, _, err := image.Decode(bytes.NewReader(img.Body))
	if err != nil {
		return err
	}
	db.Add(im, img.ContentType)
	return nil
}
