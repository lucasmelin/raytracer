package main

import (
	"os"

	"github.com/lucasmelin/raytracer/internal/display"
)

func main() {
	aspectRatio := 16.0 / 9.0
	imageWidth := 400
	imageHeight := int(float64(imageWidth) / aspectRatio)

	display.Render(os.Stdout, imageWidth, imageHeight, aspectRatio)
}
