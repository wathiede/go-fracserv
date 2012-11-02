// Copyright 2012 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Implementes Mandelbrot set, see:
//   http://eldar.mathstat.uoguelph.ca/dashlock/ftax/Mandel.html
package mandelbrot

import (
	"image"
	"image/color"

	"code.google.com/p/go-fracserv/fractal"
)

func init() {
	fractal.Register("mandelbrot", NewFractal)
}

type Mandelbrot struct {
	image.Paletted
	fractal.DefaultNavigator
	maxIterations int
	order         int
}

func NewFractal(opt fractal.Options) (fractal.Fractal, error) {
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
	nav := fractal.NewDefaultNavigator(uint(z+1+6), x*w, y*h)
	//nav := fractal.NewDefaultNavigator(float64(z+1)*200, x + int(-float64(w)/1.75), y - h/2)
	return &Mandelbrot{*image.NewPaletted(image.Rect(0, 0, w, h), p), nav,
		opt.GetIntDefault("i", 256), opt.GetIntDefault("o", 2)}, nil
}

func (m *Mandelbrot) ColorIndexAt(x, y int) uint8 {
	r, i := m.Transform(image.Pt(x, y))

	return m.ComputeMembership(r, i)
}

// Takes in a coordinate in fractal space, and returns an index to the proper
// coloring for that point
func (m *Mandelbrot) ComputeMembership(r, i float64) uint8 {
	z := complex(r, i)
	w := complex(0, 0)
	// Start at -1 so the first escaped values get the first color.
	it := -1
	for fractal.AbsLessThan(w, 2) && (it < m.maxIterations) {
		v := w
		for i := 1; i < m.order; i++ {
			v *= w
		}
		w = v + z

		it++
	}

	if fractal.AbsLessThan(w, 2) {
		// Pixel in mandelbrot set, return black
		return 0
	}

	// Black stored at m.Palette[0], so skip it
	return 1 + uint8(it%(len(m.Palette)-1))
}
