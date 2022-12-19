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

// MovingSphere represents a sphere that moves from Center0 to Center1 over t0 to t1.
type MovingSphere struct {
	Center0  geometry.Vec
	Center1  geometry.Vec
	T0       float64
	T1       float64
	Radius   float64
	Material Material
}

// NewSphere creates a new Sphere.
func NewSphere(center geometry.Vec, radius float64, material Material) *Sphere {
	return &Sphere{Center: center, Radius: radius, Material: material}
}

// NewSphere creates a new Sphere with two centers separated by times t0 and t1.
func NewMovingSphere(center0 geometry.Vec, center1 geometry.Vec, t0 float64, t1 float64, radius float64, material Material) *MovingSphere {
	return &MovingSphere{
		Center0:  center0,
		Center1:  center1,
		T0:       t0,
		T1:       t1,
		Radius:   radius,
		Material: material,
	}
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
		u, v := s.UV(hitPoint, t)
		hr := HitRecord{
			t:        t,
			p:        hitPoint,
			normal:   hitPoint.Sub(s.Center).Scale(1 / s.Radius).ToUnit(),
			Material: s.Material,
			u:        u,
			v:        v,
		}
		return true, &hr
	}

	t = (-halfb + sqrt) / a
	if t > tMin && t < tMax {
		hitPoint := r.At(t)
		u, v := s.UV(hitPoint, t)
		hr := HitRecord{
			t:        t,
			p:        hitPoint,
			normal:   hitPoint.Sub(s.Center).Scale(1 / s.Radius).ToUnit(),
			Material: s.Material,
			u:        u,
			v:        v,
		}
		return true, &hr
	}
	return false, nil
}

// Hit finds the first intersection between a ray and the moving sphere's surface.
func (s MovingSphere) Hit(r *geometry.Ray, dMin float64, dMax float64) (bool, *HitRecord) {
	oc := r.Origin.Sub(s.Center(r.Time))
	a := r.Direction.Dot(r.Direction)
	b := oc.Dot(r.Direction.Vec)
	c := oc.Dot(oc) - s.Radius*s.Radius
	discriminant := b*b - a*c
	if discriminant <= 0 {
		return false, nil
	}
	// Find the nearest root that lies in the acceptable range.
	sqrt := math.Sqrt(discriminant)
	d := (-b - sqrt) / a
	if d > dMin && d < dMax {
		hitPoint := r.At(d)
		u, v := s.UV(hitPoint, d)
		hr := HitRecord{
			t:        d,
			p:        hitPoint,
			normal:   hitPoint.Sub(s.Center(r.Time)).Scale(1 / s.Radius).ToUnit(),
			Material: s.Material,
			u:        u,
			v:        v,
		}
		return true, &hr
	}

	d = (-b + sqrt) / a
	if d > dMin && d < dMax {
		hitPoint := r.At(d)
		u, v := s.UV(hitPoint, d)
		hr := HitRecord{
			t:        d,
			p:        hitPoint,
			normal:   hitPoint.Sub(s.Center(r.Time)).Scale(1 / s.Radius).ToUnit(),
			Material: s.Material,
			u:        u,
			v:        v,
		}
		return true, &hr
	}
	return false, nil
}

func (s *MovingSphere) Center(t float64) geometry.Vec {
	p := (t - s.T0) / (s.T1 - s.T0)
	offset := s.Center1.Sub(s.Center0).Scale(p)
	return s.Center0.Add(offset)
}

// Surface returns the normal and material at point p on the Sphere.
func (s *Sphere) Surface(p geometry.Vec) (geometry.Unit, Material) {
	return p.Sub(s.Center).Scale(s.Radius).ToUnit(), s.Material
}

func (s *Sphere) Box(t0 float64, t1 float64) *AABB {
	return NewAABB(
		s.Center.Sub(geometry.NewVec(s.Radius, s.Radius, s.Radius)),
		s.Center.Add(geometry.NewVec(s.Radius, s.Radius, s.Radius)),
	)
}

func (s *MovingSphere) Box(t0 float64, t1 float64) *AABB {
	box0 := NewAABB(
		s.Center(t0).Sub(geometry.NewVec(s.Radius, s.Radius, s.Radius)),
		s.Center(t0).Add(geometry.NewVec(s.Radius, s.Radius, s.Radius)),
	)
	box1 := NewAABB(
		s.Center(t1).Sub(geometry.NewVec(s.Radius, s.Radius, s.Radius)),
		s.Center(t1).Add(geometry.NewVec(s.Radius, s.Radius, s.Radius)),
	)
	return box0.Add(box1)
}

func (s *Sphere) UV(p geometry.Vec, t float64) (float64, float64) {
	p2 := p.Sub(s.Center).Scale(1 / s.Radius)
	phi := math.Atan2(p2.Z, p2.X)
	theta := math.Asin(p2.Y)
	u := 1 - (phi+math.Pi)/(2*math.Pi)
	v := (theta + math.Pi/2) / math.Pi
	return u, v
}

func (s *MovingSphere) UV(p geometry.Vec, t float64) (float64, float64) {
	p2 := p.Sub(s.Center(t)).Scale(1 / s.Radius)
	phi := math.Atan2(p2.Z, p2.X)
	theta := math.Asin(p2.Y)
	u := 1 - (phi+math.Pi)/(2*math.Pi)
	v := (theta + math.Pi/2) / math.Pi
	return u, v
}
