package solid

import (
	"image"
	"image/color"
	"testing"
)

func TestColor(t *testing.T) {
	o := Options{
		image.Config{color.RGBAModel, 10, 10},
		color.RGBA{10, 20, 30, 40}}
	s := NewSolid(o)

	c := s.At(5, 5)
	t.Log("c:", c)
	r, g, b, a := c.RGBA()
	t.Log(r, g, b, a)
	switch {
	case r >> 8 != 10:
		t.Errorf("Red not right, expected 10, got %d\n", r)
	case g >> 8 != 20:
		t.Errorf("Green not right, expected 20, got %d\n", g)
	case b >> 8 != 30:
		t.Errorf("Blue not right, expected 30, got %d\n", b)
	case a >> 8 != 40:
		t.Errorf("Alpha not right, expected 40, got %d\n", a)
	}
}

func TestDimension(t *testing.T) {
	o := Options{
		image.Config{color.RGBAModel, 10, 10},
		color.RGBA{1, 2, 3, 4}}
	s := NewSolid(o)

	rect := image.Rect(0, 0, 10, 10)
	bounds := s.Bounds()
	if !bounds.Eq(rect) {
		t.Errorf("Expected %v got %v", rect, bounds)
	}
}
