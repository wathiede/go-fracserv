package mandelbrot

import (
	"fmt"
	"fractal"
	"image"
	"image/color"
	"strconv"
)

type Mandelbrot struct {
	*image.Paletted
	maxIterations int
}

func NewFractal(o fractal.Options) (fractal.Fractal, error) {
	w, err := strconv.Atoi(o.Get("w"))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse width %q: %s", o.Get("w"), err)
	}

	h, err := strconv.Atoi(o.Get("h"))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse height %q: %s", o.Get("h"), err)
	}

	i, err := strconv.Atoi(o.Get("i"))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse max iterations %q: %s", o.Get("i"), err)
	}

	p := color.Palette{}
	for idx := uint8(0); i<256; i++ {
		p = append(p, color.RGBA{
			uint8(idx & 0x3),
			uint8(idx >> 2 & 0x7),
			uint8(idx >> 5 & 0x3),
			0xff,
		})
	}

	r := 3.5/2.0
	if float64(w) / r > float64(h) {
		w = int(float64(h) * r)
	} else {
		h = int(float64(w) / r)
	}
	return &Mandelbrot{image.NewPaletted(image.Rect(0, 0, w, h), p), i}, nil
}

//func (m *Mandelbrot) At(x, y int) color.Color {
func (m *Mandelbrot) ColorIndexAt(x, y int) uint8 {
	b := m.Bounds()
	/*
	Normalize pixel (x,y) in range:
	x0 := (-2.5 to 1)
	y0 := (-1, 1)
	*/
	w := b.Dx()
	h := b.Dy()

	x0 := (float64(x) / float64(w) * 3.5) - 2.5
	y0 := (float64(y) / float64(h) * 2) - 1

	var tx, ty float64

	iteration := 0

	for (tx*tx + ty*ty < 4)  && (iteration < m.maxIterations) {
		xtemp := tx*tx - ty*ty + x0
		ty = 2*tx*ty + y0

		tx = xtemp

		iteration++
	}

	return uint8(iteration % 256)
}

func (m *Mandelbrot) Ratio() float32 {
	return 3.5/2.0
}
