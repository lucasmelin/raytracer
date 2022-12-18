package display

import "github.com/lucasmelin/raytracer/internal/geometry"

// Lambertian represents a Lambertian material attenuated by an Albedo.
type Lambertian struct {
	Albedo Color
}

// Scatter scatters light rays in a Lambertian pattern.
func (l Lambertian) Scatter(in geometry.Ray, p geometry.Vec, n geometry.Unit) (geometry.Ray, Color, bool) {
	target := p.Add(n.Vec).Add(geometry.RandVecInSphere())
	out := geometry.NewRay(p, target.Sub(p).ToUnit())
	return out, l.Albedo, true
}

// Metal represents a reflective material.
type Metal struct {
	Albedo Color
	Rough  float64
}

// NewMetal creates a new Metal material with a given color and roughness.
func NewMetal(albedo Color, roughness float64) Metal {
	return Metal{Albedo: albedo, Rough: roughness}
}

// Scatter reflects light rays.
func (m Metal) Scatter(in geometry.Ray, p geometry.Vec, n geometry.Unit) (geometry.Ray, Color, bool) {
	r := in.Direction.Reflect(n)
	dir := r.Add(geometry.RandVecInSphere().Scale(m.Rough)).ToUnit()
	out := geometry.NewRay(p, dir)
	ok := out.Direction.Dot(n) > 0
	return out, m.Albedo, ok
}
