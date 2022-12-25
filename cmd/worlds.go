package main

import (
	"math/rand"
	"os"

	"github.com/lucasmelin/raytracer/internal/display"
	"github.com/lucasmelin/raytracer/internal/geometry"
)

// buildWeekOneWorld sets up the world and camera for the cover of the
// Ray Tracing the Next Week book.
func buildWeekOneWorld(width int, height int) (cameraSensor, *display.BVH) {
	world := display.List{}
	rnd := rand.New(rand.NewSource(rand.Int63()))
	w := 100.0

	ground := display.NewLambertian(display.NewSolid(display.NewColor(0.48, 0.83, 0.53)))
	groundList := display.NewList()
	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			min := geometry.NewVec(-1000+float64(i)*w, 0, -1000+float64(j)*w)
			max := geometry.NewVec(w, 1+99*rnd.Float64(), w).Add(min)
			groundList.Add(display.NewBlock(min, max, ground))
		}
	}

	world.Add(display.NewBVH(0, 0, 1, groundList.Hittables...))
	world.Add(display.NewRectangle(geometry.NewVec(123, 554, 147), geometry.NewVec(423, 554, 412), display.NewLight(display.NewColor(7, 7, 7))))
	center := geometry.NewVec(400, 400, 200)
	world.Add(display.NewMovingSphere(center, center.Add(geometry.NewVec(30, 0, 0)), 0, 1, 50, display.NewLambertian(display.NewSolid(display.NewColor(0.7, 0.3, 0.1)))))
	world.Add(display.NewSphere(geometry.NewVec(260, 150, 45), 50, display.NewDielectric(1.5)))
	world.Add(display.NewSphere(geometry.NewVec(0, 150, 145), 50, display.NewMetal(display.NewColor(0.8, 0.8, 0.9), 1)))
	boundary := display.NewSphere(geometry.NewVec(360, 150, 145), 70, display.NewDielectric(1.5))
	world.Add(boundary)
	world.Add(display.NewVolume(boundary, 0.2, display.NewIsotropic(display.NewSolid(display.NewColor(0.2, 0.4, 0.9)), rnd)))
	boundary = display.NewSphere(geometry.NewVec(0, 0, 0), 5000, display.NewDielectric(1.5))
	world.Add(display.NewVolume(boundary, 0.0001, display.NewIsotropic(display.NewSolid(display.NewColor(1, 1, 1)), rnd)))
	f, err := os.Open("assets/earthtwo.jpeg")
	if err != nil {
		panic(err)
	}
	earth, err := display.NewImage(f)
	if err != nil {
		panic(err)
	}

	world.Add(display.NewSphere(geometry.NewVec(400, 200, 400), 100, display.NewLambertian(earth)))
	perlin := display.NewNoise(rnd, 0.1)
	world.Add(display.NewSphere(geometry.NewVec(220, 280, 300), 80, display.NewLambertian(perlin)))
	sphereList := display.NewList()
	white := display.NewLambertian(display.NewSolid(display.NewColor(0.73, 0.73, 0.73)))
	for i := 0; i < 1000; i++ {
		sphereList.Add(display.NewSphere(geometry.NewVec(165*rnd.Float64(), 165*rnd.Float64(), 165*rnd.Float64()), 10, white))
	}
	world.Add(display.NewTranslate(display.NewRotateY(display.NewBVH(0, 0, 1, sphereList.Hittables...), 15), geometry.NewVec(-100, 270, 395)))

	lookAt := geometry.NewVec(278, 278, 0)
	lookFrom := geometry.NewVec(478, 278, -600)
	aperture := 0.0
	distToFocus := 10.0
	camera := newCamera(
		lookFrom,
		lookAt,
		geometry.NewVec(0, 1.0, 0),
		40,
		float64(width)/float64(height),
		aperture,
		distToFocus,
	)
	return camera, display.NewBVH(0, 0, 1, world.Hittables...)
}

