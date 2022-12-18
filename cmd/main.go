package main

import (
	"math"
	"os"

	"github.com/lucasmelin/raytracer/internal/display"
	"github.com/lucasmelin/raytracer/internal/geometry"
)

func main() {
	aspectRatio := 16.0 / 9.0
	imageWidth := 400
	imageHeight := int(float64(imageWidth) / aspectRatio)

	window := display.NewWindow(imageWidth, imageHeight)

	materialLeft := display.NewLambertian(display.NewColor(0, 0, 1))
	materialRight := display.NewLambertian(display.NewColor(1, 0, 0))

	r := math.Cos(math.Pi / 4)
	world := display.NewWorld(
		display.NewSphere(geometry.NewVec(-r, 0, -1), r, materialLeft),
		display.NewSphere(geometry.NewVec(r, 0, -1), r, materialRight),
	)
	smoothness := 100
	window.Render(os.Stdout, world, smoothness)
}
