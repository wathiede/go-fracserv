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
package debug

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"sort"

	"github.com/golang/freetype"
	"github.com/wathiede/go-fracserv/fractal"
	"golang.org/x/image/font/gofont/goregular"
)

func init() {
	fractal.Register("debug", NewFractal)
}

type Debug struct {
	image.RGBA
}

func NewFractal(o fractal.Options) (fractal.Fractal, error) {
	w := o.GetIntDefault("w", 256)
	h := o.GetIntDefault("h", 256)
	im := image.NewRGBA(image.Rect(0, 0, w, h))

	c := color.RGBA{
		uint8(rand.Int()),
		uint8(rand.Int()),
		uint8(rand.Int()),
		255,
	}
	border := 4
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			switch {
			case x < border || x > w-border:
				im.SetRGBA(x, y, c)
			case y < border || y > h-border:
				im.SetRGBA(x, y, c)
			default:
				im.Set(x, y, color.Black)
			}
		}
	}

	font, err := freetype.ParseFont(goregular.TTF)
	if err != nil {
		return nil, err
	}

	fontSize := float64(16)
	lineSpacing := 1.
	ft := freetype.NewContext()
	ft.SetDPI(72)
	ft.SetFont(font)
	ft.SetFontSize(fontSize)
	ft.SetClip(im.Bounds())
	ft.SetDst(im)
	ft.SetSrc(image.White)

	pt := freetype.Pt(6, 18)
	keys := []string{}
	for k := range o.Values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		_, err = ft.DrawString(fmt.Sprintf("%s = %s", k, o.Get(k)), pt)
		if err != nil {
			return nil, err
		}
		pt.Y += ft.PointToFixed(fontSize * lineSpacing)
	}

	return &Debug{*im}, nil
}
