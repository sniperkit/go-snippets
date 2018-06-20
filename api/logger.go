package api

import (
	"net/http"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/logfmt"
	"github.com/apex/log/handlers/multi"
	"github.com/codegangsta/negroni"
)

type handler struct {
	logfmt  log.Handler
	kinesis log.Handler
}

// NewLogger is an exported function that returns a new log handler, with apex/log
func NewLogger() *handler {
	l := &handler{
		logfmt: logfmt.New(os.Stderr),
	}
	log.SetHandler(multi.New(
		l.logfmt,
	))
	return l
}

func (h *handler) LogThisDegub(fields log.Fields, level string, info string) {
	ctx := log.WithFields(fields)
	ctx.Debug(info)
}

func (h *handler) LogThis(fields log.Fields, level string, info string) {
	ctx := log.WithFields(fields)
	ctx.Info(info)
}

/*
INFO[0007] request                   agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/601.5.8 (KHTML, like Gecko) Version/9.1 Safari/601.5.8
elapsed=132.912µs method=GET path=/health proto=HTTP/1.1 remoteAddr=[::1]:57866 size=16 status=200 time=23/Jan/2016:13:03:44 -0200
*/

func (h *handler) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	t := time.Now()
	next(rw, r)
	res := rw.(negroni.ResponseWriter)
	ctx := log.WithFields(log.Fields{
		"elapsed":    time.Since(t),
		"remoteAddr": r.RemoteAddr,
		"method":     r.Method,
		"path":       r.URL.Path,
		"protocol":   r.Proto,
		"time":       t.Format("02/Jan/2006:15:04:05 -0700"),
		"status":     res.Status(),
		"size":       res.Size(),
		"agent":      r.UserAgent(),
	})
	ctx.Info("request")
}
