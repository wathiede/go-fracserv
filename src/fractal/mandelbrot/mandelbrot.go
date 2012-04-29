// Implementes Mandelbrot set, see:
//   http://eldar.mathstat.uoguelph.ca/dashlock/ftax/Mandel.html
package mandelbrot

import (
	"fractal"
	"image"
	"image/color"
	"math/cmplx"
)

type Mandelbrot struct {
	*image.Paletted
	maxIterations int
	fractal.DefaultNavigator
	order int
}

func NewFractal(opt fractal.Options) (fractal.Fractal, error) {
	it := opt.GetIntDefault("i", 256)
	o := opt.GetIntDefault("o", 2)
	w := opt.GetIntDefault("w", 256)
	h := opt.GetIntDefault("h", 256)
	x := opt.GetIntDefault("x", 0)
	y := opt.GetIntDefault("y", 0)
	z := opt.GetIntDefault("z", 0)

	p := color.Palette{color.Black}
	numColors := float64(12)
	for i := float64(0); i<numColors; i += 1 {
		degree := i/numColors * 360
		p = append(p, fractal.HSVToRGBA(degree, 1, 1))
	}

	// Center the image by considering the fractal range of (-2.5, 1), (-1, 1)
	// Explanation of magic numbers 1 and 7 being added to z
	//   +1: Google maps JS sends us zero based zoom so  we add one here to
	//       work with our math
	//   +7: The zoom factor is converted to 2^z, so having a zoom factor of 7
	//       (128x) makes the fractal range comfortably visible in pixel space
	nav := fractal.NewDefaultNavigator(float64(z+1+6), x*w, y*h)
	//nav := fractal.NewDefaultNavigator(float64(z+1)*200, x + int(-float64(w)/1.75), y - h/2)
	return &Mandelbrot{image.NewPaletted(image.Rect(0, 0, w, h), p), it, nav, o}, nil
}

//func (m *Mandelbrot) At(x, y int) color.Color {
func (m *Mandelbrot) ColorIndexAt(x, y int) uint8 {
	r, i := m.Transform(image.Pt(x, y))
	z := complex(r, i)
	w := complex(0, 0)
	// Start at -1 so the first escaped values get the first color.
	it := -1
	for (cmplx.Abs(w) < 2) && (it < m.maxIterations) {
		v := w
		for i := 1; i<m.order; i++ {
			v *= w
		}
		w = v + z

		it++
	}

	if cmplx.Abs(w) < 2 {
		// Pixel in mandelbrot set, return black
		return 0
	}

	// Black stored at m.Palette[0], so skip it
	return 1 + uint8(it % (len(m.Palette)-1))
}
