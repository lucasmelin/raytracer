package display

import "github.com/lucasmelin/raytracer/internal/geometry"

type Noise struct {
	Rnd geometry.Rnd
	per Perlin
}

func NewNoise(rnd geometry.Rnd) Noise {
	return Noise{Rnd: rnd, per: NewPerlin(rnd)}
}

func (n Noise) At(u float64, v float64, p geometry.Vec) Color {
	return NewColor(1, 1, 1).Scale(n.per.Generate(p))
}
