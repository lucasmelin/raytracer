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

	materialGround := display.Lambertian{Albedo: display.NewColor(0.8, 0.8, 0.0)}
	materialCenter := display.Lambertian{Albedo: display.NewColor(0.7, 0.3, 0.3)}
	materialLeft := display.Lambertian{Albedo: display.NewColor(0.8, 0.8, 0.8)}
	materialRight := display.Lambertian{Albedo: display.NewColor(0.8, 0.6, 0.2)}

	world := display.NewWorld(
		display.NewSphere(geometry.NewVec(0, -100.5, -1), 100, materialGround),
		display.NewSphere(geometry.NewVec(0, 0, -1), 0.5, materialCenter),
		display.NewSphere(geometry.NewVec(-1, 0, -1), 0.5, materialLeft),
		display.NewSphere(geometry.NewVec(1, 0, -1), 0.5, materialRight),
	)
	smoothness := 100
	frame.Render(os.Stdout, world, smoothness)
}
