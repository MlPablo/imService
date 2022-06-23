package storage

import (
	"image/png"
	"log"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"

	"imService/test_image"
)

func CreateStoreWithOneImage() *localdb {
	store := &localdb{Storage: make(map[fileName]file)}
	fil, _ := os.Open(test_image.PathToTestImage)
	image, _ := png.Decode(fil)
	if err := store.Add(image, "image/png", "aabcue", "25"); err != nil {
		log.Fatal(err)
	}
	return store
}

func (store *localdb) AddOneTestImage() {
	file, _ := os.Open(test_image.PathToTestImage)
	image, _ := png.Decode(file)
	if err := store.Add(image, "image/png", "aabcue", "25"); err != nil {
		log.Fatal(err)
	}
}

func TestLocaldb_Add(t *testing.T) {
	store := CreateStoreWithOneImage()
	_, ok := store.Storage[fileName{Name: "aabcue", Quality: "25"}]
	assert.Equal(t, ok, true)

	_, ok = store.Storage[fileName{Name: "2", Quality: "25"}]
	log.Println(ok)
	assert.Equal(t, ok, false)

}

func TestLocaldb_Get(t *testing.T) {
	store := CreateStoreWithOneImage()

	img, err := store.Get("aabcue", "25")
	assert.Equal(t, err, nil)
	assert.NotEqual(t, img.Image, nil)

	img, err = store.Get("10", "50")
	assert.NotEqual(t, err, nil)
	assert.Equal(t, img.Image, nil)

	img, err = store.Get("1", "17")
	assert.NotEqual(t, err, nil)
	assert.Equal(t, img.Image, nil)

	img, err = store.Get("abf", "sde")
	assert.NotEqual(t, err, nil)
	assert.Equal(t, img.Image, nil)

}
