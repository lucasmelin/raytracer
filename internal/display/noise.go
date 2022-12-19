package display

import (
	"math"

	"github.com/lucasmelin/raytracer/internal/geometry"
)

type Noise struct {
	Rnd   geometry.Rnd
	Scale float64
	per   Perlin
}

func NewNoise(rnd geometry.Rnd, scale float64) Noise {
	return Noise{Rnd: rnd, Scale: scale, per: NewPerlin(rnd)}
}

func (n Noise) At(u float64, v float64, p geometry.Vec) Color {
	scaleFactor := 0.5 * (1 + math.Sin(n.Scale*p.Z+10*n.per.turbulence(p, 7)))
	return NewColor(1, 1, 1).Scale(scaleFactor)
}
