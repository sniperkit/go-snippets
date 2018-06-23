package httputil

import (
	"net/http"
	"time"

	"github.com/urfave/negroni"
)

type RequestLogger struct {
}

func NewRequestLogger() *RequestLogger {
	return &RequestLogger{}
}

func (self *RequestLogger) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(rw, req)

	response := rw.(negroni.ResponseWriter)
	status := response.Status()
	duration := time.Since(start)

	if status < 400 {
		Logger.Debugf("[HTTP %d] %s to %v took %v", status, req.Method, req.URL, duration)
	} else {
		Logger.Warningf("[HTTP %d] %s to %v took %v", status, req.Method, req.URL, duration)
	}
}
