package display

import (
	"math"

	"github.com/lucasmelin/raytracer/internal/geometry"
)

// Checker represents a 2-color checkerboard pattern.
type Checker struct {
	Size float64
	Odd  Texture
	Even Texture
}

// Solid represents a single solid color.
type Solid struct {
	Color Color
}

// NewChecker returns a new Checker.
func NewChecker(size float64, t0 Texture, t1 Texture) Checker {
	return Checker{Size: size, Odd: t0, Even: t1}
}

// At returns the Color of the ray at the given point along the checkerboard pattern.
func (c Checker) At(u float64, v float64, p geometry.Vec) Color {
	sines := math.Sin(c.Size*p.X) * math.Sin(c.Size*p.Y) * math.Sin(c.Size*p.Z)
	if sines < 0 {
		return c.Odd.At(u, v, p)
	}
	return c.Even.At(u, v, p)
}

// NewSolid returns a new Solid.
func NewSolid(color Color) Solid {
	return Solid{Color: color}
}

// At returns the Color of the ray.
func (s Solid) At(u float64, v float64, p geometry.Vec) Color {
	return s.Color
}