// cornell is a simple Cornell box scene with two blocks made of smoke and fog.
func cornellSmoke(width int, height int) (cameraSensor, *display.BVH) {
	world := display.List{}
	rnd := rand.New(rand.NewSource(rand.Int63()))
	green := display.NewLambertian(display.NewSolid(display.NewColor(0.12, 0.45, 0.15)))
	red := display.NewLambertian(display.NewSolid(display.NewColor(0.65, 0.05, 0.05)))
	white := display.NewLambertian(display.NewSolid(display.NewColor(0.73, 0.73, 0.73)))
	light := display.NewLight(display.NewColor(15, 15, 15))
	smoke := display.NewIsotropic(display.NewSolid(display.NewColor(0, 0, 0)), rnd)
	fog := display.NewIsotropic(display.NewSolid(display.NewColor(1, 1, 1)), rnd)

	world.Hittables = append(world.Hittables,
		display.NewFlip(display.NewRectangle(
			geometry.NewVec(555, 0, 0), geometry.NewVec(555, 555, 555), green,
		)),
		display.NewRectangle(
			geometry.NewVec(0, 0, 0), geometry.NewVec(0, 555, 555), red,
		),
		display.NewRectangle(
			geometry.NewVec(113, 554, 127), geometry.NewVec(443, 554, 432), light,
		),
		display.NewRectangle(
			geometry.NewVec(0, 0, 0), geometry.NewVec(555, 0, 555), white,
		),
		display.NewFlip(display.NewRectangle(
			geometry.NewVec(0, 0, 555), geometry.NewVec(555, 555, 555), white,
		)),
		display.NewFlip(display.NewRectangle(
			geometry.NewVec(0, 555, 0), geometry.NewVec(555, 555, 555), white,
		)),

		display.NewVolume(display.NewTranslate(display.NewRotateY(display.NewBlock(geometry.NewVec(0, 0, 0), geometry.NewVec(165, 165, 165), white), -18), geometry.NewVec(130, 0, 65)), 0.01, fog),

		display.NewVolume(display.NewTranslate(display.NewRotateY(display.NewBlock(geometry.NewVec(0, 0, 0), geometry.NewVec(165, 330, 165), white), 15), geometry.NewVec(265, 0, 295)), 0.01, smoke),
	)

	lookAt := geometry.NewVec(278, 278, 0)
	lookFrom := geometry.NewVec(278, 278, -800)
	aperture := 0.1
	distToFocus := 10.0
	camera := newCamera(
		lookFrom,
		lookAt,
		geometry.NewVec(0, 1.0, 0),
		40,
		float64(width)/float64(height),
		aperture,
		distToFocus,
	)
	return camera, display.NewBVH(0, 0, 1, world.Hittables...)
}

// cornell is a simple Cornell box scene with two blocks.
func cornell(width int, height int) (cameraSensor, *display.BVH) {
	world := display.List{}
	green := display.NewLambertian(display.NewSolid(display.NewColor(0.12, 0.45, 0.15)))
	red := display.NewLambertian(display.NewSolid(display.NewColor(0.65, 0.05, 0.05)))
	white := display.NewLambertian(display.NewSolid(display.NewColor(0.73, 0.73, 0.73)))
	world.Hittables = append(world.Hittables,
		display.NewFlip(display.NewRectangle(
			geometry.NewVec(555, 0, 0), geometry.NewVec(555, 555, 555), green,
		)),
		display.NewRectangle(
			geometry.NewVec(0, 0, 0), geometry.NewVec(0, 555, 555), red,
		),
		display.NewRectangle(
			geometry.NewVec(213, 554, 227), geometry.NewVec(343, 554, 332), display.NewLight(display.NewColor(10, 10, 10)),
		),
		display.NewRectangle(
			geometry.NewVec(0, 0, 0), geometry.NewVec(555, 0, 555), white,
		),
		display.NewFlip(display.NewRectangle(
			geometry.NewVec(0, 0, 555), geometry.NewVec(555, 555, 555), white,
		)),
		display.NewFlip(display.NewRectangle(
			geometry.NewVec(0, 555, 0), geometry.NewVec(555, 555, 555), white,
		)),
		display.NewTranslate(display.NewRotateY(display.NewBlock(geometry.NewVec(0, 0, 0), geometry.NewVec(165, 165, 165), white), -18), geometry.NewVec(130, 0, 65)),
		display.NewTranslate(display.NewRotateY(display.NewBlock(geometry.NewVec(0, 0, 0), geometry.NewVec(165, 330, 165), white), 15), geometry.NewVec(265, 0, 295)),
	)

	lookAt := geometry.NewVec(278, 278, 0)
	lookFrom := geometry.NewVec(278, 278, -800)
	aperture := 0.1
	distToFocus := 10.0
	camera := newCamera(
		lookFrom,
		lookAt,
		geometry.NewVec(0, 1.0, 0),
		40,
		float64(width)/float64(height),
		aperture,
		distToFocus,
	)
	return camera, display.NewBVH(0, 0, 1, world.Hittables...)
}

