package mandelbrot

import (
	"fractal"
	"image"
	"image/color"
	"math"
	"math/cmplx"
)

type Mandelbrot struct {
	*image.Paletted
	maxIterations int
	fractal.DefaultNavigator
}

func HSVToRGBA(h, s, v float64) color.RGBA {
	hi := int(math.Mod(math.Floor(h / 60), 6))
	f := (h / 60) - math.Floor(h / 60)
	p := v * (1 - s)
	q := v * (1 - (f*s))
	t := v * (1 - ((1 - f) * s))

	RGB := func(r, g, b float64) color.RGBA {
		return color.RGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), 255}
	}

	var c color.RGBA
	switch hi {
		case 0: c = RGB(v, t, p)
		case 1: c = RGB(q, v, p)
		case 2: c = RGB(p, v, t)
		case 3: c = RGB(p, q, v)
		case 4: c = RGB(t, p, v)
		case 5: c = RGB(v, p, q)
	}
	return c
}

func NewFractal(o fractal.Options) (fractal.Fractal, error) {
	w := o.GetIntDefault("w", 256)
	h := o.GetIntDefault("h", 256)
	x := o.GetIntDefault("x", 0)
	y := o.GetIntDefault("y", 0)
	it := o.GetIntDefault("i", 256)
	z := o.GetIntDefault("z", 0)

	p := color.Palette{color.Black}
	numColors := float64(12)
	for i := float64(0); i<numColors; i += 1 {
		degree := i/numColors * 360
		p = append(p, HSVToRGBA(degree, 1, 1))
	}

	// Center the image by considering the fractal range of (-2.5, 1), (-1, 1)
	// Explanation of magic numbers 1 and 7 being added to z
	//   +1: Google maps JS sends us zero based zoom so  we add one here to
	//       work with our math
	//   +7: The zoom factor is converted to 2^z, so having a zoom factor of 7
	//       (128x) makes the fractal range comfortably visible in pixel space
	nav := fractal.NewDefaultNavigator(float64(z+1+6), x, y)
	//nav := fractal.NewDefaultNavigator(float64(z+1)*200, x + int(-float64(w)/1.75), y - h/2)
	return &Mandelbrot{image.NewPaletted(image.Rect(0, 0, w, h), p), it, nav}, nil
}

//func (m *Mandelbrot) At(x, y int) color.Color {
func (m *Mandelbrot) ColorIndexAt(x, y int) uint8 {
	/*
	For every point (x,y) in your view rectangle 
	  Let z=x+yi
	  Set n=0
	  Set w=0
	  While(n less than limit and |w|<2)
		Let w=w*w+z
		Increment n
	  End While
	  if(|w|<2) then z is a member of the approximate 
		Mandelbrot set, plot (x,y) in the Mandelbrot set color
	  otherwise z is outside the Mandelbrot set, 
		plot (x,y) in the outside color.
	End for
	*/
	r, i := m.Transform(image.Pt(x, y))
	z := complex(r, i)
	w := complex(0, 0)
	it := 0
	for (cmplx.Abs(w) < 2) && (it < m.maxIterations) {
		w = w * w + z

		it++
		// Black stored at m.Palette[0], so skip it
		if it % len(m.Palette) == 0 {
			it++
		}
	}

	if cmplx.Abs(w) < 2 {
		// Pixel in mandelbrot set, return black
		return 0
	}

	return uint8(it % len(m.Palette))
}

func (m *Mandelbrot) Transform(p image.Point) (float64, float64) {
	b := m.Bounds()
	t := m.GetTranslate()
	x := p.X + (t.X * b.Dx())
	y := p.Y + (t.Y * b.Dy())
	z := math.Pow(2, m.GetZoom())
	return float64(x)/z, float64(y)/z
}
