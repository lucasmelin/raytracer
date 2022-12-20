package display

import "github.com/lucasmelin/raytracer/internal/geometry"

// AABB represents an axis-aligned bounding box.
type AABB struct {
	Min geometry.Vec
	Max geometry.Vec
}

func NewAABB(min geometry.Vec, max geometry.Vec) *AABB {
	return &AABB{Min: min, Max: max}
}

func (ab *AABB) Hit(ray *geometry.Ray, dMin float64, dMax float64) bool {
	// Check X
	invD := 1 / ray.Direction.X
	d0 := (ab.Min.X - ray.Origin.X) * invD
	d1 := (ab.Max.X - ray.Origin.X) * invD
	if invD < 0 {
		d0, d1 = d1, d0
	}
	if d0 > dMin {
		dMin = d0
	}
	if d1 < dMax {
		dMax = d1
	}
	if dMax <= dMin {
		return false
	}
	// Check Y
	invD = 1 / ray.Direction.Y
	d0 = (ab.Min.Y - ray.Origin.Y) * invD
	d1 = (ab.Max.Y - ray.Origin.Y) * invD
	if invD < 0 {
		d0, d1 = d1, d0
	}
	if d0 > dMin {
		dMin = d0
	}
	if d1 < dMax {
		dMax = d1
	}
	if dMax <= dMin {
		return false
	}
	// Check Z
	invD = 1 / ray.Direction.Z
	d0 = (ab.Min.Z - ray.Origin.Z) * invD
	d1 = (ab.Max.Z - ray.Origin.Z) * invD
	if invD < 0 {
		d0, d1 = d1, d0
	}
	if d0 > dMin {
		dMin = d0
	}
	if d1 < dMax {
		dMax = d1
	}
	if dMax <= dMin {
		return false
	}

	return true
}

// Add combines two bounding boxes.
func (ab *AABB) Add(ab2 *AABB) *AABB {
	if ab2 == nil {
		return NewAABB(ab.Min, ab.Max)
	}
	return NewAABB(ab.Min.Min(ab2.Min), ab.Max.Max(ab2.Max))
}

func (a *AABB) Corners() []geometry.Vec {
	c := make([]geometry.Vec, 0, 8)
	for i := 0.0; i < 2; i++ {
		for j := 0.0; j < 2; j++ {
			for k := 0.0; k < 2; k++ {
				x := i*a.Min.X + (1-i)*a.Max.X
				y := j*a.Min.Y + (1-j)*a.Max.Y
				z := k*a.Min.Z + (1-k)*a.Max.Z
				c = append(c, geometry.NewVec(x, y, z))
			}
		}
	}
	return c
}

func (a *AABB) Extend(v geometry.Vec) *AABB {
	return NewAABB(a.Min.Min(v), a.Max.Max(v))
}
