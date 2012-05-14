// Implementes Glynn set, see:
//   http://eldar.mathstat.uoguelph.ca/dashlock/ftax/Mandel.html
package glynn

import (
	"fractal"
	"image"
	"image/color"
	"math/cmplx"
)

type Glynn struct {
	*image.Paletted
	maxIterations int
	fractal.DefaultNavigator
	e  float64
	mu complex128
}

func NewFractal(opt fractal.Options) (fractal.Fractal, error) {
	it := opt.GetIntDefault("i", 256)
	e := opt.GetFloat64Default("e", 1.5)
	mu := complex(-0.375, 0) //opt.GetIntDefault("mu", 1.5)
	w := opt.GetIntDefault("w", 256)
	h := opt.GetIntDefault("h", 256)
	x := opt.GetIntDefault("x", 0)
	y := opt.GetIntDefault("y", 0)
	z := opt.GetIntDefault("z", 0)

	p := color.Palette{color.Black}
	numColors := float64(12)
	for i := float64(0); i < numColors; i += 1 {
		degree := i / numColors * 360
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
	return &Glynn{image.NewPaletted(image.Rect(0, 0, w, h), p), it, nav, e, mu}, nil
}

func (g *Glynn) ColorIndexAt(x, y int) uint8 {
	/*
		Choose the parameter  mu=a+bi

		Pick a Glynn exponent, e.g. e=1.5.

		Pick an iteration cutoff limit.

		For every point (x,y) in your view rectangle 
		  Let z=x+yi
		  Set n=0
		  While(n less than limit and |z|<4)
		    Let z=z^e+mu
		    Increment n
		  End While
		  if(|z|<2) then z is a member of the approximate 
		    Glynn fractal, plot (x,y) in the Julia set color
		  otherwise z is outside the Glynn fractal, 
		    plot (x,y) in the outside color.
		End for
	*/
	r, i := g.Transform(image.Pt(x, y))
	z := complex(r, i)
	it := 0
	for (cmplx.Abs(z) < 4) && (it < g.maxIterations) {
		z = cmplx.Pow(z, g.e) + g.mu

		it++
		// Black stored at g.Palette[0], so skip it
		if it%len(g.Palette) == 0 {
			it++
		}
	}

	if cmplx.Abs(z) < 4 {
		// Pixel in glynn set, return black
		return 0
	}

	return uint8(it % len(g.Palette))
}
