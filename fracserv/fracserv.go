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
	"code.google.com/p/go-fracserv/cache"
	"code.google.com/p/go-fracserv/fractal"
	"code.google.com/p/go-fracserv/fractal/debug"
	"code.google.com/p/go-fracserv/fractal/julia"
	"code.google.com/p/go-fracserv/fractal/mandelbrot"
	"code.google.com/p/go-fracserv/fractal/solid"
	"flag"
	"fmt"
	"html/template"
	"image/png"
	"log"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"

	_ "net/http/pprof"
)

var factory map[string]func(o fractal.Options) (fractal.Fractal, error)
var PngCache cache.Cache

var (
	templateDir = flag.String("templateDir", "templates",
		"directory containing HTML pages and fragments")
	DisableCache = flag.Bool("disableCache", false,
		"disables all caching, ever requested rendered on demand")
)

type CachedPng struct {
	Timestamp time.Time
	Bytes     []byte
}

func (c CachedPng) Size() int {
	return len(c.Bytes)
}

func init() {
	factory = map[string]func(o fractal.Options) (fractal.Fractal, error){
		"debug":      debug.NewFractal,
		"solid":      solid.NewFractal,
		"mandelbrot": mandelbrot.NewFractal,
		"julia":      julia.NewFractal,
		//"lyapunov": lyapunov.NewFractal,
	}

	PngCache = *cache.NewCache()

	// Register a handler per known fractal type
	for k, _ := range factory {
		http.HandleFunc("/"+k, FracHandler)
	}
	// Catch-all handler, just serves homepage at "/", or 404s
	http.HandleFunc("/", IndexHander)
}

func drawFractalPage(w http.ResponseWriter, req *http.Request, fracType string) {
	t, err := template.ParseFiles(fmt.Sprintf("%s/%s.html", *templateDir, fracType))
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

func drawFractal(w http.ResponseWriter, req *http.Request, fracType string) {
	if *DisableCache {
		i, err := factory[fracType](fractal.Options{
			Values: req.URL.Query(),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		png.Encode(w, i)
		return
	}

	cacheKey := fsNameFromURL(req.URL)
	cacher, ok := PngCache.Get(cacheKey)
	if !ok {
		// No png in cache, create one
		i, err := factory[fracType](fractal.Options{
			Values: req.URL.Query(),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b := &bytes.Buffer{}
		png.Encode(b, i)
		cacher = CachedPng{time.Now(), b.Bytes()}
		PngCache.Add(cacheKey, cacher)

		// Async save image to disk
		// TODO make this a channel and serialize saving of images
		//go savePngFromCache(cacheKey)
	}

	cp := cacher.(CachedPng)

	// Set expire time
	req.Header.Set("Expires", time.Now().Add(time.Hour).Format(http.TimeFormat))
	// Using this instead of io.Copy, sets Last-Modified which helps given
	// the way the maps API makes lots of re-requests
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Last-Modified", cp.Timestamp.Format(http.TimeFormat))
	w.Header().Set("Expires",
		cp.Timestamp.Add(time.Hour).Format(http.TimeFormat))
	w.Write(cp.Bytes)
}

func FracHandler(w http.ResponseWriter, req *http.Request) {
	fracType := req.URL.Path[1:]
	if fracType != "" {
		//log.Println("Found fractal type", fracType)

		if len(req.URL.Query()) != 0 {
			drawFractal(w, req, fracType)
		} else {
			drawFractalPage(w, req, fracType)
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
	err = t.Execute(w, factory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
