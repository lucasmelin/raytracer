package geometry

// Sphere represents a sphere that can be rendered.
type Sphere struct {
	Center Vec
	Radius float64
}

// NewSphere creates a new Sphere.
func NewSphere(center Vec, radius float64) Sphere {
	return Sphere{Center: center, Radius: radius}
}
