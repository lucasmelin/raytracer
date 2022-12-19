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
func NewSphere(center geometry.Vec, radius float64, material Material) *Sphere {
	return &Sphere{Center: center, Radius: radius, Material: material}
}

// Hit finds the first intersection between a ray and the sphere's surface.
func (s Sphere) Hit(r *geometry.Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	oc := r.Origin.Sub(s.Center)
	a := r.Direction.Dot(r.Direction)
	halfb := oc.Dot(r.Direction.Vec)
	c := oc.Dot(oc) - s.Radius*s.Radius
	discriminant := halfb*halfb - a*c
	if discriminant <= 0 {
		return false, nil
	}
	// Find the nearest root that lies in the acceptable range.
	sqrt := math.Sqrt(discriminant)
	t := (-halfb - sqrt) / a
	if t > tMin && t < tMax {
		hitPoint := r.At(t)
		hr := HitRecord{
			t:        t,
			p:        hitPoint,
			normal:   hitPoint.Sub(s.Center).Scale(1 / s.Radius).ToUnit(),
			Material: s.Material,
		}
		return true, &hr
	}

	t = (-halfb + sqrt) / a
	if t > tMin && t < tMax {
		hitPoint := r.At(t)
		hr := HitRecord{
			t:        t,
			p:        hitPoint,
			normal:   hitPoint.Sub(s.Center).Scale(1 / s.Radius).ToUnit(),
			Material: s.Material,
		}
		return true, &hr
	}
	return false, nil
}

// Surface returns the normal and material at point p on the Sphere.
func (s *Sphere) Surface(p geometry.Vec) (geometry.Unit, Material) {
	return p.Sub(s.Center).Scale(s.Radius).ToUnit(), s.Material
}
