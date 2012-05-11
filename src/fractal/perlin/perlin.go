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

// Implementes Perlin noise, see:
//   http://freespace.virgin.net/hugo.elias/models/m_perlin.htm
package perlin

import (
	"fractal"
	"log"
	"image"
	"image/color"
	"math"
)

type Perlin struct {
	image.Gray
	fractal.DefaultNavigator
	octaves uint
	persistence float64
}

func NewFractal(opt fractal.Options) (fractal.Fractal, error) {
	w := opt.GetIntDefault("w", 256)
	h := opt.GetIntDefault("h", 256)
	x := opt.GetIntDefault("x", 0)
	y := opt.GetIntDefault("y", 0)
	z := opt.GetIntDefault("z", 0)
	octaves := uint(opt.GetIntDefault("o", 2))
	persistence := opt.GetFloat64Default("p", 1/2)

	p := color.Palette{color.Black}
	numColors := float64(12)
	for i := float64(0); i < numColors; i += 1 {
		degree := i / numColors * 360
		p = append(p, fractal.HSVToRGBA(degree, 1, 1))
	}

	nav := fractal.NewDefaultNavigator(uint(z), x*w, y*h)
	return &Perlin{*image.NewGray(image.Rect(0, 0, w, h)), nav, octaves, persistence}, nil
}

/*
  function Noise1(integer x, integer y)
    n = x + y * 57
    n = (n<<13) ^ n;
    return ( 1.0 - ( (n * (n * n * 15731 + 789221) + 1376312589) & 7fffffff) / 1073741824.0);    
  end function

  function SmoothNoise_1(float x, float y)
    corners = ( Noise(x-1, y-1)+Noise(x+1, y-1)+Noise(x-1, y+1)+Noise(x+1, y+1) ) / 16
    sides   = ( Noise(x-1, y)  +Noise(x+1, y)  +Noise(x, y-1)  +Noise(x, y+1) ) /  8
    center  =  Noise(x, y) / 4
    return corners + sides + center
  end function

  function InterpolatedNoise_1(float x, float y)

      integer_X    = int(x)
      fractional_X = x - integer_X

      integer_Y    = int(y)
      fractional_Y = y - integer_Y

      v1 = SmoothedNoise1(integer_X,     integer_Y)
      v2 = SmoothedNoise1(integer_X + 1, integer_Y)
      v3 = SmoothedNoise1(integer_X,     integer_Y + 1)
      v4 = SmoothedNoise1(integer_X + 1, integer_Y + 1)

      i1 = Interpolate(v1 , v2 , fractional_X)
      i2 = Interpolate(v3 , v4 , fractional_X)

      return Interpolate(i1 , i2 , fractional_Y)

  end function


  function PerlinNoise_2D(float x, float y)

      total = 0
      p = persistence
      n = Number_Of_Octaves - 1

      loop i from 0 to n

          frequency = 2i
          amplitude = pi

          total = total + InterpolatedNoisei(x * frequency, y * frequency) * amplitude

      end of i loop

      return total

  end function
*/

func (p *Perlin) At(x, y int) color.Color {
	r, i := p.Transform(image.Pt(x, y))

	return p.perlinNoise2D(r, i)
}

func linearInterpolate(a, b, x float64) float64 {
	switch {
	case x < 0, x > 1:
		log.Printf("x out of range for %f, %f: %f", a, b, x)
	}
	return  a*(1-x) + b*x
}

func cosineInterpolate(a, b, x float64) float64 {
	ft := x * math.Pi
	f := (1 - math.Cos(ft)) * .5

	return  a*(1-f) + b*f
}

func noise(x, y float64) float64 {
	_, f := math.Modf(x / 256)
	return math.Abs(f)
	n := int(x + y * 57)
	n = (n<<13) ^ n
	return  1 - float64((n * (n * n * 15731 + 789221) + 1376312589) & 0x7fffffff) / 1073741824.0
}

func smoothedNoise(x, y float64) float64 {
	corners := ( noise(x-1, y-1)+noise(x+1, y-1)+noise(x-1, y+1)+noise(x+1, y+1) ) / 16
	sides   := ( noise(x-1, y)  +noise(x+1, y)  +noise(x, y-1)  +noise(x, y+1) ) /  8
	center  :=  noise(x, y) / 4
	return corners + sides + center
}

func interpolatedNoise(x, y float64) float64 {
	interpolate := linearInterpolate //cosineInterpolate

	integer_X, fractional_X := math.Modf(x)
	integer_Y, fractional_Y := math.Modf(y)

	v1 := smoothedNoise(integer_X,     integer_Y)
	v2 := smoothedNoise(integer_X + 1, integer_Y)
	v3 := smoothedNoise(integer_X,     integer_Y + 1)
	v4 := smoothedNoise(integer_X + 1, integer_Y + 1)

	i1 := interpolate(v1 , v2 , fractional_X)
	i2 := interpolate(v3 , v4 , fractional_X)

	return interpolate(i1 , i2 , fractional_Y)
}

func (p *Perlin) perlinNoise2D(x, y float64) color.Gray {
	total := float64(0)

	for i := uint(0); i < p.octaves; i++ {
		frequency := float64(uint(1) << i)
		amplitude := math.Pow(p.persistence, float64(i))

		total = total + interpolatedNoise(x * frequency, y * frequency) * amplitude
	}

	/*
	if total > 1 {
		log.Printf("total greater than 1 for %f, %f: %f", x, y, total)
	}
	if total < 0 {
		log.Printf("total less than 0 for %f, %f: %f", x, y, total)
	}
	*/
	return color.Gray{uint8(total * 255)}
}
