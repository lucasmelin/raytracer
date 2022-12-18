package geometry

import (
	"fmt"
	"io"
	"math"
	"math/rand"
)

// Unit represents a unit vector of length 1.
type Unit struct {
	Vec
}

// NewUnit creates a new Unil.
func NewUnit(e0, e1, e2 float64) Unit {
	return Unit{
		Vec: NewVec(e0, e1, e2),
	}
}

// Vec represents a 3-element vector.
type Vec struct {
	E [3]float64
}

// NewVec creates a new Vec.
func NewVec(e0, e1, e2 float64) Vec {
	return Vec{
		E: [3]float64{e0, e1, e2},
	}
}

// X returns the vector's first element.
func (v Vec) X() float64 {
	return v.E[0]
}

// Y returns the vector's second element.
func (v Vec) Y() float64 {
	return v.E[1]
}

// Z returns the vector's third element.
func (v Vec) Z() float64 {
	return v.E[2]
}

// Inv returns the inverse of the vector as a new vector.
func (v Vec) Inv() Vec {
	return NewVec(-v.X(), -v.Y(), -v.Z())
}

// Len returns the vector's length.
func (v Vec) Len() float64 {
	return math.Sqrt(v.LenSquared())
}

// LenSquared returns the square of the vector's length.
func (v Vec) LenSquared() float64 {
	xSq := v.X() * v.X()
	ySq := v.Y() * v.Y()
	zSq := v.Z() * v.Z()
	return xSq + ySq + zSq
}

// StreamReader streams in space-separated vector elements from a Reader r.
func (v Vec) StreamReader(r io.Reader) error {
	_, err := fmt.Fscan(r, v.E[0], v.E[1], v.E[2])
	return err
}

// StreamWriter streams out space-separated vector elements to a Writer w.
func (v Vec) StreamWriter(w io.Writer) error {
	_, err := fmt.Fprint(w, v.E[0], v.E[1], v.E[2])
	return err
}

// Add returns the sum of two vectors.
func (v Vec) Add(v2 Vec) Vec {
	newX := v.X() + v2.X()
	newY := v.Y() + v2.Y()
	newZ := v.Z() + v2.Z()
	return NewVec(newX, newY, newZ)
}

// Sub returns the difference of two vectors.
func (v Vec) Sub(v2 Vec) Vec {
	newX := v.X() - v2.X()
	newY := v.Y() - v2.Y()
	newZ := v.Z() - v2.Z()
	return NewVec(newX, newY, newZ)
}

// Mul returns the multiplication of two vectors.
func (v Vec) Mul(v2 Vec) Vec {
	newX := v.X() * v2.X()
	newY := v.Y() * v2.Y()
	newZ := v.Z() * v2.Z()
	return NewVec(newX, newY, newZ)
}

// Div returns the division of two vectors.
func (v Vec) Div(v2 Vec) Vec {
	newX := v.X() / v2.X()
	newY := v.Y() / v2.Y()
	newZ := v.Z() / v2.Z()
	return NewVec(newX, newY, newZ)
}

// Dot returns the dot-product of two vectors.
func (v Vec) Dot(v2 Vec) float64 {
	newX := v.X() * v2.X()
	newY := v.Y() * v2.Y()
	newZ := v.Z() * v2.Z()
	return newX + newY + newZ
}

// Cross returns the cross-product of two vectors.
func (v Vec) Cross(v2 Vec) Vec {
	newX := v.Y()*v2.Z() - v.Z()*v2.Y()
	newY := v.Z()*v2.X() - v.X()*v2.Z()
	newZ := v.X()*v2.Y() - v.Y()*v2.X()
	return NewVec(newX, newY, newZ)
}

// ToUnit converts the vector to a unit vector.
func (v Vec) ToUnit() Unit {
	scalingFactor := 1.0 / v.Len()
	newX := v.X() * scalingFactor
	newY := v.Y() * scalingFactor
	newZ := v.Z() * scalingFactor
	return NewUnit(newX, newY, newZ)
}

// Scale returns the vector scaled by a scalar.
func (v Vec) Scale(n float64) Vec {
	newX := v.X() * n
	newY := v.Y() * n
	newZ := v.Z() * n
	return NewVec(newX, newY, newZ)
}

// Zero returns whether this is a zero vector.
func (v Vec) Zero() bool {
	return v.X() == 0 && v.Y() == 0 && v.Z() == 0
}

// Dot returns the dot product of two unit vectors.
func (u Unit) Dot(u2 Unit) float64 {
	return u.X()*u2.X() + u.Y()*u2.Y() + u.Z()*u2.Z()
}

// Reflect reflects this unit vector about a normal vector n.
func (u Unit) Reflect(n Unit) Unit {
	return Unit{Vec: u.Sub(n.Scale(2 * u.Dot(n)))}
}

// Inv returns the inverse of this unit vector as a new vector.
func (u Unit) Inv() Unit {
	return Unit{Vec: u.Vec.Inv()}
}

// RandVecInSphere creates a random Vec within a unit sphere.
func RandVecInSphere() Vec {
	for {
		v := NewVec(rand.Float64(), rand.Float64(), rand.Float64()).Scale(2).Sub(NewVec(1, 1, 1))
		if v.LenSquared() < 1 {
			return v
		}
	}
}
