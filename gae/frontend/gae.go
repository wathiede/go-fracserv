package appengine

import (
	"fmt"
	"net/http"
	_ "code.google.com/p/go-fracserv/fracserv"
)

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}
