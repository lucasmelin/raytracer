package display

import "github.com/lucasmelin/raytracer/internal/geometry"

// World holds a list of Hittables.
type World struct {
	Hittables []Hittable
}

// NewWorld created a new list of Hittables.
func NewWorld(h ...Hittable) *World {
	return &World{Hittables: h}
}

// Hit returns the first intersection between the Ray r and and of the Hittables in the World.
func (w *World) Hit(r geometry.Ray, tMin float64, tMax float64) (t float64, s Surfacer) {
	closest := tMax
	for _, h := range w.Hittables {
		if ht, hs := h.Hit(r, tMin, closest); ht > 0 {
			closest, t = ht, ht
			s = hs
		}
	}
	return t, s
}

// Add adds a Hittable to the World.
func (w *World) Add(h ...Hittable) int {
	w.Hittables = append(w.Hittables, h...)
	return len(w.Hittables)
}

// Surfacer represents something that can return surface normals and materials.
type Surfacer interface {
	Surface(p geometry.Vec) (n geometry.Unit, m Material)
}

// Material represents a material that scatters light.
type Material interface {
	Scatter(in geometry.Unit, n geometry.Unit) (geometry.Unit, Color, bool)
}

// Hittable represents a surface that can be Hit by a Ray.
type Hittable interface {
	Hit(ray geometry.Ray, tMin float64, tMax float64) (float64, Surfacer)
}
