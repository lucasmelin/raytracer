package display

import "github.com/lucasmelin/raytracer/internal/geometry"

type Texture interface {
	At(u float64, v float64, p geometry.Vec) Color
}
