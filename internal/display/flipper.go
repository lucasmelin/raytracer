package display

import "github.com/lucasmelin/raytracer/internal/geometry"

// Flip contains a HitBoxer that can be flipped.
type Flip struct {
	Child HitBoxer
}

// NewFlip creates a new Flip.
func NewFlip(child HitBoxer) *Flip {
	return &Flip{Child: child}
}

// Hit calculates if the ray hits the HitBoxer. If so, the normal of the HitRecord is inverted.
func (f *Flip) Hit(r *geometry.Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	hit, record := f.Child.Hit(r, tMin, tMax)
	if hit {
		record.normal = record.normal.Inv()
	}
	return hit, record
}

// Box returns the bounding box that encloses the child surface.
func (f *Flip) Box(t0 float64, t1 float64) *AABB {
	return f.Child.Box(t0, t1)
}
