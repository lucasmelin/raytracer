package display

import "github.com/lucasmelin/raytracer/internal/geometry"

type Texture interface {
	At(u float64, v float64, p geometry.Vec) Color
}

type Perlin struct {
	Rnd      geometry.Rnd
	rndFloat []float64
	permX    []int
	permY    []int
	permZ    []int
}

func NewPerlin(rnd geometry.Rnd) Perlin {
	return Perlin{
		Rnd:      rnd,
		rndFloat: perlinGen(rnd),
		permX:    perlinGenPerm(rnd),
		permY:    perlinGenPerm(rnd),
		permZ:    perlinGenPerm(rnd),
	}
}

func (per Perlin) Generate(p geometry.Vec) float64 {
	i := int(4*p.X) & 255
	j := int(4*p.Y) & 255
	k := int(4*p.Z) & 255
	return per.rndFloat[per.permX[i]^per.permY[j]^per.permZ[k]]
}

func perlinGen(rnd geometry.Rnd) []float64 {
	p := make([]float64, 256)
	for i := 0; i < len(p); i++ {
		p[i] = rnd.Float64()
	}
	return p
}

func perlinPermute(rnd geometry.Rnd, p []int, n int) []int {
	for i := n - 1; i > 0; i-- {
		target := int(rnd.Float64() * float64(i+1))
		p[i], p[target] = p[target], p[i]
	}
	return p
}

func perlinGenPerm(rnd geometry.Rnd) []int {
	p := make([]int, 256)
	for i := 0; i < len(p); i++ {
		p[i] = i
	}
	p = perlinPermute(rnd, p, 256)
	return p
}
