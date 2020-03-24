package main

import (
	"image"
	"os"

	_ "image/png"

	"github.com/faiface/pixel"
)

// loadPicture loads picture from file
func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
