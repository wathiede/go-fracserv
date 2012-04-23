package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	//"path/filepath"
)

var port string

func init() {
	flag.StringVar(&port, "port", "8000", "webserver listen port")
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

func IndexServer(w http.ResponseWriter, req *http.Request) {
	fracType := req.URL.Path[1:]
	if fracType != "" {
		log.Println("Found fractal type", fracType)
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
