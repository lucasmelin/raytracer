package display

import "github.com/lucasmelin/raytracer/internal/geometry"

// Translate contains a HitBoxer that is translated according to a vector.
type Translate struct {
	Child  HitBoxer
	Offset geometry.Vec
}

// NewTranslate returns a new Translate.
func NewTranslate(child HitBoxer, offset geometry.Vec) *Translate {
	return &Translate{Child: child, Offset: offset}
}

// Hit calculates if the ray hits the HitBoxer. If so, the point of the HitRecord is translated.
func (t *Translate) Hit(r *geometry.Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	r2 := geometry.NewRay(r.Origin.Sub(t.Offset), r.Direction, r.Time, r.Rnd)
	didHit, record := t.Child.Hit(r2, tMin, tMax)
	if didHit {
		record.p = record.p.Add(t.Offset)
	}
	return didHit, record
}

// Box returns the translated bounding box for the Child HitBoxer.
func (t *Translate) Box(t0 float64, t1 float64) *AABB {
	box := t.Child.Box(t0, t1)
	return NewAABB(box.Min.Add(t.Offset), box.Max.Add(t.Offset))
}
