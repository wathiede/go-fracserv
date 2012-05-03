package main

import (
	"flag"
	"fmt"
	"fractal"
	"fractal/debug"
	"fractal/julia"
	"fractal/mandelbrot"
	"fractal/solid"
	"html/template"
	"image/png"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

var factory map[string]func(o fractal.Options) (fractal.Fractal, error)
var port string
var cacheDir string

func init() {
	flag.StringVar(&port, "port", "8000", "webserver listen port")
	flag.StringVar(&cacheDir, "cacheDir", "/tmp/fractals",
		"directory to store rendered tiles. Directory must exist")
	flag.Parse()

	factory = map[string]func(o fractal.Options) (fractal.Fractal, error){
		"debug": debug.NewFractal,
		"solid": solid.NewFractal,
		"mandelbrot": mandelbrot.NewFractal,
		"julia": julia.NewFractal,
		//"glynn": glynn.NewFractal,
		//"lyapunov": lyapunov.NewFractal,
	}
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

	// Setup handler for js, img, css files
	http.Handle("/"+s, http.StripPrefix("/"+s, http.FileServer(http.Dir(s))))
	// Register a handler per known fractal type
	for k, _ := range factory {
		http.HandleFunc("/" + k, FracHandler)
	}
	// Catch-all handler, just serves homepage at "/", or 404s
	http.HandleFunc("/", IndexHander)
	log.Fatal(http.ListenAndServe(":" + port, nil))
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

func drawFractal(w http.ResponseWriter, req *http.Request, fracType string) {
	cleanup := func(r rune) rune {
		switch r {
		case '?':
			return '/'
		case '&':
			return ','
		}
		return r
	}
	cachefn := cacheDir + strings.Map(cleanup, req.URL.RequestURI())
	d := path.Dir(cachefn)
	if _, err := os.Stat(d); err != nil {
		log.Printf("Creating cache dir for %q", d)
		err = os.Mkdir(d, 0700)
	}
	_, err := os.Stat(cachefn)
	if err != nil {
		// No file, create one
		i, err := factory[fracType](fractal.Options{req.URL.Query()})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		outf, err := os.OpenFile(cachefn, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Failed to open tile for save: %s", err)
			// Just send png from memory
			png.Encode(w, i)
			return
		}
		// Save to disk and serve below with http.ServeFile
		png.Encode(outf, i)
		outf.Close()
	}
	// TODO(wathiede): log cache hits as expvar

	// Set expire time
	req.Header.Set("Expires", time.Now().Add(time.Hour).Format(http.TimeFormat))
	// Using this instead of io.Copy, sets Last-Modified which helps given
	// the way the maps API makes lots of re-requests
	http.ServeFile(w, req, cachefn)
}

func FracHandler(w http.ResponseWriter, req *http.Request) {
	// TODO(wathiede): check fracType against keys in factory before handing
	// off to handler to prevent directory traversal exploit from malformed
	// urls.
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
