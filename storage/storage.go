package storage

import "image"

// Storage set interface of storages
type Storage interface {
	Add(image.Image, string) error
	Delete(string) error
	Get(string, string) (file, error)
	GetCurrentId() string
}

// NewStorage creates new storage
func NewStorage() Storage {
	return &localdb{Storage: make(map[fileName]file), CurrentId: 1}
}
