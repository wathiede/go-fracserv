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
package fracserv

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/wathiede/go-fracserv/cache"
	"github.com/wathiede/go-fracserv/fractal"
	// Handy tile type that draws text for query parameters on tile
	_ "github.com/wathiede/go-fracserv/fractal/debug"
	_ "github.com/wathiede/go-fracserv/fractal/julia"
	_ "github.com/wathiede/go-fracserv/fractal/mandelbrot"
	_ "github.com/wathiede/go-fracserv/fractal/perlin"
	// Illustrates basic parameter handling
	_ "github.com/wathiede/go-fracserv/fractal/solid"
)

var ImageCache cache.Cache

var (
	templateDir = flag.String("templateDir", "templates",
		"directory containing HTML pages and fragments")
	DisableCache = flag.Bool("disableCache", false,
		"disables all caching, ever requested rendered on demand")
	jpegTiles = flag.Bool("jpegTiles", false,
		"render jpeg instead of png tiles")
)

type CachedImage struct {
	Timestamp time.Time
	Bytes     []byte
}

func (c CachedImage) Size() int {
	return len(c.Bytes)
}

func encodeImage(w io.Writer, m image.Image) error {
	var err error
	if *jpegTiles {
		err = jpeg.Encode(w, m, &jpeg.Options{
			Quality: 100,
		})
	} else {
		err = png.Encode(w, m)
	}
	return err
}

func init() {
	ImageCache = *cache.NewCache()

	// Register a handler per known fractal type
	fractal.Do(func(name string, newFunc fractal.FractalNew) {
		log.Print("Registering ", name)
		http.HandleFunc("/"+name, FracHandlerNew(name, newFunc))
	})
	// Catch-all handler, just serves homepage at "/", or 404s
	http.HandleFunc("/", IndexHander)
}

func drawFractalPage(w http.ResponseWriter, req *http.Request, name string) {
	t, err := template.ParseFiles(fmt.Sprintf("%s/%s.html", *templateDir, name))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func fsNameFromURL(u *url.URL) string {
	fn := strings.TrimLeft(u.Path, "/") + "/"
	keys := []string{}
	q := u.Query()

	for k := range q {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	p := []string{}
	for _, k := range keys {
		p = append(p, k+"="+q[k][0])
	}

	return fn + strings.Join(p, ",")
}

func drawFractal(w http.ResponseWriter, req *http.Request, newFunc fractal.FractalNew) {
	if *DisableCache {
		i, err := newFunc(fractal.Options{Values: req.URL.Query()})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = encodeImage(w, i)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	cacheKey := fsNameFromURL(req.URL)
	cacher, ok := ImageCache.Get(cacheKey)
	if !ok {
		// No png in cache, create one
		i, err := newFunc(fractal.Options{Values: req.URL.Query()})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b := &bytes.Buffer{}
		err = encodeImage(b, i)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cacher = CachedImage{time.Now(), b.Bytes()}
		ImageCache.Add(cacheKey, cacher)

		// Async save image to disk
		// TODO make this a channel and serialize saving of images
		//go savePngFromCache(cacheKey)
	}

	cp := cacher.(CachedImage)

	// Set expire time
	w.Header().Set("Expires", time.Now().Add(time.Hour).Format(http.TimeFormat))
	// Using this instead of io.Copy, sets Last-Modified which helps given
	// the way the maps API makes lots of re-requests
	if *jpegTiles {
		w.Header().Set("Content-Type", "image/jpeg")
	} else {
		w.Header().Set("Content-Type", "image/png")
	}
	w.Header().Set("Last-Modified", cp.Timestamp.Format(http.TimeFormat))
	w.Header().Set("Expires",
		cp.Timestamp.Add(time.Hour).Format(http.TimeFormat))
	w.Write(cp.Bytes)
}

func FracHandlerNew(name string, newFunc fractal.FractalNew) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fracType := req.URL.Path[1:]
		if fracType != "" {
			//log.Println("Found fractal type", fracType)

			if len(req.URL.Query()) != 0 {
				drawFractal(w, req, newFunc)
			} else {
				drawFractalPage(w, req, name)
			}
		}
	}
}

func IndexHander(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		log.Println("404:", req.URL)
		http.NotFound(w, req)
		return
	}

	t, err := template.ParseFiles(path.Join(*templateDir, "index.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fractals := []string{}
	fractal.Do(func(name string, _ fractal.FractalNew) {
		fractals = append(fractals, name)
	})
	sort.Strings(fractals)

	err = t.Execute(w, fractals)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
