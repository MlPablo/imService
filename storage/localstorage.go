package storage

import (
	"errors"
	"image"
	"log"
)

// localdb struct of local DataBase
type localdb struct {
	Storage map[fileName]file
}

// Add is method for localdb. It takes an image and creates 4 of them with different resize (25, 50, 75, 100) and adds it to localdb
// Stores images in format "id_quality.format"
func (store *localdb) Add(image image.Image, ext, id, quality string) error {
	store.Storage[fileName{Name: id, Quality: quality}] = file{image, ext}
	return nil
}

// Get return file{image, extension} by id and quality
func (store *localdb) Get(name, quality string) (file, error) {
	if v, ok := store.Storage[fileName{Name: name, Quality: quality}]; ok {
		return v, nil
	}
	for _, f := range store.Storage {
		log.Println(f)
	}
	return file{}, errors.New("wrong quality or id not exists")
}
