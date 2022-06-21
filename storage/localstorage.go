package storage

import (
	"errors"
	"github.com/nfnt/resize"
	"image"
	"strconv"
)

// localdb struct of local DataBase
type localdb struct {
	Storage   map[fileName]file
	CurrentId int
}

// Add is method for localdb. It takes an image and creates 4 of them with different resize (25, 50, 75, 100) and adds it to localdb
// Stores images in format "id_quality.format"
func (store *localdb) Add(value image.Image, ext string) error {
	x := value.Bounds().Size().X
	for res := 25; res <= 100; res += 25 {
		// Resize Image
		Size := uint((x * res) / 100)
		if Size < 100 {
			Size = 100
		}
		file := file{resize.Resize(Size, 0, value, resize.Lanczos3), ext}
		// Add to storage resized image {id: id, quality: quality}
		store.Storage[fileName{Name: strconv.Itoa(store.CurrentId), Quality: strconv.Itoa(res)}] = file
	}
	store.CurrentId++
	return nil
}

// Delete deletes image by id (so deletes 4 images cause one image have fore different qualities)
func (store *localdb) Delete(str string) error {
	if _, ok := store.Storage[fileName{Name: str, Quality: "100"}]; ok {
		for i := 25; i <= 100; i += 25 {
			delete(store.Storage, fileName{Name: str, Quality: strconv.Itoa(i)})
		}
		return nil
	}
	return errors.New("not exists")
}

// Get return file{image, extension} by id and quality
func (store *localdb) Get(name, quality string) (file, error) {
	if v, ok := store.Storage[fileName{Name: name, Quality: quality}]; ok {
		return v, nil
	}
	return file{}, errors.New("wrong quality or id not exists")
}

// GetCurrentId return current ID of images( it used for telling users their ID's )
func (store *localdb) GetCurrentId() string {
	return strconv.Itoa(store.CurrentId)
}
