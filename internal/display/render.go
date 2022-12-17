package display

import (
	"fmt"
	"io"
	"os"
)

const asciiColorPalette = "P3"
const maxColor = 255

func Render(out io.Writer, width int, height int) {
	header := fmt.Sprintf("%s\n%d %d\n%d", asciiColorPalette, width, height, maxColor)
	fmt.Fprintln(out, header)

	for j := height - 1; j >= 0; j-- {
		fmt.Fprintf(os.Stderr, "\nScanlines remaining: %d", j)
		for i := 0; i < width; i++ {
			red := float64(i) / float64(width-1)
			green := float64(j) / float64(height-1)
			blue := 0.25
			c := NewColor(red, green, blue)
			WriteColor(out, c)
		}
	}
	fmt.Fprintf(os.Stderr, "\nDone\n")
}

func toHue(value float64) int {
	return int(255.99 * value)
}
