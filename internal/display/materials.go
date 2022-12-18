package display

import (
	"math"
	"math/rand"

	"github.com/lucasmelin/raytracer/internal/geometry"
)

// Lambertian represents a Lambertian material attenuated by an Albedo.
type Lambertian struct {
	Albedo Color
}

// NewLambertian creates a new Lambertian material with a given color.
func NewLambertian(albedo Color) Lambertian {
	return Lambertian{Albedo: albedo}
}

// Scatter scatters light rays in a Lambertian pattern.
func (l Lambertian) Scatter(in geometry.Unit, n geometry.Unit) (geometry.Unit, Color, bool) {
	out := n.Add(geometry.RandVecInSphere()).ToUnit()
	return out, l.Albedo, true
}

// Metal represents a reflective material.
type Metal struct {
	Albedo Color
	Rough  float64
}

// NewMetal creates a new Metal material with a given color and roughness.
func NewMetal(albedo Color, roughness float64) Metal {
	return Metal{Albedo: albedo, Rough: roughness}
}

// Scatter reflects light rays.
func (m Metal) Scatter(in geometry.Unit, n geometry.Unit) (geometry.Unit, Color, bool) {
	r := in.Reflect(n)
	out := r.Add(geometry.RandVecInSphere().Scale(m.Rough)).ToUnit()
	ok := out.Dot(n) > 0
	return out, m.Albedo, ok
}

// Dielectric represents a clear material.
type Dielectric struct {
	RefIndex float64
}

// NewDielectric creates a new material with a given index of refraction.
func NewDielectric(refIndex float64) Dielectric {
	return Dielectric{RefIndex: refIndex}
}

// Scatter reflects or refracts light rays based on the index of refraction.
func (d Dielectric) Scatter(in geometry.Unit, n geometry.Unit) (geometry.Unit, Color, bool) {
	outNormal := n
	ratio := 1 / d.RefIndex
	cosTheta := -in.Dot(n) / in.Len()

	if in.Dot(n) > 0 {
		outNormal = n.Inv()
		ratio = d.RefIndex
		cosTheta = d.RefIndex * in.Dot(n) / in.Len()
	}

	out, refracted := refract(in, outNormal, ratio)

	if !refracted || schlick(cosTheta, ratio) > rand.Float64() {
		out = in.Reflect(n)
	}
	return out, NewColor(1, 1, 1), true
}

func refract(u geometry.Unit, n geometry.Unit, ratio float64) (geometry.Unit, bool) {
	dt := u.Dot(n)
	disc := 1 - ratio*ratio*(1-dt*dt)
	if disc <= 0 {
		return geometry.Unit{}, false
	}
	u2 := (u.Sub(n.Scale(dt)).Scale(ratio)).Sub(n.Scale(math.Sqrt(disc))).ToUnit()
	return u2, true
}

func schlick(cos float64, refIndex float64) float64 {
	r := (1 - refIndex) / (1 + refIndex)
	r = r * r
	return r + (1-r)*math.Pow((1-cos), 5)
}
