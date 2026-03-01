package infrastructure

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

type ImageWriter struct{}

func NewImageWriter() *ImageWriter {
	return &ImageWriter{}
}

func (iw *ImageWriter) SaveImage(img *image.RGBA, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	return nil
}
