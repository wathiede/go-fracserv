package solid

import (
	"fmt"
	"fractal"
	"image"
	"image/color"
	"strconv"
)

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

	w, err := strconv.Atoi(o.Get("w"))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse width %q: %s", o.Get("w"), err)
	}
	h, err := strconv.Atoi(o.Get("h"))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse height %q: %s", o.Get("h"), err)
	}

	return &Solid{*image.NewUniform(color.RGBA{r, g, b, 0xff}),
		image.Rect(0, 0, w, h)}, nil
}

func (s *Solid) Ratio() float32 {
	return 1
}

func (s *Solid) Bounds() image.Rectangle {
	return s.Rect
}
