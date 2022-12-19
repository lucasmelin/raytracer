package display

import (
	"math"
	"sort"

	"github.com/lucasmelin/raytracer/internal/geometry"
)

// HitBoxer represents something that can be hit by a ray
type HitBoxer interface {
	Hit(r *geometry.Ray, tMin float64, tMax float64) (bool, *HitRecord)
	Box(t0 float64, t1 float64) *AABB
}

// List holds a list of Hittables.
type List struct {
	Hittables []HitBoxer
}

type BVH struct {
	Left  HitBoxer
	Right HitBoxer
	box   *AABB
}

// NewList creates a new list of Hittables.
func NewList(h ...HitBoxer) *List {
	return &List{Hittables: h}
}

// Add adds a Hittable to the List.
func (w *List) Add(h ...HitBoxer) int {
	w.Hittables = append(w.Hittables, h...)
	return len(w.Hittables)
}

// Hit returns the first intersection between the Ray r and the Hittables in the List.
func (w *List) Hit(r *geometry.Ray, tMin float64, tMax float64) (bool, *HitRecord) {
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

func (w *List) Box(t0 float64, t1 float64) *AABB {
	box := &AABB{}
	for _, h := range w.Hittables {
		box = h.Box(t0, t1).Add(box)
	}
	if box == nil {
		panic("No boxes defined")
	}
	return box
}

type HitRecord struct {
	t        float64       // which t generated the hit
	p        geometry.Vec  // which point when hit
	normal   geometry.Unit // normal at that point
	Material Material      // the material associated to this record
}

func NewBVH(depth int, time0 float64, time1 float64, h ...HitBoxer) *BVH {
	b := BVH{}
	switch len(h) {
	case 1:
		b.Left, b.Right = h[0], h[0]
		b.box = h[0].Box(time0, time1)
		return &b
	case 2:
		b.Left, b.Right = h[0], h[1]
		b.box = h[0].Box(time0, time1).Add(h[1].Box(time0, time1))
		return &b
	default:
		axis := depth % 3
		switch axis {
		case 0:
			sort.Slice(h, func(i, j int) bool {
				h0 := h[i].Box(time0, time1)
				h1 := h[j].Box(time0, time1)
				return h0.Min.X < h1.Min.X
			})
		case 1:
			sort.Slice(h, func(i, j int) bool {
				h0 := h[i].Box(time0, time1)
				h1 := h[j].Box(time0, time1)
				return h0.Min.Y < h1.Min.Y
			})
		case 2:
			sort.Slice(h, func(i, j int) bool {
				h0 := h[i].Box(time0, time1)
				h1 := h[j].Box(time0, time1)
				return h0.Min.Z < h1.Min.Z
			})
		}
		split := int(math.Floor(float64(len(h))/2 + 1))
		b.Left = NewBVH(depth+1, time0, time1, h[0:split]...)
		b.Right = NewBVH(depth+1, time0, time1, h[split:]...)
		b.box = b.Left.Box(time0, time1).Add(b.Right.Box(time0, time1))
		return &b
	}
}

func (b *BVH) Hit(ray *geometry.Ray, dMin float64, dMax float64) (bool, *HitRecord) {
	if !b.box.Hit(ray, dMin, dMax) {
		return false, nil
	}
	lDist, lBounce := b.Left.Hit(ray, dMin, dMax)
	rDist, rBounce := b.Right.Hit(ray, dMin, dMax)

	if lDist && rDist {
		if lBounce.t < rBounce.t {
			return lDist, lBounce
		}
		return rDist, rBounce
	}

	if lDist {
		return lDist, lBounce
	}
	if rDist {
		return rDist, rBounce
	}
	return false, nil
}

func (b *BVH) Box(t0 float64, t1 float64) *AABB {
	return b.box
}