// simpleLight is a scene with a Perlin-textured sphere and a rectangle light.
func simpleLight(width int, height int) (cameraSensor, *display.BVH) {
	world := display.List{}
	rnd := rand.New(rand.NewSource(rand.Int63()))
	perlin := display.NewNoise(rnd, 4)

	world.Hittables = append(world.Hittables,
		display.NewSphere(
			geometry.NewVec(0, -1000, 0), 1000, display.NewLambertian(perlin),
		),
		display.NewSphere(
			geometry.NewVec(0, 2, 0), 2, display.NewLambertian(perlin),
		),
		display.NewSphere(
			geometry.NewVec(0, 7, 0), 2, display.NewLight(display.NewColor(0, 2, 4)),
		),
		display.NewRectangle(
			geometry.NewVec(3, 1, -2), geometry.NewVec(5, 3, -2), display.NewLight(display.NewColor(4, 4, 4)),
		),
	)

	lookAt := geometry.NewVec(0, 2, 0)
	lookFrom := geometry.NewVec(24, 4, 6)
	aperture := 0.4
	distToFocus := 10.0
	camera := newCamera(
		lookFrom,
		lookAt,
		geometry.NewVec(0, 1.0, 0),
		40,
		float64(width)/float64(height),
		aperture,
		distToFocus,
	)
	return camera, display.NewBVH(0, 0, 1, world.Hittables...)
}

// jupiter is a simple sphere with a projection map of Jupiter.
func jupiter(width int, height int) (cameraSensor, *display.Sphere) {
	f, err := os.Open("assets/jupiter.jpeg")
	if err != nil {
		panic(err)
	}
	t, err := display.NewImage(f)
	if err != nil {
		panic(err)
	}
	lookAt := geometry.Vec{}
	lookFrom := geometry.NewVec(13, 2, 3)
	aperture := 0.1
	distToFocus := 10.0
	camera := newCamera(
		lookFrom,
		lookAt,
		geometry.NewVec(0, 1.0, 0),
		20,
		float64(width)/float64(height),
		aperture,
		distToFocus,
	)
	return camera, display.NewSphere(geometry.NewVec(0, 0, 0), 2, display.NewLambertian(t))
}

func buildTwoPerlinSpheresWorld(width, height int) (cameraSensor, *display.BVH) {
	world := display.List{}
	rnd := rand.New(rand.NewSource(rand.Int63()))
	perlin := display.NewNoise(rnd, 5)
	world.Hittables = append(world.Hittables,
		&display.Sphere{
			Center:   geometry.Vec{Y: -1000.0},
			Radius:   1000,
			Material: display.NewLambertian(perlin),
		},
		&display.Sphere{
			Center:   geometry.Vec{Y: 2.0},
			Radius:   2,
			Material: display.NewLambertian(perlin),
		},
	)
	lookAt := geometry.Vec{}
	lookFrom := geometry.NewVec(13, 2, 3)
	aperture := 0.1
	distToFocus := 10.0
	camera := newCamera(
		lookFrom,
		lookAt,
		geometry.NewVec(0, 1.0, 0),
		20,
		float64(width)/float64(height),
		aperture,
		distToFocus,
	)

	return camera, display.NewBVH(0, 0, 1, world.Hittables...)
}

