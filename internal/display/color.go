package display

import (
	"math"

	"github.com/lucasmelin/raytracer/internal/geometry"
)

var (
	White = NewColor(1.0, 1.0, 1.0)
	Black = Color{}
)

// Color represents an RGB color value.
type Color struct {
	geometry.Vec
}

func NewColor(e0, e1, e2 float64) Color {
	return Color{
		Vec: geometry.NewVec(e0, e1, e2),
	}
}

// Red returns the color's first element.
func (c Color) Red() float64 {
	return c.Vec.X
}

// Green returns the color's second element.
func (c Color) Green() float64 {
	return c.Vec.Y
}

// Blue returns the color's third element.
func (c Color) Blue() float64 {
	return c.Vec.Z
}

// Scale returns color scaled by a scalar.
func (c Color) Scale(n float64) Color {
	return Color{Vec: c.Vec.Scale(n)}
}

// Mul returns the multiplication of two colors.
func (c Color) Mult(c2 Color) Color {
	return Color{Vec: c.Vec.Mul(c2.Vec)}
}

// Add returns the sum of two colors.
func (c Color) Add(c2 Color) Color {
	return Color{Vec: c.Vec.Add(c2.Vec)}
}

// PixelValue converts a Color into a pixel value.
func (c Color) PixelValue() uint32 {
	r := uint32(math.Min(255.0, c.X*255.99))
	g := uint32(math.Min(255.0, c.Y*255.99))
	b := uint32(math.Min(255.0, c.Z*255.99))

	return ((r & 0xFF) << 16) | ((g & 0xFF) << 8) | (b & 0xFF)
}
