package display

import "github.com/lucasmelin/raytracer/internal/geometry"

// Flip contains a HitBoxer that can be flipped.
type Flip struct {
	HitBoxer
}

// NewFlip creates a new Flip.
func NewFlip(child HitBoxer) *Flip {
	return &Flip{HitBoxer: child}
}

// Hit calculates if the ray hits the HitBoxer. If so, the normal of the HitRecord is inverted.
func (f *Flip) Hit(r *geometry.Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	hit, record := f.HitBoxer.Hit(r, tMin, tMax)
	if hit {
		record.normal = record.normal.Inv()
	}
	return hit, record
}
