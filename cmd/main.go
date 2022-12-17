package main

import (
	"os"

	"github.com/lucasmelin/raytracer/internal/display"
	"github.com/lucasmelin/raytracer/internal/geometry"
)

func main() {
	aspectRatio := 16.0 / 9.0
	imageWidth := 400
	imageHeight := int(float64(imageWidth) / aspectRatio)

	frame := display.Frame{Width: imageWidth, Height: imageHeight}
	sphere := geometry.NewSphere(geometry.NewVec(0, 0, -1), 0.5)
	frame.Render(os.Stdout, aspectRatio, sphere)
}
