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
package main

import (
	"bytes"
	"code.google.com/p/go-fracserv/cache"
	"flag"
	"fmt"
	"code.google.com/p/go-fracserv/fractal"
	"code.google.com/p/go-fracserv/fractal/debug"
	"code.google.com/p/go-fracserv/fractal/julia"
	"code.google.com/p/go-fracserv/fractal/mandelbrot"
	"code.google.com/p/go-fracserv/fractal/solid"
	"html/template"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "net/http/pprof"
)

var factory map[string]func(o fractal.Options) (fractal.Fractal, error)
var port string
var cacheDir string
var disableCache bool
var pngCache cache.Cache

type cachedPng struct {
	Timestamp time.Time
	Bytes     []byte
}

func (c cachedPng) Size() int {
	return len(c.Bytes)
}

func init() {
	flag.StringVar(&port, "port", "8000", "webserver listen port")
	flag.StringVar(&cacheDir, "cacheDir", "/tmp/fractals",
		"directory to store rendered tiles. Directory must exist")
	flag.BoolVar(&disableCache, "disableCache", false,
		"never serve from disk cache")
	flag.Parse()

	factory = map[string]func(o fractal.Options) (fractal.Fractal, error){
		"debug":      debug.NewFractal,
		"solid":      solid.NewFractal,
		"mandelbrot": mandelbrot.NewFractal,
		"julia":      julia.NewFractal,
		//"glynn": glynn.NewFractal,
		//"lyapunov": lyapunov.NewFractal,
	}

	pngCache = *cache.NewCache()
}

func main() {
	fmt.Printf("Listening on:\n")
	host, err := os.Hostname()
	if err != nil {
		log.Fatal("Failed to get hostname from os:", err)
	}
	fmt.Printf("  http://%s:%s/\n", host, port)

	s := "static/"
	_, err = os.Open(s)
	if os.IsNotExist(err) {
		log.Fatalf("Directory %s not found, please run for directory containing %s\n", s, s)
	}

	go loadCache()

	// Setup handler for js, img, css files
	http.Handle("/"+s, http.StripPrefix("/"+s, http.FileServer(http.Dir(s))))
	// Register a handler per known fractal type
	for k, _ := range factory {
		http.HandleFunc("/"+k, FracHandler)
	}
	// Catch-all handler, just serves homepage at "/", or 404s
	http.HandleFunc("/", IndexHander)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func loadCache() {
	if disableCache {
		log.Printf("Caching disable, not loading cache")
		return
	}

	files, err := filepath.Glob(cacheDir + "/*/*")
	if err != nil {
		log.Printf("Error globing cachedir %q: %s", cacheDir, err)
	}

	for idx, fn := range files {
		if idx%1000 == 0 {
			log.Printf("Loading %d/%d cached tiles...", idx, len(files))
		}

		s, err := os.Stat(fn)
		if err != nil {
			log.Printf("Error stating tile %q: %s", fn, err)
			continue
		}

		b, err := ioutil.ReadFile(fn)
		if err != nil {
			log.Printf("Error reading tile %q: %s", fn, err)
		}
		cacher := cachedPng{s.ModTime(), b}
		pngCache.Add(path.Join(path.Base(path.Dir(fn)), path.Base(fn)), cacher)
	}
	log.Printf("Loaded %d cached tiles.", len(files))
}

func drawFractalPage(w http.ResponseWriter, req *http.Request, fracType string) {
	t, err := template.ParseFiles(fmt.Sprintf("templates/%s.html", fracType))
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

func savePngFromCache(cacheKey string) {
	cacher, ok := pngCache.Get(cacheKey)
	if !ok {
		log.Printf("Attempt to save %q to disk, but image not in cache",
			cacheKey)
		return
	}

	cachefn := cacheDir + cacheKey
	d := path.Dir(cachefn)
	if _, err := os.Stat(d); err != nil {
		log.Printf("Creating cache dir for %q", d)
		err = os.Mkdir(d, 0700)
	}

	_, err := os.Stat(cachefn)
	if err == nil {
		log.Printf("Attempt to save %q to %q, but file already exists",
			cacheKey, cachefn)
		return
	}

	outf, err := os.OpenFile(cachefn, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open tile %q for save: %s", cachefn, err)
		return
	}
	cp := cacher.(cachedPng)
	outf.Write(cp.Bytes)
	outf.Close()

	err = os.Chtimes(cachefn, cp.Timestamp, cp.Timestamp)
	if err != nil {
		log.Printf("Error setting atime and mtime on %q: %s", cachefn, err)
	}
}

func drawFractal(w http.ResponseWriter, req *http.Request, fracType string) {
	if disableCache {
		i, err := factory[fracType](fractal.Options{req.URL.Query()})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		png.Encode(w, i)
		return
	}

	cacheKey := fsNameFromURL(req.URL)
	cacher, ok := pngCache.Get(cacheKey)
	if !ok {
		// No png in cache, create one
		i, err := factory[fracType](fractal.Options{req.URL.Query()})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b := &bytes.Buffer{}
		png.Encode(b, i)
		cacher = cachedPng{time.Now(), b.Bytes()}
		pngCache.Add(cacheKey, cacher)

		// Async save image to disk
		// TODO make this a channel and serialize saving of images
		go savePngFromCache(cacheKey)
	}

	cp := cacher.(cachedPng)

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

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, factory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
