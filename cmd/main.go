package main

import (
	"os"

	"github.com/lucasmelin/raytracer/internal/display"
)

func main() {
	display.Render(os.Stdout, 256, 256)
}
