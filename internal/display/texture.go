package display

import (
	"math"

	"github.com/lucasmelin/raytracer/internal/geometry"
)

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

func (per Perlin) GenerateTrilinear(p geometry.Vec) float64 {
	u := p.X - math.Floor(p.X)
	v := p.Y - math.Floor(p.Y)
	w := p.Z - math.Floor(p.Z)

	i := int(math.Floor(p.X))
	j := int(math.Floor(p.Y))
	k := int(math.Floor(p.Z))
	c := make([]float64, 8)

	for di := 0; di < 2; di++ {
		for dj := 0; dj < 2; dj++ {
			for dk := 0; dk < 2; dk++ {
				x := per.permX[(i+di)&255]
				y := per.permX[(j+dj)&255]
				z := per.permX[(k+dk)&255]
				c[4*di+2*dj+dk] = per.rndFloat[x^y^z]
			}
		}
	}
	return trilinear(c, u, v, w)
}

func trilinear(c []float64, u float64, v float64, w float64) float64 {
	var sum float64
	for i := 0.0; i < 2; i++ {
		for j := 0.0; j < 2; j++ {
			for k := 0.0; k < 2; k++ {
				xyz := c[4*int(i)+2*int(j)+int(k)]
				sum += (i*u + (1-i)*(1-u)) * (j*v + (1-j)*(1-v)) *
					(k*w + (1-k)*(1-w)) *
					xyz
			}
		}
	}
	return sum
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
