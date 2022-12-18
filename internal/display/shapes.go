package display

import (
	"math"

	"github.com/lucasmelin/raytracer/internal/geometry"
)

// Sphere represents a sphere with a Center and a Radius.
type Sphere struct {
	Center   geometry.Vec
	Radius   float64
	Material Material
}

// NewSphere creates a new Sphere.
func NewSphere(center geometry.Vec, radius float64, material Material) Sphere {
	return Sphere{Center: center, Radius: radius, Material: material}
}

// Hit finds the first intersection between a ray and the sphere's surface.
func (s Sphere) Hit(r geometry.Ray, tMin, tMax float64) (float64, Surfacer) {
	oc := r.Origin.Sub(s.Center)
	a := r.Direction.Dot(r.Direction.Vec)
	halfb := oc.Dot(r.Direction.Vec)
	c := oc.Dot(oc) - s.Radius*s.Radius
	discriminant := halfb*halfb - a*c
	if discriminant <= 0 {
		return 0, s
	}
	// Find the nearest root that lies in the acceptable range.
	sqrt := math.Sqrt(discriminant)
	t := (-halfb - sqrt) / a
	if t > tMin && t < tMax {
		return t, s
	}
	t = (-halfb + sqrt) / a
	if t > tMin && t < tMax {
		return t, s
	}
	return 0, s
}

// Surface returns the normal and material at point p on the Sphere.
func (s Sphere) Surface(p geometry.Vec) (geometry.Unit, Material) {
	return p.Sub(s.Center).ToUnit(), s.Material
}
