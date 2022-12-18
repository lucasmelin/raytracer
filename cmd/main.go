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
	ground := geometry.NewSphere(geometry.NewVec(0, -100.5, -1), 100)
	world := display.NewWorld(sphere, ground)
	smoothness := 100
	frame.Render(os.Stdout, world, smoothness)
}
