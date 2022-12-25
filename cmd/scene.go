package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/lucasmelin/raytracer/internal/display"
	"github.com/lucasmelin/raytracer/internal/geometry"
)

const (
	bias        = 0.001
	renderDepth = 10
)

// pixels represents the array of pixels to renderPixel.
type pixels []uint32

// scene represents the scene to render.
type scene struct {
	width, height int
	raysPerPixel  []int // array index represents the renderPixel pass
	camera        cameraSensor
	hitBoxer      display.HitBoxer
}

// pixel represents the pixel to be processed.
//
// x and y are the coordinates, k is the index in the pixels array, color is the color
// that has been computed by casting raysPerPixel through x/y coordinates.
type pixel struct {
	x            int
	y            int
	k            int
	color        display.Color
	raysPerPixel int
}

// split will split an array into an array of arrays with n number of elements.
func split(buf []*pixel, n int) [][]*pixel {
	var chunk []*pixel
	chunks := make([][]*pixel, 0, len(buf)/n+1)
	for len(buf) >= n {
		chunk, buf = buf[:n], buf[n:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf)
	}
	return chunks
}

// renderPixel casts rays one at a time through a pixel and accumulates the color for the pixel.
//
// Returns the normalized and gamma corrected value while updating the pixel for further ray casting.
func (scene *scene) renderPixel(rnd geometry.Rnd, pixel *pixel, raysPerPixel int) uint32 {
	c := pixel.color

	for s := 0; s < raysPerPixel; s++ {
		u := (float64(pixel.x) + rnd.Float64()) / float64(scene.width)
		v := (float64(pixel.y) + rnd.Float64()) / float64(scene.height)
		r := scene.camera.ray(rnd, u, v)
		c = c.Add(rayColor(r, scene.hitBoxer, 0))
	}

	pixel.color = c
	pixel.raysPerPixel += raysPerPixel

	// Normalize the color
	c = c.Scale(1.0 / float64(pixel.raysPerPixel))

	// Gamma correct
	c = display.NewColor(math.Sqrt(c.Red()), math.Sqrt(c.Green()), math.Sqrt(c.Blue()))

	return c.PixelValue()
}

// render returns the array of pixels to be computed asynchronously and a channel
// for signaling that the processing is complete.
// The image is split into lines, with each line being processed in a separate goroutine.
// The image is progressively rendered using the passes defined in raysPerPixel.
func (scene *scene) render(parallelCount int) (pixels, chan struct{}) {
	pixels := make([]uint32, scene.width*scene.height)
	completed := make(chan struct{})

	go func() {
		allPixelsToProcess := make([]*pixel, scene.width*scene.height)

		// Initializes the pixels, starting with black for no light.
		k := 0
		for j := scene.height - 1; j >= 0; j-- {
			for i := 0; i < scene.width; i++ {
				allPixelsToProcess[k] = &pixel{x: i, y: j, k: k}
				k++
			}
		}

		// Split the scene into lines
		lines := split(allPixelsToProcess, scene.width)

		// Compute the total numbers of rays to cast.
		totalRaysPerPixel := 0
		for _, rpp := range scene.raysPerPixel {
			totalRaysPerPixel += rpp
		}

		totalStart := time.Now()
		accumulatedRaysPerPixel := 0

		// Loop for each phase of the renderPixel.
		for _, rpp := range scene.raysPerPixel {

			loopStart := time.Now()

			// Create a channel for dispatching the line to process to each go routine.
			pixelsToProcess := make(chan []*pixel)

			// Dispatch the lines to process
			go func() {
				for _, p := range lines {
					pixelsToProcess <- p
				}
				// signal the end
				close(pixelsToProcess)
			}()

			// Wait until all goroutines have completed
			wg := sync.WaitGroup{}

			for c := 0; c < parallelCount; c++ {
				wg.Add(1)
				go func() {
					rnd := rand.New(rand.NewSource(rand.Int63()))

					// Process a line of pixels
					for ps := range pixelsToProcess {

						// Display the line without gamma correction so that it's more visible.
						for _, p := range ps {
							if p.raysPerPixel > 0 {
								col := p.color.Scale(1.0 / float64(p.raysPerPixel))
								pixels[p.k] = col.PixelValue()
							}
						}

						// render every pixel in the line one-by-one.
						for _, p := range ps {
							pixels[p.k] = scene.renderPixel(rnd, p, rpp)
						}
					}
					wg.Done()
				}()
			}

			// Wait for the entire renderPixel pass.
			wg.Wait()

			// Compute stats for the renderPixel pass.
			accumulatedRaysPerPixel += rpp

			loopEnd := time.Now()
			totalTime := loopEnd.Sub(totalStart)
			estimatedTotalTime := time.Duration(float64(totalTime) * float64(totalRaysPerPixel) / float64(accumulatedRaysPerPixel))
			erm := estimatedTotalTime - totalTime

			fmt.Printf("Processed %v rays per pixel in %v\nTotal %v in %v\nEst. Remaining Time: %s\n", rpp, time.Since(loopStart), accumulatedRaysPerPixel, totalTime, erm.Round(time.Second))
		}

		// signal completion
		completed <- struct{}{}
	}()

	return pixels, completed
}

// rayColor computes the color of the ray and scatters more rays according to the properties of the hittable.
func rayColor(r *geometry.Ray, hb display.HitBoxer, depth int) display.Color {
	// If we've exceeded the ray bounce limit, no more light is gathered.
	if depth >= renderDepth {
		return display.Black
	}
	if hit, hr := hb.Hit(r, bias, math.MaxFloat64); hit {
		if wasScattered, attenuation, scattered := hr.Material.Scatter(r, hr); wasScattered {
			indirect := attenuation.Mul(rayColor(scattered, hb, depth+1))
			return hr.Material.Emit(hr).Add(indirect)
		}
		return hr.Material.Emit(hr)
	}
	return display.Black
}
