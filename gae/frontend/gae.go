package appengine

import (
	"fmt"
	"net/http"
	_ "code.google.com/p/go-fracserv/fracserv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}
