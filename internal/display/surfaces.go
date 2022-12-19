package display

import "github.com/lucasmelin/raytracer/internal/geometry"

// Hittable defines the interface of objects that can be hit by a ray
type Hittable interface {
	Hit(r *geometry.Ray, tMin float64, tMax float64) (bool, *HitRecord)
}

// World holds a list of Hittables.
type World struct {
	Hittables []Hittable
}

// NewWorld creates a new list of Hittables.
func NewWorld(h ...Hittable) *World {
	return &World{Hittables: h}
}

// Add adds a Hittable to the World.
func (w *World) Add(h ...Hittable) int {
	w.Hittables = append(w.Hittables, h...)
	return len(w.Hittables)
}

// Hit returns the first intersection between the Ray r and the Hittables in the World.
func (w *World) Hit(r *geometry.Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	var res *HitRecord
	hitSomething := false
	closest := tMax
	for _, h := range w.Hittables {
		if hit, hr := h.Hit(r, tMin, closest); hit {
			hitSomething = true
			res = hr
			closest = hr.t
		}
	}

	return hitSomething, res
}

type HitRecord struct {
	t        float64       // which t generated the hit
	p        geometry.Vec  // which point when hit
	normal   geometry.Unit // normal at that point
	Material Material      // the material associated to this record
}
