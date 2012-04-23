package solid

import (
	"image"
	"image/color"
)

type Solid struct {
	image.Image
	Rect image.Rectangle
}

type Options struct {
	image.Config
	C color.Color
}

func NewSolid(o Options) image.Image {
	return &Solid{image.NewUniform(o.C),
		image.Rect(0, 0, o.Width, o.Height)}
}

func (s *Solid) Bounds() image.Rectangle {
	return s.Rect
}
