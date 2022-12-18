package display

import (
	"fmt"
	"io"
	"math"

	"github.com/lucasmelin/raytracer/internal/geometry"
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
	return c.E[0]
}

// Green returns the color's second element.
func (c Color) Green() float64 {
	return c.E[1]
}

// Blue returns the color's third element.
func (c Color) Blue() float64 {
	return c.E[2]
}

// WriteColor writes the color to the writer.
func WriteColor(out io.Writer, c Color) {
	hueRed := toHue(c.Red())
	hueGreen := toHue(c.Green())
	hueBlue := toHue(c.Blue())

	fmt.Fprintln(out, hueRed, hueGreen, hueBlue)
}

// Add returns the sum of two colors.
func (c Color) Plus(c2 Color) Color {
	return Color{Vec: c.Vec.Add(c2.Vec)}
}

// Scale returns color scaled by a scalar.
func (c Color) Scale(n float64) Color {
	return Color{Vec: c.Vec.Scale(n)}
}

// Gamma raises each color channel to 1/n.
func (c Color) Gamma(n float64) Color {
	k := 1 / n
	return NewColor(math.Pow(c.Red(), k), math.Pow(c.Green(), k), math.Pow(c.Blue(), k))
}
