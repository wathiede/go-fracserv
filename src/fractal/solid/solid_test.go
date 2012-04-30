package solid

import (
	"fractal"
	"image"
	"net/url"
	"testing"
)

func TestColor(t *testing.T) {
	o := fractal.Options{
		url.Values{
			"w": []string{"10"},
			"h": []string{"10"},
			"c": []string{"224488"},
			//color.RGBA{10, 20, 30, 40}
		},
	}
	s, err := NewFractal(o)
	if err != nil {
		t.Errorf("Failed to create NewFractal: %s", err)
	}

	c := s.At(5, 5)
	t.Log("c:", c)
	r, g, b, a := c.RGBA()
	t.Log(r, g, b, a)
	switch {
	case r>>8 != 0x22:
		t.Errorf("Red not right, expected 10, got %d\n", r)
	case g>>8 != 0x44:
		t.Errorf("Green not right, expected 20, got %d\n", g)
	case b>>8 != 0x88:
		t.Errorf("Blue not right, expected 30, got %d\n", b)
	case a>>8 != 0xff:
		t.Errorf("Alpha not right, expected 40, got %d\n", a)
	}
}

func TestDimension(t *testing.T) {
	o := fractal.Options{
		url.Values{
			"w": []string{"10"},
			"h": []string{"10"},
			"c": []string{"224488"},
		},
	}
	s, err := NewFractal(o)
	if err != nil {
		t.Errorf("Failed to create NewFractal: %s", err)
	}

	rect := image.Rect(0, 0, 10, 10)
	bounds := s.Bounds()
	if !bounds.Eq(rect) {
		t.Errorf("Expected %v got %v", rect, bounds)
	}
}
