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

	"github.com/veandco/go-sdl2/sdl"
)

// raysPerPixelList is used to define the number of rays per-pixel, per phase.
type raysPerPixelList []int

// String allows for printing the raysPerPixelList.
func (r *raysPerPixelList) String() string {
	return fmt.Sprint(*r)
}

// Set parses the number of rays to use per pixel, per phase.
func (r *raysPerPixelList) Set(value string) error {
	for _, e := range strings.Split(value, ",") {
		i, err := strconv.Atoi(e)
		if err != nil {
			return fmt.Errorf("could not parse %q is int: %w", e, err)
		}
		*r = append(*r, i)
	}
	return nil
}

// options defines the command line options.
type options struct {
	Width        int
	Height       int
	RaysPerPixel raysPerPixelList
	Output       string
	Seed         int64
	CPU          int
}

// disp will update the display with the pixels as they get rendered by each goroutine.
func disp(window *sdl.Window, screen *sdl.Surface, scene *scene, pixels pixels) {
	// Create an img from the generated pixels.
	img, err := sdl.CreateRGBSurfaceFrom(
		// https://pkg.go.dev/unsafe#Pointer
		unsafe.Pointer(&pixels[0]),
		int32(scene.width),
		int32(scene.height),
		32,
		scene.width*int(unsafe.Sizeof(pixels[0])), 0, 0, 0, 0)
	if err != nil {
		panic(err)
	}
	defer img.Free()

	// Copy to img to the screen.
	err = img.Blit(nil, screen, nil)
	if err != nil {
		panic(err)
	}

	// Update the surface to display.
	if err = window.UpdateSurface(); err != nil {
		panic(err)
	}
}

// saveImage saves the image to a file in png format.
func saveImage(pixels pixels, options options) (bool, error) {
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

func main() {
	options := options{}

	flag.IntVar(&options.Width, "w", 800, "width in pixels")
	flag.IntVar(&options.Height, "h", 400, "height in pixels")
	flag.IntVar(&options.CPU, "cpu", runtime.NumCPU(), "number of CPU to use (default number of available CPUs)")
	flag.Int64Var(&options.Seed, "seed", 1992, "seed for random number generator")
	flag.Var(&options.RaysPerPixel, "r", "comma separated list of rays-per-pixel")
	flag.StringVar(&options.Output, "o", "image.png", "path to output file")

	flag.Parse()

	if len(options.RaysPerPixel) == 0 {
		// Default 10 rays on the first pass, 90 rays on the subsequent pass.
		options.RaysPerPixel = []int{10, 90}
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

	camera, bvh := buildWeekOneWorld(options.Width, options.Height)

	scene := &scene{
		width:        options.Width,
		height:       options.Height,
		raysPerPixel: options.RaysPerPixel,
		camera:       camera,
		hitBoxer:     bvh,
	}
	pixels, completed := scene.render(options.CPU)

	// Show the initial renderPixel pass.
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
				fmt.Println("render cancelled")
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
				fmt.Println("render complete")
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
