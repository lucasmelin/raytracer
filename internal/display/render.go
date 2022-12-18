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
	height     float64
	width      float64
	origin     geometry.Vec
	horizontal geometry.Vec
	vertical   geometry.Vec
	lowerLeft  geometry.Vec
}

func newCamera(lookFrom geometry.Vec, lookAt geometry.Vec, verticalUp geometry.Unit, verticalFov float64, aspectRatio float64) camera {
	theta := verticalFov * math.Pi / 180
	halfH := math.Tan(theta / 2)
	halfW := aspectRatio * halfH

	w := lookFrom.Sub(lookAt).ToUnit()
	u := verticalUp.Cross(w.Vec).ToUnit()
	v := w.Cross(u.Vec).ToUnit()

	viewportHeight := 2 * halfH
	viewportWidth := 2 * halfW
	camera := camera{
		height:     viewportHeight,
		width:      aspectRatio * viewportHeight,
		origin:     lookFrom,
		horizontal: u.Scale(viewportWidth),
		vertical:   v.Scale(viewportHeight),
	}

	camera.lowerLeft = camera.origin.Sub(u.Scale(halfW)).Sub(v.Scale(halfH)).Sub(w.Vec)

	return camera
}

// Ray returns a Ray passing through a given coordinate u, v.
func (c camera) Ray(u float64, v float64) geometry.Ray {
	r := geometry.NewRay(
		c.origin,
		c.lowerLeft.Add((c.horizontal.Scale(u)).Add(c.vertical.Scale(v))).Sub(c.origin).ToUnit(),
	)
	return r
}

func (w Window) Render(out io.Writer, h Hittable, samples int) {
	header := fmt.Sprintf("%s\n%d %d\n%d", asciiColorPalette, w.Width, w.Height, maxColor)
	fmt.Fprintln(out, header)

	// Camera
	cam := newCamera(geometry.NewVec(-2, 2, 1), geometry.NewVec(0, 0, -1), geometry.NewUnit(0, 1, 0), 20, float64(w.Width)/float64(w.Height))

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
