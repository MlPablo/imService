package storage

import (
	"image/png"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/go-playground/assert/v2"

	"imService/test_image"
)

func CreateStoreWithOneImage() *localdb {
	store := &localdb{Storage: make(map[fileName]file), CurrentId: 1}
	file, _ := os.Open(test_image.PathToTestImage)
	image, _ := png.Decode(file)
	if err := store.Add(image, "image/png"); err != nil {
		log.Fatal(err)
	}
	return store
}

func (store *localdb) AddOneTestImage() {
	file, _ := os.Open(test_image.PathToTestImage)
	image, _ := png.Decode(file)
	if err := store.Add(image, "image/png"); err != nil {
		log.Fatal(err)
	}
}

func TestLocaldb_Add(t *testing.T) {
	store := CreateStoreWithOneImage()
	for x := 25; x <= 100; x += 25 {
		_, ok := store.Storage[fileName{Name: "1", Quality: strconv.Itoa(x)}]
		assert.Equal(t, ok, true)
	}

	_, ok := store.Storage[fileName{Name: "2", Quality: "25"}]
	assert.Equal(t, ok, false)

	store.AddOneTestImage()
	_, ok = store.Storage[fileName{Name: "2", Quality: "25"}]
	assert.Equal(t, ok, true)

}

func TestLocaldb_Get(t *testing.T) {
	store := CreateStoreWithOneImage()
	store.AddOneTestImage()
	store.AddOneTestImage()

	img, err := store.Get("3", "50")
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
