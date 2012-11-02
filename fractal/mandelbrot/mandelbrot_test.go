package mandelbrot

import (
	"image/png"
	"io/ioutil"
	"math"
	"strconv"
	"testing"

	"code.google.com/p/go-fracserv/fractal"
)

func BenchmarkTiles(b *testing.B) {
	// Randomly chosen tile of moderate complexity
	// /mandelbrot?w=128&h=128&x=-44&y=2&z=5&o=2&i=50&name=&url=
	o := fractal.NewOptions()
	o.Set("i", "50")
	o.Set("o", "2")
	o.Set("w", "128")
	o.Set("h", "128")
	o.Set("z", "5")

	size := int(math.Sqrt(float64(b.N)))
	xTiles := size
	yTiles := size

	for y := -yTiles / 2; y < yTiles/2; y++ {
		for x := -xTiles / 2; x < xTiles/2; x++ {
			o.Set("x", strconv.Itoa(x))
			o.Set("y", strconv.Itoa(x))

			f, err := NewFractal(o)
			if err != nil {
				b.Fatal(err)
			}
			png.Encode(ioutil.Discard, f)
		}
	}
}
