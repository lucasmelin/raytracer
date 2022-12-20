package display

import (
	"math"

	"github.com/lucasmelin/raytracer/internal/geometry"
)

// Material represents a material that scatters light.
type Material interface {
	Scatter(r *geometry.Ray, rec *HitRecord) (wasScattered bool, attenuation *Color, scattered *geometry.Ray)
	Emit(rec *HitRecord) Color
}

type nonEmitter struct{}

func (n nonEmitter) Emit(rec *HitRecord) Color {
	return Black
}

// Lambertian represents a Lambertian material attenuated by an Albedo.
type Lambertian struct {
	Albedo Texture
	nonEmitter
}

// NewLambertian creates a new Lambertian material with a given color.
func NewLambertian(albedo Texture) Lambertian {
	return Lambertian{Albedo: albedo}
}

// Scatter scatters light rays in a Lambertian pattern.
func (l Lambertian) Scatter(r *geometry.Ray, rec *HitRecord) (bool, *Color, *geometry.Ray) {
	out := rec.normal.Add(geometry.RandVecInSphere(r.Rnd)).ToUnit()
	attenuation := l.Albedo.At(rec.u, rec.v, rec.p)
	return true, &attenuation, geometry.NewRay(rec.p, out.ToUnit(), r.Time, r.Rnd)
}

// Metal represents a reflective material.
type Metal struct {
	Albedo Color
	Rough  float64
	nonEmitter
}

// NewMetal creates a new Metal material with a given color and roughness.
func NewMetal(albedo Color, roughness float64) Metal {
	return Metal{Albedo: albedo, Rough: roughness}
}

// Scatter reflects light rays.
func (m Metal) Scatter(r *geometry.Ray, rec *HitRecord) (bool, *Color, *geometry.Ray) {
	reflected := r.Direction.ToUnit().Reflect(rec.normal)
	out := reflected.Add(geometry.RandVecInSphere(r.Rnd).Scale(m.Rough))
	return true, &m.Albedo, geometry.NewRay(rec.p, out.ToUnit(), r.Time, r.Rnd)
}

// Dielectric represents a clear material.
type Dielectric struct {
	RefIndex float64
	nonEmitter
}

// NewDielectric creates a new material with a given index of refraction.
func NewDielectric(refIndex float64) Dielectric {
	return Dielectric{RefIndex: refIndex}
}

// Scatter reflects or refracts light rays based on the index of refraction.
func (d Dielectric) Scatter(r *geometry.Ray, rec *HitRecord) (bool, *Color, *geometry.Ray) {
	in := r.Direction
	n := rec.normal

	outNormal := n
	ratio := 1 / d.RefIndex
	cosTheta := -in.Dot(n) / in.Len()

	if in.Dot(n) > 0 {
		outNormal = n.Inv()
		ratio = d.RefIndex
		cosTheta = d.RefIndex * in.Dot(n) / in.Len()
	}

	refracted, out := geometry.Refract(in, outNormal, ratio)

	if !refracted || schlick(cosTheta, ratio) > r.Rnd.Float64() {
		a := in.Reflect(n)
		out = &a
	}
	return true, &White, geometry.NewRay(rec.p, out.ToUnit(), r.Time, r.Rnd)
}

func schlick(cos float64, refIndex float64) float64 {
	r := (1 - refIndex) / (1 + refIndex)
	r = r * r
	return r + (1-r)*math.Pow((1-cos), 5)
}

type Light struct {
	Solid Solid
}

func NewLight(c Color) *Light {
	return &Light{Solid: NewSolid(c)}
}

func (l Light) Scatter(r *geometry.Ray, rec *HitRecord) (bool, *Color, *geometry.Ray) {
	return false, &Color{}, &geometry.Ray{}
}

func (l Light) Emit(rec *HitRecord) Color {
	return l.Solid.At(rec.u, rec.v, rec.p)
}

type Isotropic struct {
	albedo Solid
	Rnd    geometry.Rnd
	nonEmitter
}

func NewIsotropic(albedo Solid, rnd geometry.Rnd) *Isotropic {
	return &Isotropic{albedo: albedo, Rnd: rnd}
}

func (i *Isotropic) Scatter(r *geometry.Ray, rec *HitRecord) (bool, *Color, *geometry.Ray) {
	color := i.albedo.At(rec.u, rec.v, rec.p)
	return true, &color, r
}
