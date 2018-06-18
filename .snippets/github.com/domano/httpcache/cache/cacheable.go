package cache

import (
	"bytes"
	"net/http"
)

type Cacheable interface {
	Key() string
	Size() int
}
type CacheableResponseWriter interface {
	http.ResponseWriter
	Cacheable
}

type cacheableResponseWriter struct {
	http.ResponseWriter
	key     string
	header  int
	buffer  *bytes.Buffer
	buffErr error
}

func NewCachableResponseWriter(key string, rw http.ResponseWriter) CacheableResponseWriter {
	return &cacheableResponseWriter{ResponseWriter: rw, key: key, buffer: &bytes.Buffer{}}
}

func (hc *cacheableResponseWriter) Key() string {
	return hc.key
}

func (hc *cacheableResponseWriter) Size() int {
	return hc.buffer.Len()
}

func (hc *cacheableResponseWriter) Write(p []byte) (int, error) {
	if hc.buffer.Len() < 0 || hc.buffErr == nil {
		_, err := hc.buffer.Write(p)
		hc.buffErr = err
	}
	return hc.ResponseWriter.Write(p)
}
