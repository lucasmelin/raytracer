package geometry

// Ray represents a ray defined by its origin and direction
type Ray struct {
	Origin    Vec
	Direction Unit
	Time      float64
	Rnd       Rnd
}

// NewRay creates a new ray with an origin and direction.
func NewRay(origin Vec, direction Unit, time float64, rnd Rnd) *Ray {
	return &Ray{
		Origin:    origin,
		Direction: direction,
		Time:      time,
		Rnd:       rnd,
	}
}

// At returns the ray at the given point.
func (r *Ray) At(t float64) Vec {
	return r.Origin.Add(r.Direction.Scale(t))
}
