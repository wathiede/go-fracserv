package mandelbrot

import (
	"fmt"
	"fractal"
	"image"
	"image/color"
	"math/rand"
	"strconv"
)

type Mandelbrot struct {
	*image.Paletted
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

	p := color.Palette{}
	for i := 0; i<256; i++ {
		p = append(p, color.RGBA{
			uint8(rand.Int()) % 255,
			uint8(rand.Int()) % 255,
			uint8(rand.Int()) % 255, 255})
	}

	r := 3.5/2.0
	if float64(w) / r > float64(h) {
		w = int(float64(h) * r)
	} else {
		h = int(float64(w) / r)
	}
	return &Mandelbrot{image.NewPaletted(image.Rect(0, 0, w, h), p)}, nil
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

	var iteration uint8 = 0
	var max_iteration uint8 = 255

	for (tx*tx + ty*ty < 4)  && (iteration < max_iteration) {
		xtemp := tx*tx - ty*ty + x0
		ty = 2*tx*ty + y0

		tx = xtemp

		iteration++
	}

	return iteration
}

func (m *Mandelbrot) Ratio() float32 {
	return 3.5/2.0
}
