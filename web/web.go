package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"net/http"
)

var port string

func init() {
	flag.StringVar(&port, "port", "8000", "webserver listen port")
}

func main() {
	fmt.Printf("Listening on:\n")
	host, err := os.Hostname()
	if err != nil {
		log.Fatalf("Failed to get hostname from os: %s\n", err)
	}
	fmt.Printf("  http://%s:%s/\n", host, port)

	http.HandleFunc("/", IndexServer)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}

// hello world, the web server
func IndexServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

