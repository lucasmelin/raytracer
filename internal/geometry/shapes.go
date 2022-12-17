package geometry

import "math"

// Hittable represents a surface that can be Hit by a Ray.
type Hittable interface {
	Hit(ray Ray, tMin float64, tMax float64) (float64, Vec, Unit)
}

// Sphere represents a sphere with a Center and a Radius.
type Sphere struct {
	Center Vec
	Radius float64
}

// NewSphere creates a new Sphere.
func NewSphere(center Vec, radius float64) Sphere {
	return Sphere{Center: center, Radius: radius}
}

// Hit finds the first intersection between a ray and the sphere's surface.
func (s Sphere) Hit(r Ray, tMin, tMax float64) (t float64, p Vec, n Unit) {
	oc := r.Origin.Sub(s.Center)
	a := r.Direction.Dot(r.Direction.Vec)
	halfb := oc.Dot(r.Direction.Vec)
	c := oc.Dot(oc) - s.Radius*s.Radius
	discriminant := halfb*halfb - a*c
	if discriminant <= 0 {
		return 0, p, n
	}
	// Find the nearest root that lies in the acceptable range.
	sqrt := math.Sqrt(discriminant)
	t = (-halfb - sqrt) / a
	if t > tMin && t < tMax {
		p = r.At(t)
		return t, p, p.Sub(s.Center).ToUnit()
	}
	t = (-halfb + sqrt) / a
	if t > tMin && t < tMax {
		p = r.At(t)
		return t, p, p.Sub(s.Center).ToUnit()
	}
	return 0, p, n
}
