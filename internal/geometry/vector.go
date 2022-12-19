package geometry

import "math"

// Unit represents a unit vector of length 1.
type Unit struct {
	Vec
}

// NewUnit creates a new Unit.
func NewUnit(e0, e1, e2 float64) Unit {
	return Unit{
		Vec: Vec{
			X: e0,
			Y: e1,
			Z: e2,
		},
	}
}

// Vec represents a 3-element vector.
type Vec struct {
	X float64
	Y float64
	Z float64
}

// NewVec creates a new Vec.
func NewVec(e0, e1, e2 float64) Vec {
	return Vec{e0, e1, e2}
}

// Inv returns the inverse of the vector as a new vector.
func (v Vec) Inv() Vec {
	return Vec{-v.X, -v.Y, -v.Z}
}

// Len returns the vector's length.
func (v Vec) Len() float64 {
	return math.Sqrt(v.LenSquared())
}

// LenSquared returns the square of the vector's length.
func (v Vec) LenSquared() float64 {
	xSq := v.X * v.X
	ySq := v.Y * v.Y
	zSq := v.Z * v.Z
	return xSq + ySq + zSq
}

// Add returns the sum of two vectors.
func (v Vec) Add(v2 Vec) Vec {
	return Vec{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
		Z: v.Z + v2.Z,
	}
}

// Sub returns the difference of two vectors.
func (v Vec) Sub(v2 Vec) Vec {
	return Vec{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
		Z: v.Z - v2.Z,
	}
}

// Mul returns the multiplication of two vectors.
func (v Vec) Mul(v2 Vec) Vec {
	return Vec{
		X: v.X * v2.X,
		Y: v.Y * v2.Y,
		Z: v.Z * v2.Z,
	}
}

// Div returns the division of two vectors.
func (v Vec) Div(v2 Vec) Vec {
	return Vec{
		X: v.X / v2.X,
		Y: v.Y / v2.Y,
		Z: v.Z / v2.Z,
	}
}

// Dot returns the dot product of two unit vectors.
func (v Vec) Dot(v2 Vec) float64 {
	newX := v.X * v2.X
	newY := v.Y * v2.Y
	newZ := v.Z * v2.Z
	return newX + newY + newZ
}

// Cross returns the cross-product of two vectors.
func (v Vec) Cross(v2 Vec) Vec {
	newX := v.Y*v2.Z - v.Z*v2.Y
	newY := v.Z*v2.X - v.X*v2.Z
	newZ := v.X*v2.Y - v.Y*v2.X
	return Vec{newX, newY, newZ}
}

// ToUnit converts the vector to a unit vector.
func (v Vec) ToUnit() Unit {
	scalingFactor := 1.0 / v.Len()
	newX := v.X * scalingFactor
	newY := v.Y * scalingFactor
	newZ := v.Z * scalingFactor
	return NewUnit(newX, newY, newZ)
}

// Scale returns the vector scaled by a scalar.
func (v Vec) Scale(n float64) Vec {
	newX := v.X * n
	newY := v.Y * n
	newZ := v.Z * n
	return Vec{X: newX, Y: newY, Z: newZ}
}

// Zero returns whether this is a zero vector.
func (v Vec) Zero() bool {
	return v.X == 0 && v.Y == 0 && v.Z == 0
}

// Dot returns the dot product of two unit vectors.
func (u Unit) Dot(u2 Unit) float64 {
	return u.X*u2.X + u.Y*u2.Y + u.Z*u2.Z
}

// Reflect reflects this unit vector about a normal vector n.
func (u Unit) Reflect(n Unit) Unit {
	return Unit{Vec: u.Sub(n.Scale(2 * u.Dot(n)))}
}

// Reflect reflects this unit vector about a normal vector n.
func (v Vec) Reflect(n Vec) Vec {
	return v.Sub(n.Scale(2 * v.Dot(n)))
}

// Inv returns the inverse of this unit vector as a new vector.
func (u Unit) Inv() Unit {
	return Unit{Vec: u.Vec.Inv()}
}

// Refract returns a refracted vector.
func Refract(u Unit, n Unit, ratio float64) (bool, *Unit) {
	dt := u.Dot(n)
	disc := 1 - ratio*ratio*(1-dt*dt)
	if disc <= 0 {
		return false, &Unit{}
	}
	u2 := (u.Sub(n.Scale(dt)).Scale(ratio)).Sub(n.Scale(math.Sqrt(disc))).ToUnit()
	return true, &u2
}

// Randgeometry.VecInSphere creates a random geometry.Vec within a unit sphere.
func RandVecInSphere(rnd Rnd) Vec {
	for {
		v := Vec{rnd.Float64(), rnd.Float64(), rnd.Float64()}.Scale(2).Sub(Vec{1, 1, 1})
		if v.LenSquared() < 1 {
			return v
		}
	}
}

// Randgeometry.VecInDisk creates a random geometry.Vec within a unit disk.
func RandVecInDisk(rnd Rnd) Vec {
	xy := Vec{1, 1, 0}
	for {
		v := Vec{rnd.Float64(), rnd.Float64(), 0}.Scale(2).Sub(xy)
		if v.Dot(v) < 1 {
			return v
		}
	}
}

// Min returns a new Vector using the smallest elements of two vectors.
func (v Vec) Min(v2 Vec) Vec {
	if v2.X < v.X {
		v.X = v2.X
	}
	if v2.Y < v.Y {
		v.Y = v2.Y
	}
	if v2.Z < v.Z {
		v.Z = v2.Z
	}
	return v
}

// Max returns a new Vector using the largest elements of two vectors.
func (v Vec) Max(v2 Vec) Vec {
	if v2.X > v.X {
		v.X = v2.X
	}
	if v2.Y > v.Y {
		v.Y = v2.Y
	}
	if v2.Z > v.Z {
		v.Z = v2.Z
	}
	return v
}

// RandUnit returns a random unit vector.
func RandUnit(rnd Rnd) Unit {
	return NewVec(2*rnd.Float64()-1, 2*rnd.Float64()-1, 2*rnd.Float64()-1).ToUnit()
}
