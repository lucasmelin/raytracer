package display

import (
	"math"

	"github.com/lucasmelin/raytracer/internal/geometry"
)

type RotateY struct {
	Child    HitBoxer
	sinTheta float64
	cosTheta float64
	box      *AABB
}

func NewRotateY(child HitBoxer, angle float64) *RotateY {
	radians := angle * math.Pi / 180
	ry := RotateY{
		Child:    child,
		sinTheta: math.Sin(radians),
		cosTheta: math.Cos(radians),
		box:      child.Box(0, 1),
	}
	for _, p := range ry.box.Corners() {
		p2 := ry.right(p)
		ry.box = ry.box.Extend(p2)
	}
	return &ry
}

func (r *RotateY) Box(t0 float64, t1 float64) *AABB {
	return r.box
}

func (r *RotateY) right(dir geometry.Vec) geometry.Vec {
	x := r.cosTheta*dir.X + r.sinTheta*dir.Z
	z := -r.sinTheta*dir.X + r.cosTheta*dir.Z
	return geometry.NewVec(x, dir.Y, z)
}

func (r *RotateY) left(dir geometry.Vec) geometry.Vec {
	x := r.cosTheta*dir.X - r.sinTheta*dir.Z
	z := r.sinTheta*dir.X + r.cosTheta*dir.Z
	return geometry.NewVec(x, dir.Y, z)
}

func (r *RotateY) Hit(ray *geometry.Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	ray2 := geometry.NewRay(r.left(ray.Origin), geometry.Unit{r.left(ray.Direction.Vec)}, ray.Time, ray.Rnd)
	didHit, record := r.Child.Hit(ray2, tMin, tMax)
	if didHit {
		record.normal = geometry.Unit{Vec: r.right(record.normal.Vec)}
		record.p = r.right(record.p)
	}
	return didHit, record
}
