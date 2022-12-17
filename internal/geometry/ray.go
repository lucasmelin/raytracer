package geometry

// Ray defines a ray of light with an origin and direction.
type Ray struct {
	Origin    Vec
	Direction Unit
}

// NewRay creates a new ray with an origin and direction.
func NewRay(origin Vec, direction Unit) Ray {
	return Ray{Origin: origin, Direction: direction}
}

// At returns the ray at the given point.
func (r Ray) At(t float64) Vec {
	return r.Origin.Add(r.Direction.Scale(t))
}
