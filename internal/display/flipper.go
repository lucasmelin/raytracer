package display

import "github.com/lucasmelin/raytracer/internal/geometry"

type Flip struct {
	HitBoxer
}

func NewFlip(child HitBoxer) *Flip {
	return &Flip{HitBoxer: child}
}

func (f *Flip) Hit(r *geometry.Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	hit, record := f.HitBoxer.Hit(r, tMin, tMax)
	if hit {
		record.normal = record.normal.Inv()
	}
	return hit, record
}
