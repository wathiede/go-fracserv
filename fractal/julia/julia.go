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
// Implementes Julia set, see:
//   http://eldar.mathstat.uoguelph.ca/dashlock/ftax/Julia.html
package julia

import (
	"code.google.com/p/go-fracserv/fractal"
	"fmt"
	"image"
	"image/color"
	"math/cmplx"
)

type Method int

const (
	method_unset = iota
	method_zSquared
	method_consine
)

type Julia struct {
	image.Paletted
	maxIterations int
	fractal.DefaultNavigator
	mu     complex128
	method Method
}

func NewFractal(opt fractal.Options) (fractal.Fractal, error) {
	it := opt.GetIntDefault("i", 256)
	mu_r := opt.GetFloat64Default("mu_r", 0.36237)
	mu_i := opt.GetFloat64Default("mu_i", 0.32)
	mu := complex(mu_r, mu_i)
	method := Method(opt.GetIntDefault("method", 1))
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
	return &Julia{*image.NewPaletted(image.Rect(0, 0, w, h), p), it, nav, mu, method}, nil
}

func (j *Julia) ColorIndexAt(x, y int) uint8 {
	r, i := j.Transform(image.Pt(x, y))

	return j.ComputeMembership(r, i)
}

func (j *Julia) ComputeMembership(r, i float64) uint8 {
	/*
		For every point (x,y) in your view rectangle 
		  Let z=x+yi
		  Set n=0
		  While(n less than limit and |z|<2)
		    Let z=z*z+mu
		    Increment n
		  End While
		  if(|z|<2) then z is a member of the approximate 
		    Julia set, plot (x,y) in the Julia set color
		  otherwise z is outside the Julia set, 
		    plot (x,y) in the outside color.
		End for

	*/
	z := complex(r, i)
	// Start at -1 so the first escaped values get the first color.
	it := -1
	switch j.method {
	case method_unset:
		panic("Julia method not set")
	case method_zSquared:
		for fractal.AbsLessThan(z, 2) && (it < j.maxIterations) {
			z = z*z + j.mu
			it++
		}
		if fractal.AbsLessThan(z, 2) {
			// Pixel in julia set, return black
			return 0
		}
	case method_consine:
		for fractal.AbsLessThan(z, 12) && (it < j.maxIterations) {
			z = cmplx.Cos(z) + j.mu
			it++
		}
		if fractal.AbsLessThan(z, 12) {
			// Pixel in julia set, return black
			return 0
		}
	default:
		panic(fmt.Sprintf("Unknown julia method %q", j.method))
	}

	// Black stored at j.Palette[0], so skip it
	return 1 + uint8(it%(len(j.Palette)-1))
}
