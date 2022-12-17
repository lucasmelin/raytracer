package display

import "github.com/lucasmelin/raytracer/internal/geometry"

// World holds a list of Hittables.
type World struct {
	Hittables []geometry.Hittable
}

// NewWorld created a new list of Hittables.
func NewWorld(h ...geometry.Hittable) World {
	return World{Hittables: h}
}

// Hit returns the first intersection between the Ray r and and of the Hittables in the World.
func (w World) Hit(r geometry.Ray, tMin float64, tMax float64) (float64, geometry.Vec, geometry.Unit) {
	closest := tMax
	t := 0.0
	p := geometry.Vec{}
	n := geometry.Unit{}
	for _, h := range w.Hittables {
		if ht, hp, hn := h.Hit(r, tMin, closest); ht > 0 {
			closest, t = ht, ht
			p = hp
			n = hn
		}
	}
	return t, p, n
}
