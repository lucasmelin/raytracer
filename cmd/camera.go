package main

import (
	"math"

	"github.com/lucasmelin/raytracer/internal/geometry"
)

// Camera is an interface that computes a ray.
// A random number source is provided for consistency, testing and performance benefits.
type Camera interface {
	ray(rnd geometry.Rnd, u, v float64) *geometry.Ray
}

// camera contains a set of image coordinates.
type camera struct {
	origin          geometry.Vec
	lowerLeftCorner geometry.Vec
	horizontal      geometry.Vec
	vertical        geometry.Vec
	u               geometry.Unit
	v               geometry.Unit
	lensRadius      float64
}

func NewCamera(lookFrom geometry.Vec, lookAt geometry.Vec, vup geometry.Vec, vfov float64, aspect float64, aperture float64, focusDist float64) Camera {
	theta := vfov * math.Pi / 180.0
	halfHeight := math.Tan(theta / 2.0)
	halfWidth := aspect * halfHeight

	origin := lookFrom
	w := lookFrom.Sub(lookAt).ToUnit()
	u := vup.Cross(w.Vec).ToUnit()
	v := w.Cross(u.Vec).ToUnit()

	lowerLeftCorner := origin.Add(u.Scale(-(halfWidth * focusDist))).Add(v.Scale(-(halfHeight * focusDist))).Add(w.Scale(-focusDist))
	horizontal := u.Scale(2 * halfWidth * focusDist)
	vertical := v.Scale(2 * halfHeight * focusDist)

	return camera{origin, lowerLeftCorner, horizontal, vertical, u, v, aperture / 2.0}
}

// ray returns a Ray that represents a ray of light.
func (c camera) ray(rnd geometry.Rnd, u float64, v float64) *geometry.Ray {
	rd := geometry.RandVecInDisk(rnd).Scale(c.lensRadius)
	offset := c.u.Scale(rd.X).Add(c.v.Scale(rd.Y))
	source := c.origin.Add(offset)
	dest := c.lowerLeftCorner.Add(c.horizontal.Scale(u).Add(c.vertical.Scale(v)))

	return geometry.NewRay(source, dest.Sub(source).ToUnit(), floatInRange(rnd, 0, 1), rnd)
}

func floatInRange(rnd geometry.Rnd, min float64, max float64) float64 {
	return min + rnd.Float64()*(max-min)
}
