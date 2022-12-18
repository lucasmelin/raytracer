package display

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"

	"github.com/lucasmelin/raytracer/internal/geometry"
)

const (
	asciiColorPalette = "P3"
	maxColor          = 255
	aspectRatio       = 16.0 / 9.0
	bias              = 0.001
	maxDepth          = 50
)

// Frame collects the results of the ray traces on a Width by Height grid.
type Frame struct {
	Width  int
	Height int
}

// camera contains a viewport and a focal length.
type camera struct {
	viewport
	focalLength float64
	aspectRatio float64
}

// viewport contains a set of image coordinates.
type viewport struct {
	height          float64
	width           float64
	origin          geometry.Vec
	horizontal      geometry.Vec
	vertical        geometry.Vec
	lowerLeftCorner geometry.Vec
}

func newCamera() camera {
	viewportHeight := 2.0
	viewportWidth := aspectRatio * viewportHeight
	camera := camera{
		viewport: viewport{
			height:     viewportHeight,
			width:      aspectRatio * viewportHeight,
			origin:     geometry.NewVec(0, 0, 0),
			horizontal: geometry.NewVec(viewportWidth, 0, 0),
			vertical:   geometry.NewVec(0, viewportHeight, 0),
		},
		focalLength: 1.0,
		aspectRatio: aspectRatio,
	}
	camera.lowerLeftCorner = camera.origin.Sub(camera.horizontal.Scale(0.5)).Sub(camera.vertical.Scale(0.5)).Sub(geometry.NewVec(0, 0, camera.focalLength))
	return camera
}

// Ray returns a Ray passing through a given coordinate u, v.
func (c camera) Ray(u float64, v float64) geometry.Ray {
	r := geometry.NewRay(
		c.origin,
		c.lowerLeftCorner.Add((c.horizontal.Scale(u)).Add(c.vertical.Scale(v))).ToUnit(),
	)
	return r
}

func (f Frame) Render(out io.Writer, h Hittable, samples int) {
	header := fmt.Sprintf("%s\n%d %d\n%d", asciiColorPalette, f.Width, f.Height, maxColor)
	fmt.Fprintln(out, header)

	// Camera
	cam := newCamera()

	fmt.Fprintf(os.Stderr, "Rendering image %d X %d", f.Width, f.Height)

	for y := f.Height - 1; y >= 0; y-- {
		fmt.Fprintf(os.Stderr, "\nScanlines remaining: %d", y)
		for x := 0; x < f.Width; x++ {
			c := NewColor(0, 0, 0)
			for s := 0; s < samples; s++ {
				u := (float64(x) + rand.Float64()) / float64(f.Width-1)
				v := (float64(y) + rand.Float64()) / float64(f.Height-1)
				r := cam.Ray(u, v)
				c = c.Plus(rayColor(r, h, maxDepth))
			}
			c = c.Scale(1.0 / float64(samples)).Gamma(2)
			WriteColor(out, c)
		}
	}
	fmt.Fprintf(os.Stderr, "\nDone\n")
}

func toHue(value float64) int {
	return int(256 * clamp(value, 0.0, 0.999))
}

func clamp(x, min, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

// rayColor linearly blends white and blue depending on the height of the Y coordinate.
func rayColor(r geometry.Ray, h Hittable, depth int) Color {
	// If we've exceeded the ray bounce limit, no more light is gathered.
	if depth <= 0 {
		return NewColor(0, 0, 0)
	}
	if t, s := h.Hit(r, bias, math.MaxFloat64); t > 0 {
		p := r.At(t)
		n, m := s.Surface(p)

		r2, attenuation := m.Scatter(r, p, n)
		if attenuation.Zero() {
			return attenuation
		}
		return rayColor(r2, h, depth-1).Mul(attenuation)
	}
	t := 0.5 * (r.Direction.Y() + 1.0)
	white := NewColor(1.0, 1.0, 1.0).Scale(1 - t)
	blue := NewColor(0.5, 0.7, 1.0).Scale(t)
	return white.Plus(blue)
}
