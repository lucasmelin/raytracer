package main

import (
	"math/rand"
	"os"

	"github.com/lucasmelin/raytracer/internal/display"
	"github.com/lucasmelin/raytracer/internal/geometry"
)

func main() {
	window := display.NewWindow(1200, 800)

	ground := display.NewLambertian(display.NewColor(0.5, 0.5, 0.5))

	world := display.NewWorld(
		display.NewSphere(geometry.NewVec(0, -1000, 0), 1000, ground),
		display.NewSphere(geometry.NewVec(0, 1, 0), 1, display.NewDielectric(1.5)),
		display.NewSphere(geometry.NewVec(-4, 1, 0), 1, display.NewLambertian(display.NewColor(0.4, 0.2, 0.1))),
		display.NewSphere(geometry.NewVec(4, 1, 0), 1, display.NewMetal(display.NewColor(0.7, 0.6, 0.5), 0)),
	)

	for a := -11.0; a < 11; a++ {
		for b := -11.0; b < 11; b++ {
			center := geometry.NewVec(a+0.9*rand.Float64(), 0.2, b+0.9*rand.Float64())
			if center.Sub(geometry.NewVec(4, 0.2, 0)).Len() <= 0.9 {
				continue
			}
			world.Add(display.NewSphere(center, 0.2, material()))
		}
	}
	smoothness := 100
	window.Render(os.Stdout, world, smoothness)
}

func material() display.Material {
	m := rand.Float64()
	if m < 0.8 {
		// diffuse
		c := display.NewColor(rand.Float64()*rand.Float64(), rand.Float64()*rand.Float64(), rand.Float64()*rand.Float64())
		return display.NewLambertian(c)
	}
	if m < 0.95 {
		// metal
		c := display.NewColor(0.5*(1+rand.Float64()), 0.5*(1+rand.Float64()), 0.5*(1+rand.Float64()))
		return display.NewMetal(c, 0.5*rand.Float64())
	}
	//glass
	return display.NewDielectric(1.5)
}
