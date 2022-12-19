package display

import (
	"image"
	"io"
	"math"

	_ "image/jpeg"

	"github.com/lucasmelin/raytracer/internal/geometry"
)

type Texture interface {
	At(u float64, v float64, p geometry.Vec) Color
}

type Image struct {
	X    int
	Y    int
	Data image.Image
}

type Perlin struct {
	Rnd     geometry.Rnd
	rndUnit []geometry.Unit
	permX   []int
	permY   []int
	permZ   []int
}

func NewPerlin(rnd geometry.Rnd) Perlin {
	return Perlin{
		Rnd:     rnd,
		rndUnit: perlinGen(rnd),
		permX:   perlinGenPerm(rnd),
		permY:   perlinGenPerm(rnd),
		permZ:   perlinGenPerm(rnd),
	}
}

func (per Perlin) GenerateTrilinear(p geometry.Vec) float64 {
	u := p.X - math.Floor(p.X)
	v := p.Y - math.Floor(p.Y)
	w := p.Z - math.Floor(p.Z)

	i := int(math.Floor(p.X))
	j := int(math.Floor(p.Y))
	k := int(math.Floor(p.Z))
	c := make([]geometry.Unit, 8)

	for di := 0; di < 2; di++ {
		for dj := 0; dj < 2; dj++ {
			for dk := 0; dk < 2; dk++ {
				x := per.permX[(i+di)&255]
				y := per.permX[(j+dj)&255]
				z := per.permX[(k+dk)&255]
				c[4*di+2*dj+dk] = per.rndUnit[x^y^z]
			}
		}
	}
	return interp(c, u, v, w)
}

func interp(c []geometry.Unit, u float64, v float64, w float64) float64 {
	u2 := u * u * (3 - 2*u)
	v2 := v * v * (3 - 2*v)
	w2 := w * w * (3 - 2*w)

	var sum float64
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				weight := geometry.NewVec(u-float64(i), v-float64(j), w-float64(k))
				xyz := c[4*i+2*j+k]

				sum += (float64(i)*u2 + (1-float64(i))*(1-u2)) * (float64(j)*v2 + (1-float64(j))*(1-v2)) * (float64(k)*w2 + (1-float64(k))*(1-w2)) * xyz.Vec.Dot(weight)
			}
		}
	}
	return sum
}

func perlinGen(rnd geometry.Rnd) []geometry.Unit {
	p := make([]geometry.Unit, 256)
	for i := 0; i < len(p); i++ {
		p[i] = geometry.RandUnit(rnd)
	}
	return p
}

func perlinPermute(rnd geometry.Rnd, p []int, n int) []int {
	for i := n - 1; i > 0; i-- {
		target := int(rnd.Float64() * float64(i+1))
		p[i], p[target] = p[target], p[i]
	}
	return p
}

func perlinGenPerm(rnd geometry.Rnd) []int {
	p := make([]int, 256)
	for i := 0; i < len(p); i++ {
		p[i] = i
	}
	p = perlinPermute(rnd, p, 256)
	return p
}

func (per Perlin) turbulence(p geometry.Vec, depth int) float64 {
	var sum float64
	p2 := p
	weight := 1.0
	for i := 0; i < depth; i++ {
		sum += weight * per.GenerateTrilinear(p2)
		weight *= 0.5
		p2 = p2.Scale(2)
	}
	return math.Abs(sum)
}

func NewImage(rc io.ReadCloser) (*Image, error) {
	defer rc.Close()
	im, _, err := image.Decode(rc)
	if err != nil {
		return nil, err
	}
	bounds := im.Bounds()
	i := Image{
		X:    bounds.Max.X,
		Y:    bounds.Max.Y,
		Data: im,
	}
	return &i, nil
}

func (i *Image) At(u float64, v float64, p geometry.Vec) Color {
	x := int(u * float64(i.X))
	y := int((1 - v) * float64(i.Y))
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x > i.X-1 {
		x = i.X - 1
	}
	if y > i.Y-1 {
		y = i.Y - 1
	}
	c := i.Data.At(x, y)
	r, g, b, _ := c.RGBA()
	divisor := float64(65535)
	return NewColor(float64(r)/divisor, float64(g)/divisor, float64(b)/divisor)
}
