// +build !appengine !appenginedev
package main

import (
	"code.google.com/p/go-fracserv/fracserv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

var (
	port     = flag.Int("port", 8000, "webserver listen port")
	cacheDir = flag.String("cacheDir", "/tmp/fractals",
		"directory to store rendered tiles. Directory must exist")
	staticDir = flag.String("staticDir", "static",
		"directory containg statically served web page resources, i.e. javascript, css and image files")
)

func main() {
	flag.Parse()

	s := *staticDir
	_, err := os.Stat(s)
	if os.IsNotExist(err) {
		log.Fatalf("Directory %s not found, please run for directory containing %s\n", s, s)
	}
	// Setup handler for js, img, css files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(s))))

	fmt.Printf("Listening on:\n")
	host, err := os.Hostname()
	if err != nil {
		log.Fatal("Failed to get hostname from os:", err)
	}
	fmt.Printf("  http://%s:%d/\n", host, *port)

	go loadCache()

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}

func loadCache() {
	if *fracserv.DisableCache {
		log.Printf("Caching disable, not loading cache")
		return
	}

	files, err := filepath.Glob(path.Join(*cacheDir, "*/*"))
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
		cacher := fracserv.CachedPng{
			Timestamp: s.ModTime(),
			Bytes: b,
		}
		fracserv.PngCache.Add(path.Join(path.Base(path.Dir(fn)), path.Base(fn)), cacher)
	}
	log.Printf("Loaded %d cached tiles.", len(files))
}

func savePngFromCache(cacheKey string) {
	cacher, ok := fracserv.PngCache.Get(cacheKey)
	if !ok {
		log.Printf("Attempt to save %q to disk, but image not in cache",
			cacheKey)
		return
	}

	cachefn := path.Join(*cacheDir, cacheKey)
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
	cp := cacher.(fracserv.CachedPng)
	outf.Write(cp.Bytes)
	outf.Close()

	err = os.Chtimes(cachefn, cp.Timestamp, cp.Timestamp)
	if err != nil {
		log.Printf("Error setting atime and mtime on %q: %s", cachefn, err)
	}
}
