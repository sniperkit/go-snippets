package httpmiddleware

import (
	"net/http"
)

type SingleHost struct {
	handler     http.Handler // the http.Handler interface
	allowedHost string
}

// ServeHTTP is the middleware
func (this *SingleHost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	println(r.Host)
	if r.Host == this.allowedHost {
		this.handler.ServeHTTP(w, r)
	} else {
		w.WriteHeader(403)
	}
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world!"))
}

func main() {
	single := &SingleHost{
		handler:     http.HandlerFunc(myHandler), // HandlerFunc is a type and the type is a function
		allowedHost: "localhost:8080",
	}
	http.ListenAndServe(":8080", single)
}
