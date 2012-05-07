package appengine

import (
	_ "code.google.com/p/go-fracserv/fracserv"
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}
