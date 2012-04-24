package fractal

import (
	"image"
	"net/url"
)

type Options struct {
	url.Values
}
type Fractal image.Image
