package display

import "github.com/lucasmelin/raytracer/internal/geometry"

type Noise struct {
	Rnd   geometry.Rnd
	Scale float64
	per   Perlin
}

func NewNoise(rnd geometry.Rnd, scale float64) Noise {
	return Noise{Rnd: rnd, Scale: scale, per: NewPerlin(rnd)}
}

func (n Noise) At(u float64, v float64, p geometry.Vec) Color {
	scaleFactor := n.per.turbulence(p.Scale(n.Scale), 7)
	return NewColor(1, 1, 1).Scale(scaleFactor)
}
