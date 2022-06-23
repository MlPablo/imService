package imageProccesor

import (
	"image"

	"github.com/nfnt/resize"
)

func Lancoz(image image.Image, quality int) image.Image {
	x := image.Bounds().Size().X
	Size := uint((x * quality) / 100)
	if Size < 100 {
		Size = 100
	}
	return resize.Resize(Size, 0, image, resize.Lanczos3)
}
