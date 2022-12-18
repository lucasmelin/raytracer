package display

import "github.com/lucasmelin/raytracer/internal/geometry"

// Lambertian represents a Lambertian material attenuated by an Albedo.
type Lambertian struct {
	Albedo Color
}

// Scatter scatters light rays in a Lambertian pattern.
func (l Lambertian) Scatter(in geometry.Ray, p geometry.Vec, n geometry.Unit) (geometry.Ray, Color) {
	target := p.Add(n.Vec).Add(geometry.RandVecInSphere())
	out := geometry.NewRay(p, target.Sub(p).ToUnit())
	return out, l.Albedo
}
