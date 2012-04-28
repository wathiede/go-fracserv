package fractal

import (
	"log"
	"image"
	"net/url"
	"strconv"
)

type Options struct {
	url.Values
}

// Converts value with key 'k' to int, in absence or failure 'd' is returned
func (o Options) GetIntDefault(k string, d int) int {
	v, err := strconv.Atoi(o.Get(k))
	if err != nil {
		log.Printf("Failed to parse %s: %s", k, o.Get(k), err)
		return d
	}
	return v
}

type Fractal interface {
	image.Image
}

type Navigator interface {
	// Convert pixel space to fractal space
	Transform(p image.Point) (float64, float64)
	// Set offset in pixel space
	Translate(offset image.Point)
	// Set the zoom depth for this transformation
	Zoom(z float64)
}

type DefaultNavigator struct {
	z float64
	offset image.Point
}

func NewDefaultNavigator(z float64, xoff, yoff int) DefaultNavigator {
	return DefaultNavigator{z, image.Pt(xoff, yoff)}
}

func (n *DefaultNavigator) Transform(p image.Point) (float64, float64) {
	o := p.Add(n.offset)
	return float64(o.X) / n.z, float64(o.Y) / n.z
}

func (n *DefaultNavigator) Translate(offset image.Point) {
	n.offset = offset
}

func (n *DefaultNavigator) Zoom(z float64) {
	n.z = z
}
