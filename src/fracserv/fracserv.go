package main

import (
	"flag"
	"fmt"
	"fractal"
	"fractal/debug"
	"fractal/lyapunov"
	"fractal/mandelbrot"
	"fractal/solid"
	"html/template"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	//"path/filepath"
)

var factory map[string]func(o fractal.Options) (fractal.Fractal, error)
var port string


func init() {
	flag.StringVar(&port, "port", "8000", "webserver listen port")
	factory = map[string]func(o fractal.Options) (fractal.Fractal, error){
		"debug": debug.NewFractal,
		"solid": solid.NewFractal,
		"mandelbrot": mandelbrot.NewFractal,
		"lyapunov": lyapunov.NewFractal,
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

	http.Handle("/"+s, http.StripPrefix("/"+s, http.FileServer(http.Dir(s))))
	http.HandleFunc("/", IndexServer)
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
	cachefn := "cache" + req.URL.RequestURI()
	f, err := os.Open(cachefn)
	if err != nil {
		// No file, create one
		i, err := factory[fracType](fractal.Options{req.URL.Query()})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f, err := os.OpenFile(cachefn, os.O_CREATE|os.O_WRONLY, 0644)
		var mw io.Writer
		if err != nil {
			log.Printf("Failed to save tile: %s", err)
			mw = io.MultiWriter(w)
		} else {
			defer f.Close()
			mw = io.MultiWriter(w, f)
		}

		png.Encode(mw, i)
	} else {
		defer f.Close()
		log.Printf("cache hit: %q", cachefn)
		io.Copy(w, f)
	}
}

func IndexServer(w http.ResponseWriter, req *http.Request) {
	fracType := req.URL.Path[1:]
	if fracType != "" {
		//log.Println("Found fractal type", fracType)

		if len(req.URL.Query()) != 0 {
			drawFractal(w, req, fracType)
		} else {
			drawFractalPage(w, req, fracType)
		}
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
