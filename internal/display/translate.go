package display

import "github.com/lucasmelin/raytracer/internal/geometry"

type Translate struct {
	Child  HitBoxer
	Offset geometry.Vec
}

func NewTranslate(child HitBoxer, offset geometry.Vec) *Translate {
	return &Translate{Child: child, Offset: offset}
}

func (t *Translate) Hit(r *geometry.Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	r2 := geometry.NewRay(r.Origin.Sub(t.Offset), r.Direction, r.Time, r.Rnd)
	didHit, record := t.Child.Hit(r2, tMin, tMax)
	if didHit {
		record.p = record.p.Add(t.Offset)
	}
	return didHit, record
}

func (t *Translate) Box(t0 float64, t1 float64) *AABB {
	box := t.Child.Box(t0, t1)
	return NewAABB(box.Min.Add(t.Offset), box.Max.Add(t.Offset))
}
