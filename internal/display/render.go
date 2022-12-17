package display

import (
	"fmt"
	"io"
	"math"
	"os"

	"github.com/lucasmelin/raytracer/internal/geometry"
)

const asciiColorPalette = "P3"
const maxColor = 255

// Frame collects the results of the ray traces on a Width by Height grid.
type Frame struct {
	Width  int
	Height int
}

// camera contains a viewport and a focal length.
type camera struct {
	viewport
	focalLength float64
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

func newCamera(aspectRatio float64, viewportWidth float64, viewportHeight float64) camera {
	camera := camera{
		viewport: viewport{
			height:     viewportHeight,
			width:      viewportWidth,
			origin:     geometry.NewVec(0, 0, 0),
			horizontal: geometry.NewVec(viewportWidth, 0, 0),
			vertical:   geometry.NewVec(0, viewportHeight, 0),
		},
		focalLength: 1.0,
	}
	camera.lowerLeftCorner = camera.origin.Sub(camera.horizontal.Scale(0.5)).Sub(camera.vertical.Scale(0.5)).Sub(geometry.NewVec(0, 0, camera.focalLength))
	return camera
}

func (f Frame) Render(out io.Writer, aspectRatio float64, h geometry.Hittable) {
	header := fmt.Sprintf("%s\n%d %d\n%d", asciiColorPalette, f.Width, f.Height, maxColor)
	fmt.Fprintln(out, header)

	// Camera
	vh := 2.0
	vw := aspectRatio * vh
	cam := newCamera(aspectRatio, vw, vh)

	fmt.Fprintf(os.Stderr, "Rendering image %d X %d", f.Width, f.Height)

	for j := f.Height - 1; j >= 0; j-- {
		fmt.Fprintf(os.Stderr, "\nScanlines remaining: %d", j)
		for i := 0; i < f.Width; i++ {
			u := float64(i) / float64(f.Width-1)
			v := float64(j) / float64(f.Height-1)
			r := geometry.NewRay(
				cam.origin,
				cam.lowerLeftCorner.Add((cam.horizontal.Scale(u)).Add(cam.vertical.Scale(v))).ToUnit(),
			)
			c := rayColor(r, h)
			WriteColor(out, c)
		}
	}
	fmt.Fprintf(os.Stderr, "\nDone\n")
}

func toHue(value float64) int {
	return int(255.99 * value)
}

// rayColor linearly blends white and blue depending on the height of the Y coordinate.
func rayColor(r geometry.Ray, s geometry.Hittable) Color {
	if t, _, n := s.Hit(r, 0, math.MaxFloat64); t > 0 {
		return NewColor(n.X()+1, n.Y()+1, n.Z()+1).Scale(0.5)
	}
	t := 0.5 * (r.Direction.Y() + 1.0)
	white := NewColor(1.0, 1.0, 1.0).Scale(1 - t)
	blue := NewColor(0.5, 0.7, 1.0).Scale(t)
	return white.Plus(blue)
}
