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

	window := display.NewWindow(imageWidth, imageHeight)

	materialGround := display.NewLambertian(display.NewColor(0.8, 0.8, 0.0))
	materialCenter := display.NewLambertian(display.NewColor(0.1, 0.2, 0.5))
	materialLeft := display.NewDielectric(1.5)
	materialRight := display.NewMetal(display.NewColor(0.8, 0.6, 0.2), 0)

	world := display.NewWorld(
		display.NewSphere(geometry.NewVec(0, -100.5, -1), 100, materialGround),
		display.NewSphere(geometry.NewVec(0, 0, -1), 0.5, materialCenter),
		display.NewSphere(geometry.NewVec(-1, 0, -1), 0.5, materialLeft),
		display.NewSphere(geometry.NewVec(-1, 0, -1), -0.45, materialLeft),
		display.NewSphere(geometry.NewVec(1, 0, -1), 0.5, materialRight),
	)
	smoothness := 100
	window.Render(os.Stdout, world, smoothness)
}
