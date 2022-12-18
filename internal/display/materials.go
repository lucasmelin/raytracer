package display

import (
	"math"

	"github.com/lucasmelin/raytracer/internal/geometry"
)

// Lambertian represents a Lambertian material attenuated by an Albedo.
type Lambertian struct {
	Albedo Color
}

// NewLambertian creates a new Lambertian material with a given color.
func NewLambertian(albedo Color) Lambertian {
	return Lambertian{Albedo: albedo}
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

// Dielectric represents a clear material.
type Dielectric struct {
	RefIndex float64
}

// NewDielectric creates a new material with a given index of refraction.
func NewDielectric(refIndex float64) Dielectric {
	return Dielectric{RefIndex: refIndex}
}

// Scatter reflects or refracts light rays based on the index of refraction.
func (d Dielectric) Scatter(in geometry.Ray, p geometry.Vec, n geometry.Unit) (geometry.Ray, Color, bool) {
	outNormal := n
	ratio := 1 / d.RefIndex
	if in.Direction.Dot(n) > 0 {
		outNormal = n.Inv()
		ratio = d.RefIndex
	}

	r, refracted := refract(in.Direction, outNormal, ratio)
	if !refracted {
		r = in.Direction.Reflect(n)
	}
	return geometry.NewRay(p, r), NewColor(1, 1, 1), true
}

func refract(u geometry.Unit, n geometry.Unit, ratio float64) (geometry.Unit, bool) {
	dt := u.Dot(n)
	disc := 1 - ratio*ratio*(1-dt*dt)
	if disc <= 0 {
		return geometry.Unit{}, false
	}
	u2 := (u.Sub(n.Scale(dt)).Scale(ratio)).Sub(n.Scale(math.Sqrt(disc))).ToUnit()
	return u2, true
}
