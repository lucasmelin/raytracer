package main

import (
	"fmt"
	"io"
	"os"
)

const asciiColorPalette = "P3"
const maxColor = 255

func main() {
	render(os.Stdout, 256, 256)
}

func render(out io.Writer, width int, height int) {
	header := fmt.Sprintf("%s\n%d %d\n%d", asciiColorPalette, width, height, maxColor)
	fmt.Fprintln(out, header)

	for j := height - 1; j >= 0; j-- {
		for i := 0; i < width; i++ {
			red := float64(i) / float64(width-1)
			green := float64(j) / float64(height-1)
			blue := 0.25

			hueRed := toHue(red)
			hueGreen := toHue(green)
			hueBlue := toHue(blue)

			fmt.Fprintln(out, hueRed, hueGreen, hueBlue)
		}
	}
}

func toHue(value float64) int {
	return int(255.99 * value)
}