// buildFinalWorld sets up the world and camera for the cover of the
// Ray Tracing in One Weekend book.
func buildFinalWorld(width, height int) (cameraSensor, *display.BVH) {
	world := display.List{}
	maxSpheres := 500

	checkered := display.NewChecker(10,
		display.NewSolid(display.NewColor(0.2, 0.3, 0.1)),
		display.NewSolid(display.NewColor(0.9, 0.9, 0.9)),
	)

	world.Hittables = append(world.Hittables,
		&display.Sphere{
			Center:   geometry.Vec{Y: -1000.0},
			Radius:   1000,
			Material: display.NewLambertian(checkered),
		},
		&display.Sphere{
			Center:   geometry.NewVec(0, 1, 0),
			Radius:   1.0,
			Material: display.NewDielectric(1.5),
		},
		&display.Sphere{
			Center:   geometry.NewVec(-4, 1, 0),
			Radius:   1.0,
			Material: display.NewLambertian(display.NewSolid(display.NewColor(0.4, 0.2, 0.1))),
		},
		&display.Sphere{
			Center:   geometry.NewVec(4, 1, 0),
			Radius:   1.0,
			Material: display.NewMetal(display.NewColor(0.7, 0.6, 0.5), 0),
		},
	)
	for a := -11; a < 11 && len(world.Hittables) < maxSpheres; a++ {
		for b := -11; b < 11 && len(world.Hittables) < maxSpheres; b++ {
			chooseMaterial := rand.Float64()
			center := geometry.NewVec(float64(a)+0.9*rand.Float64(), 0.2, float64(b)+0.9*rand.Float64())
			if center.Sub(geometry.NewVec(4.0, 0.2, 0)).Len() > 0.9 {
				switch {
				case chooseMaterial < 0.8:
					// Lambertian
					rnd := rand.New(rand.NewSource(rand.Int63()))
					center2 := center.Add(geometry.NewVec(0, geometry.FloatInRange(rnd, 0, 0.1), 0))
					world.Hittables = append(world.Hittables,
						display.NewMovingSphere(
							center,
							center2,
							0.0,
							1.0,
							0.2,
							display.NewLambertian(
								display.NewSolid(
									display.NewColor(
										rand.Float64()*rand.Float64(),
										rand.Float64()*rand.Float64(),
										rand.Float64()*rand.Float64(),
									),
								),
							),
						),
					)
				case chooseMaterial < 0.95:
					// Metal
					world.Hittables = append(world.Hittables,
						&display.Sphere{
							Center: center,
							Radius: 0.2,
							Material: display.NewMetal(
								display.NewColor(
									0.5*(1+rand.Float64()),
									0.5*(1+rand.Float64()),
									0.5*(1+rand.Float64()),
								),
								0.5*rand.Float64(),
							),
						},
					)
				default:
					// Dielectric
					world.Hittables = append(world.Hittables,
						&display.Sphere{
							Center:   center,
							Radius:   0.2,
							Material: display.NewDielectric(1.5),
						},
					)
				}
			}
		}
	}

	lookAt := geometry.Vec{}
	lookFrom := geometry.NewVec(13, 2, 3)
	aperture := 0.1
	distToFocus := 10.0
	camera := newCamera(
		lookFrom,
		lookAt,
		geometry.NewVec(0, 1.0, 0),
		20,
		float64(width)/float64(height),
		aperture,
		distToFocus,
	)

	return camera, display.NewBVH(0, 0, 1, world.Hittables...)
}
