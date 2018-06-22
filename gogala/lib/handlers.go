package lib

import (
	"compress/gzip"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func CacheHandler(days int, next http.Handler) http.Handler {

	if days < 1 {
		days = 1
	}
	age := days * 24 * 60 * 60 * 1000
	t := time.Now().Add(time.Duration(time.Hour * 24 * time.Duration(days)))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(age))
		w.Header().Set("Expires", t.Format(time.RFC1123Z))

		next.ServeHTTP(w, r)
	})
}

type GZipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w GZipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GZipHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gz := gzip.NewWriter(w)
		defer gz.Close()

		gw := GZipResponseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(gw, r)

	})
}
