package main

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"

	"github.com/zenazn/goji"
)

func main() {

	// use gzip handle for every request
	goji.Use(gzipHandler)

	goji.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "helloworld..........")
	})

	goji.Serve()
}

// gzip handler

type gzipResponseWriter struct {
	*gzip.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	h := w.ResponseWriter.Header()
	if h.Get("Content-Type") == "" {
		h.Set("Content-Type", http.DetectContentType(b))
	}

	return w.Writer.Write(b)
}

func gzipHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}

		// gzip content
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")
		gw := gzip.NewWriter(w)
		defer gw.Close()

		w = &gzipResponseWriter{gw, w}

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
