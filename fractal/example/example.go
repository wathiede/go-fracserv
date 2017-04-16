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

// TODO REPLACEME:
// This package implements an example fractal.  It takes a single parameter
// "wrap", and uses it to modify the output.  Your fractal will probably want
// to do something more interesting in ComputeMembership
package example

import (
	"image"
	"image/color"
	"math"

	"github.com/wathiede/go-fracserv/fractal"
)

func init() {
	fractal.Register("example", NewFractal)
}

type Example struct {
	image.Paletted
	fractal.DefaultNavigator
	wrapFactor int
}

func NewFractal(opt fractal.Options) (fractal.Fractal, error) {
	w := opt.GetIntDefault("w", 256)
	h := opt.GetIntDefault("h", 256)
	x := opt.GetIntDefault("x", 0)
	y := opt.GetIntDefault("y", 0)
	z := opt.GetIntDefault("z", 0)

	// Setup color palette
	p := color.Palette{color.Black}
	numColors := float64(12)
	for i := float64(0); i < numColors; i += 1 {
		degree := i / numColors * 360
		p = append(p, fractal.HSVToRGBA(degree, 1, 1))
	}

	// Navigator provides a way to scale between pixel space and the number
	// space for your fractal type.  It also takes the zoom factor into
	// account so the Google Maps API will work properly.  If you need to do
	// something more complicated in your navigation of the fractal space,
	// create your own Navigator instead of using the DefaultNavigator
	nav := fractal.NewDefaultNavigator(uint(z), x*w, y*h)

	return &Example{*image.NewPaletted(image.Rect(0, 0, w, h), p), nav,
		opt.GetIntDefault("wrap", 4)}, nil
}

func (f *Example) ColorIndexAt(x, y int) uint8 {
	r, i := f.Transform(image.Pt(x, y))

	return f.ComputeMembership(r, i)
}

// Takes in a coordinate in fractal space, and returns an index to the proper
// coloring for that point
func (f *Example) ComputeMembership(r, i float64) uint8 {
	// See mandelbrot for an example of a fractal that uses and escape
	// threshold to limit the amount of computation done per pixel.  For our
	// example case we just do something simple.  Normally the fractal
	// algorithm would be implemented here.

	// TODO Implement your fractal here
	idx := int(math.Floor(r * i / float64(f.wrapFactor)))
	return uint8(int(idx) % len(f.Palette))
}
