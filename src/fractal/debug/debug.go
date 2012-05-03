package debug

import (
	"code.google.com/p/freetype-go/freetype"
	"fmt"
	"fractal"
	"io/ioutil"
	"image"
	"image/color"
	"math/rand"
	"sort"
)

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

	// Read the font data.
	fontBytes, err := ioutil.ReadFile("static/fonts/ProggyClean.ttf")
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
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
		pt.Y += ft.PointToFix32(fontSize * lineSpacing)
	}


	return &Debug{*im}, nil
}
