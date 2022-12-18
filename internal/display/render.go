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

// Window collects the results of the ray traces on a Width by Height grid.
type Window struct {
	Width  int
	Height int
}

// NewWindow creates a new Window given a width and height.
func NewWindow(width int, height int) Window {
	return Window{Width: width, Height: height}
}

// camera contains a set of image coordinates.
type camera struct {
	origin     geometry.Vec
	horizontal geometry.Vec
	vertical   geometry.Vec
	lowerLeft  geometry.Vec
	u          geometry.Unit
	v          geometry.Unit
	w          geometry.Unit
	lensRadius float64
}

func newCamera(lookFrom geometry.Vec, lookAt geometry.Vec, verticalUp geometry.Unit, verticalFov float64, aspectRatio float64, aperture float64, focus float64) camera {
	theta := verticalFov * math.Pi / 180
	halfH := math.Tan(theta / 2)
	halfW := aspectRatio * halfH

	camera := camera{}

	camera.w = lookFrom.Sub(lookAt).ToUnit()
	camera.u = verticalUp.Cross(camera.w.Vec).ToUnit()
	camera.v = camera.w.Cross(camera.u.Vec).ToUnit()

	width := camera.u.Scale(halfW * focus)
	height := camera.v.Scale(halfH * focus)
	dist := camera.w.Scale(focus)

	camera.lensRadius = aperture / 2

	camera.origin = lookFrom
	camera.lowerLeft = camera.origin.Sub(width).Sub(height).Sub(dist)
	camera.horizontal = width.Scale(2)
	camera.vertical = height.Scale(2)

	return camera
}

// Ray returns a Ray passing through a given coordinate s, t.
func (c camera) Ray(s float64, t float64) geometry.Ray {
	rd := geometry.RandVecInDisk().Scale(c.lensRadius)
	offset := c.u.Scale(rd.X()).Add(c.v.Scale(rd.Y()))
	source := c.origin.Add(offset)
	dest := c.lowerLeft.Add(c.horizontal.Scale(s).Add(c.vertical.Scale(t)))
	return geometry.NewRay(
		source,
		dest.Sub(source).ToUnit(),
	)
}

func (w Window) Render(out io.Writer, h Hittable, samples int) {
	header := fmt.Sprintf("%s\n%d %d\n%d", asciiColorPalette, w.Width, w.Height, maxColor)
	fmt.Fprintln(out, header)

	// Camera
	from := geometry.NewVec(3, 3, 2)
	at := geometry.NewVec(0, 0, -1)
	focus := from.Sub(at).Len()
	cam := newCamera(from, at, geometry.NewUnit(0, 1, 0), 20, float64(w.Width)/float64(w.Height), 2, focus)

	fmt.Fprintf(os.Stderr, "Rendering image %d X %d", w.Width, w.Height)

	for y := w.Height - 1; y >= 0; y-- {
		fmt.Fprintf(os.Stderr, "\nScanlines remaining: %d", y)
		for x := 0; x < w.Width; x++ {
			c := NewColor(0, 0, 0)
			for s := 0; s < samples; s++ {
				u := (float64(x) + rand.Float64()) / float64(w.Width-1)
				v := (float64(y) + rand.Float64()) / float64(w.Height-1)
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

		scattered, attenuation, ok := m.Scatter(r.Direction, n)
		if !ok {
			return NewColor(0, 0, 0)
		}
		r2 := geometry.NewRay(p, scattered)
		return rayColor(r2, h, depth-1).Mul(attenuation)
	}
	t := 0.5 * (r.Direction.Y() + 1.0)
	white := NewColor(1.0, 1.0, 1.0).Scale(1 - t)
	blue := NewColor(0.5, 0.7, 1.0).Scale(t)
	return white.Plus(blue)
}
