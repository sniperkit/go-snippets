package router

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/eapache/go-resiliency/breaker"
	"github.com/pressly/chi"
	"github.com/prometheus/client_golang/prometheus"
)

//Example of router package, the router wrapper is based on pressly/chi

var (
	// for prometheus monitoring
	prometheusSummaryVec *prometheus.SummaryVec
)

func SetMonitoring(namespace string) {
	prometheusSummaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: namespace,
		Name:      "handler_request_milisecond",
		Help:      "Average of handler response time at one time",
	}, []string{"handler", "method", "httpcode"})
	if err := prometheus.Register(prometheusSummaryVec); err != nil {
		log.Printf("Failed to register prometheus metrics: %s", err.Error())
	}
}

// Router struct
type Router struct {
	opt Option
	r   chi.Router
	cb  *breaker.Breaker
}

// Option for router
type Option struct {
	Timeout time.Duration
	// Instrumenet white-box monitoring from the application, provide data about handler
	Instrument bool
}

// New router
func New(opt Option) *Router {
	chiRouter := chi.NewRouter()
	rtr := &Router{
		r:   chiRouter,
		opt: opt,
		// this is only an example
		cb: breaker.New(10, 2, time.Second*15),
	}
	return rtr
}

const (
	BreakerOpen = "rtr_breaker_open"
)

func OpenBreaker(r *http.Request) {
	ctx := context.WithValue(r.Context(), BreakerOpen, 1)
	r = r.WithContext(ctx)
}

func (rtr *Router) circuitbreaker(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := rtr.cb.Run(func() error {
			h(w, r)
			v := r.Context().Value(BreakerOpen)
			if v == nil || v != 1 {
				return nil
			}
			return errors.New("Circuit breaker open")
		})
		if err == nil || err != breaker.ErrBreakerOpen {
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Service unavailable"))
	}
}

// timeout middleware
// the timeout middleware should cover timeout budget
// this is needed because the budget will determine it should open a breaker or not
func (rtr *Router) timeout(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opt := rtr.opt
		// cancel context
		if opt.Timeout > 0 {
			ctx, cancel := context.WithTimeout(r.Context(), opt.Timeout*time.Second)
			defer cancel()
			r = r.WithContext(ctx)
		}

		doneChan := make(chan bool)
		go func() {
			h(w, r)
			doneChan <- true
		}()
		select {
		case <-r.Context().Done():
			// only an example response
			resp := map[string]interface{}{
				"errors": []string{"Request timed out"},
			}
			jsonResp, _ := json.Marshal(resp)
			w.WriteHeader(http.StatusRequestTimeout)
			// only an example response
			w.Write(jsonResp)
			return
		case <-doneChan:
			return
		}
	}
}

// responseWriterDelegator to delegate the current writer
// this is a 100% from prometheus delegator with some modification
// the modification is needed because namespace is required
type responseWriterDelegator struct {
	http.ResponseWriter

	// handler, method string
	status      int
	written     int64
	wroteHeader bool
}

func (r *responseWriterDelegator) WriteHeader(code int) {
	r.status = code
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseWriterDelegator) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	n, err := r.ResponseWriter.Write(b)
	r.written += int64(n)
	return n, err
}

func sanitizeStatusCode(status int) string {
	code := strconv.Itoa(status)
	return code
}

// monitor middleware provides white-box monitoring for application
// htp.ResponseWriter is delegated to create a custom metrics in this function
func (rtr *Router) monitor(pattern string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		delegator := &responseWriterDelegator{ResponseWriter: w}
		defer func(t time.Time) {
			method := r.Method
			statusCode := sanitizeStatusCode(delegator.status)
			if prometheusSummaryVec != nil {
				prometheusSummaryVec.With(prometheus.Labels{"handler": pattern, "method": method, "httpcode": statusCode}).Observe(time.Since(t).Seconds() * 1000)
			}
		}(time.Now())
		h(delegator, r)
	}
}

// URLParam get param from rest request
func URLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// URLParamFromCtx get param from rest request from context
func URLParamFromCtx(ctx context.Context, key string) string {
	return chi.URLParamFromCtx(ctx, key)
}

// Get function
func (rtr *Router) Get(pattern string, h http.HandlerFunc) {
	rtr.r.Get(pattern, rtr.monitor(pattern, rtr.timeout(h)))
}

// Post function
func (rtr *Router) Post(pattern string, h http.HandlerFunc) {
	rtr.r.Post(pattern, rtr.monitor(pattern, rtr.timeout(h)))
}

// Put function
func (rtr *Router) Put(pattern string, h http.HandlerFunc) {
	rtr.r.Put(pattern, rtr.monitor(pattern, rtr.timeout(h)))
}

// Delete function
func (rtr *Router) Delete(pattern string, h http.HandlerFunc) {
	rtr.r.Delete(pattern, rtr.monitor(pattern, rtr.timeout(h)))
}

// Patch function
func (rtr *Router) Patch(pattern string, h http.HandlerFunc) {
	rtr.r.Patch(pattern, rtr.monitor(pattern, rtr.timeout(h)))
}
