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
package solid

import (
	"code.google.com/p/go-fracserv/fractal"
	"fmt"
	"image"
	"image/color"
)

func init() {
	fractal.Register("solid", NewFractal)
}

type Solid struct {
	image.Uniform
	Rect image.Rectangle
}

func NewFractal(o fractal.Options) (fractal.Fractal, error) {
	var r, g, b uint8
	c := o.Get("c")
	_, err := fmt.Sscanf(c, "%2x%2x%2x", &r, &g, &b)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse color %q: %s", c, err)
	}

	w := o.GetIntDefault("w", 256)
	h := o.GetIntDefault("h", 256)

	return &Solid{*image.NewUniform(color.RGBA{r, g, b, 0xff}),
		image.Rect(0, 0, w, h)}, nil
}

func (s *Solid) Bounds() image.Rectangle {
	return s.Rect
}
