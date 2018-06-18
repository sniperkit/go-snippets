package main

import (
	"fmt"
	"net/http"
)

func OurLoggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(*r.URL)
		h.ServeHTTP(w, r) // http.Handler 需要实现 ServeHTTP(ResponseWriter, *Request) 这个方法
	})
}

func main() {
	fileHandler := http.FileServer(http.Dir("/Users/path/code"))
	wrappedHandler := OurLoggingHandler(fileHandler)
	http.ListenAndServe(":9000", wrappedHandler)
}
