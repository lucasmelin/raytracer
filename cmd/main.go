// Simple ray tracer based on the Ray Tracing book series by Peter Shirley (Kindle)
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"unsafe"

	"github.com/lucasmelin/raytracer/internal/display"
	"github.com/lucasmelin/raytracer/internal/geometry"
	"github.com/veandco/go-sdl2/sdl"
)

// RaysPerPixelList is used to define the number of rays per-pixel, per phase.
type RaysPerPixelList []int

// String allows for printing the RaysPerPixelList.
func (r *RaysPerPixelList) String() string {
	return fmt.Sprint(*r)
}

// Set parses the number of rays to use per pixel, per phase.
func (r *RaysPerPixelList) Set(value string) error {
	for _, e := range strings.Split(value, ",") {
		i, err := strconv.Atoi(e)
		if err != nil {
			return fmt.Errorf("could not parse %q is int: %w", e, err)
		}
		*r = append(*r, i)
	}
	return nil
}

// Options defines the command line options.
type Options struct {
	Width        int
	Height       int
	RaysPerPixel RaysPerPixelList
	Output       string
	Seed         int64
	CPU          int
}

// disp will update the display with the pixels as they get rendered by each goroutine.
func disp(window *sdl.Window, screen *sdl.Surface, scene *Scene, pixels Pixels) {
	// Create an image from the generated pixels.
	image, err := sdl.CreateRGBSurfaceFrom(
		// https://pkg.go.dev/unsafe#Pointer
		unsafe.Pointer(&pixels[0]),
		int32(scene.width),
		int32(scene.height),
		32,
		scene.width*int(unsafe.Sizeof(pixels[0])), 0, 0, 0, 0)
	if err != nil {
		panic(err)
	}
	defer image.Free()

	// Copy to image tothe screen.
	err = image.Blit(nil, screen, nil)
	if err != nil {
		panic(err)
	}

	// Update the surface to display.
	if err = window.UpdateSurface(); err != nil {
		panic(err)
	}
}

func cornell(width int, height int) (Camera, *display.BVH) {
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
			geometry.NewVec(213, 554, 227), geometry.NewVec(343, 554, 332), display.NewLight(display.NewColor(4, 4, 4)),
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
	)

	lookAt := geometry.NewVec(278, 278, 0)
	lookFrom := geometry.NewVec(278, 278, -800)
	aperture := 0.1
	distToFocus := 10.0
	camera := NewCamera(
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

func simpleLight(width int, height int) (Camera, *display.BVH) {
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
	camera := NewCamera(
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

func jupiter(width int, height int) (Camera, *display.Sphere) {
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
	camera := NewCamera(
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

func buildTwoPerlinSpheresWorld(width, height int) (Camera, *display.BVH) {
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
	camera := NewCamera(
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

// buildFinalWorld sets up the world and camera for the final scene/cover of the book.
func buildFinalWorld(width, height int) (Camera, *display.BVH) {
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
	camera := NewCamera(
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

// saveImage saves the image to a file in png format.
func saveImage(pixels Pixels, options Options) (bool, error) {
	if options.Output != "" {
		f, err := os.OpenFile(options.Output, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return true, err
		}

		img := image.NewNRGBA(image.Rect(0, 0, options.Width, options.Height))

		k := 0
		for y := 0; y < options.Height; y++ {
			for x := 0; x < options.Width; x++ {
				p := pixels[k]
				img.Set(x, y, color.NRGBA{
					R: uint8(p >> 16 & 0xFF),
					G: uint8(p >> 8 & 0xFF),
					B: uint8(p & 0xFF),
					A: 255,
				})
				k++
			}
		}
		if err := png.Encode(f, img); err != nil {
			f.Close()
			return true, err
		}
		if err := f.Close(); err != nil {
			return true, err
		}
		return true, nil
	}

	return false, nil

}

// main parses the options, set up the Window/Screen, builds the world and renders the scene.
// As the scene gets rendered the screen gets refreshed regularly to show progress. When the image is fully
// rendered, it saves it to a file (if the output option is set)
func main() {
	options := Options{}

	flag.IntVar(&options.Width, "w", 800, "width in pixels")
	flag.IntVar(&options.Height, "h", 400, "height in pixels")
	flag.IntVar(&options.CPU, "cpu", runtime.NumCPU(), "number of CPU to use (default number of available CPUs)")
	flag.Int64Var(&options.Seed, "seed", 1337, "seed for random number generator")
	flag.Var(&options.RaysPerPixel, "r", "comma separated list of rays-per-pixel")
	flag.StringVar(&options.Output, "o", "image.png", "path to output file")

	flag.Parse()

	if len(options.RaysPerPixel) == 0 {
		// Default 1 ray on the first pass, 199 rays on the subsequent pass.
		options.RaysPerPixel = []int{1, 399}
	}

	rand.Seed(options.Seed)

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		newErr := fmt.Errorf("could not initialize SDL: %w", err)
		panic(newErr)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"Raytracer",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(options.Width),
		int32(options.Height),
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		newErr := fmt.Errorf("could not create window using SDL: %w", err)
		panic(newErr)
	}
	defer window.Destroy()

	screen, err := window.GetSurface()
	if err != nil {
		newErr := fmt.Errorf("could not retrieve window using SDL: %w", err)
		panic(newErr)
	}

	// Fill the screen so that it is blank.
	if err = screen.FillRect(&sdl.Rect{W: int32(options.Width), H: int32(options.Height)}, 0x00000000); err != nil {
		newErr := fmt.Errorf("could not blank out screen: %w", err)
		panic(newErr)
	}

	camera, bvh := cornell(options.Width, options.Height)

	scene := &Scene{
		width:        options.Width,
		height:       options.Height,
		raysPerPixel: options.RaysPerPixel,
		camera:       camera,
		hitBoxer:     bvh,
	}
	pixels, completed := scene.Render(options.CPU)

	// Show the initial render pass.
	if err = window.UpdateSurface(); err != nil {
		newErr := fmt.Errorf("could not display screen: %w", err)
		panic(newErr)
	}

	updateDisplay := true
	for {
		// Poll for quit event from SDL in case the window is closed.
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Render cancelled")
				os.Exit(1)
			}
		}

		// Wait for a few ms between iterations.
		sdl.Delay(15)

		if updateDisplay {
			disp(window, screen, scene, pixels)

			// Check if the image is completely rendered.
			select {
			case <-completed:
				updateDisplay = false
				fmt.Println("Render complete")
				saved, err := saveImage(pixels, options)
				if err != nil {
					fmt.Printf("Error saving image: %q", err.Error())
				}
				if saved {
					fmt.Printf("Image saved to %s\n", options.Output)
				}
			default:
				break
			}

		}
	}
}
