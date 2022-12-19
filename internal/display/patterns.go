package display

import (
	"math"

	"github.com/lucasmelin/raytracer/internal/geometry"
)

type Checker struct {
	Size float64
	Odd  Texture
	Even Texture
}

type Solid struct {
	C Color
}

func NewChecker(size float64, t0 Texture, t1 Texture) Checker {
	return Checker{Size: size, Odd: t0, Even: t1}
}

func (c Checker) At(u float64, v float64, p geometry.Vec) Color {
	sines := math.Sin(c.Size*p.X) * math.Sin(c.Size*p.Y) * math.Sin(c.Size*p.Z)
	if sines < 0 {
		return c.Odd.At(u, v, p)
	}
	return c.Even.At(u, v, p)
}

func NewSolid(color Color) Solid {
	return Solid{C: color}
}

func (s Solid) At(u float64, v float64, p geometry.Vec) Color {
	return s.C
}
