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

func Render(out io.Writer, width int, height int, aspectRatio float64) {
	header := fmt.Sprintf("%s\n%d %d\n%d", asciiColorPalette, width, height, maxColor)
	fmt.Fprintln(out, header)

	// Camera
	vh := 2.0
	vw := aspectRatio * vh
	cam := newCamera(aspectRatio, vw, vh)

	fmt.Fprintf(os.Stderr, "Rendering image %d X %d", width, height)

	for j := height - 1; j >= 0; j-- {
		fmt.Fprintf(os.Stderr, "\nScanlines remaining: %d", j)
		for i := 0; i < width; i++ {
			u := float64(i) / float64(width-1)
			v := float64(j) / float64(height-1)
			r := geometry.NewRay(
				cam.origin,
				cam.lowerLeftCorner.Add((cam.horizontal.Scale(u)).Add(cam.vertical.Scale(v))).ToUnit(),
			)
			c := rayColor(r)
			WriteColor(out, c)
		}
	}
	fmt.Fprintf(os.Stderr, "\nDone\n")
}

func toHue(value float64) int {
	return int(255.99 * value)
}

// rayColor linearly blends white and blue depending on the height of the Y coordinate.
func rayColor(r geometry.Ray) Color {
	s := geometry.NewSphere(geometry.NewVec(0, 0, -1), 0.5)
	if t, ok := hitSphere(s, r); ok {
		n := r.At(t).Sub(s.Center).ToUnit()
		return NewColor(n.X()+1, n.Y()+1, n.Z()+1).Scale(0.5)
	}
	t := 0.5 * (r.Direction.Y() + 1.0)
	white := NewColor(1.0, 1.0, 1.0).Scale(1 - t)
	blue := NewColor(0.5, 0.7, 1.0).Scale(t)
	return white.Plus(blue)
}

func hitSphere(sphere geometry.Sphere, r geometry.Ray) (float64, bool) {
	oc := r.Origin.Sub(sphere.Center)
	a := r.Direction.Dot(r.Direction.Vec)
	b := 2.0 * oc.Dot(r.Direction.Vec)
	c := oc.Dot(oc) - sphere.Radius*sphere.Radius
	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return -1.0, false
	} else {
		return (-b - math.Sqrt(discriminant)) / (2.0 * a), true
	}
}
