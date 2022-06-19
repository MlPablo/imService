package storage

import "image"

// fileName is struct that have id of image, and it's quality
type fileName struct {
	Name    string
	Quality string
}

// fileName is struct that stores image and it extension (png, jpg), and it's quality
type file struct {
	Image image.Image
	Ext   string
}
