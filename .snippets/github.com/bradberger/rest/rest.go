// Package rest implements a simple REST microframework designed for Google App Engine
package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/bradberger/context"
	"github.com/gorilla/mux"
)

// Context key standardized definitions
var (
	ContextKeyNamepace       context.Key = "namespace"
	ContextKeyRequestBody    context.Key = "request.body"
	ContextKeyRequestVars    context.Key = "request.vars"
	ContextKeyResponseWriter context.Key = "http.responsewriter"
	ContextKeyRequest        context.Key = "request"
	ContextKeyResponseCode   context.Key = "response.code"
	ContextKeyInitialized    context.Key = "initialized"
	ContextKeyResponseBody   context.Key = "response.body"
	ContextKeyEnvironment    context.Key = "environment"
)

// AppHandler is the wrapper for all HTTP requests. It provides a valid context, authorization information, and route parameters.
// The returned interface is written to the response, unless an error is returned.
type AppHandler func(ctx context.Context) error

// Handler is a chainable set of AppHandler middleware funcs
func Handler(fn ...AppHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := Init(w, r)
		for i := range fn {
			if err := fn[i](ctx); err != nil {
				OnError(ctx, GetErrorCode(ctx, err), err)
				return
			}
		}
	})
}

// Decode decodes the http request body in JSON format into the dst variable.
func Decode(ctx context.Context, dst interface{}) error {
	b := ctx.Value(ContextKeyRequestBody)
	if b == nil {
		return errors.New("no request body")
	}
	return json.Unmarshal(b.([]byte), dst)
}

// Body returns the body of the request as a byte slice
func Body(ctx context.Context) []byte {
	b := ctx.Value(ContextKeyRequestBody).([]byte)
	if b == nil {
		return []byte{}
	}
	return b
}

// BodyString returns the body of the request as a string
func BodyString(ctx context.Context) string {
	return string(Body(ctx))
}

// FormValue returns the form value (or mux.Vars value) from the request
func FormValue(ctx context.Context, key string) string {
	v := ctx.Value(ContextKeyRequestVars)
	if v == nil {
		return ""
	}
	return v.(map[string]string)[key]
}

// ResponseWriter returns the response writer for the given context
func ResponseWriter(ctx context.Context) http.ResponseWriter {
	return ctx.Value(ContextKeyResponseWriter).(http.ResponseWriter)
}

// Request returns the http.Request asso
func Request(ctx context.Context) *http.Request {
	return ctx.Value(ContextKeyRequest).(*http.Request)
}

// Header returns the ResponseWriter headers
func Header(ctx context.Context) http.Header {
	return ResponseWriter(ctx).Header()
}

// Headers returns the http.Request headers
func Headers(ctx context.Context) http.Header {
	return Request(ctx).Header
}

// GetCode returns the http response code associated with the request
func GetCode(ctx context.Context) int {
	if v := ctx.Value(context.ContextKeyResponseCode); v != nil {
		return v.(int)
	}
	return http.StatusOK
}

// FormFile matches the "net/http".Request.FormFile api
func FormFile(ctx context.Context, key string) (multipart.File, *multipart.FileHeader, error) {
	return Request(ctx).FormFile(key)
}

func setRequest(ctx context.Context, r *http.Request) context.Context {
	return setValue(ctx, ContextKeyRequest, r)
}

func setValue(ctx context.Context, key, val interface{}) context.Context {
	return context.WithValue(ctx, key, val)
}

func setVars(ctx context.Context) (context.Context, error) {

	v := map[string]string{}
	r := Request(ctx)

	err := r.ParseForm()
	if err != nil {
		return ctx, err
	}

	for i := range r.Form {
		v[i] = r.FormValue(i)
	}
	for i := range mux.Vars(r) {
		v[i] = mux.Vars(r)[i]
	}

	ctx = setValue(ctx, ContextKeyRequestVars, v)
	return ctx, nil
}

func setBody(ctx context.Context) (context.Context, error) {
	r := Request(ctx)
	bodyBytes, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	// Reset the body so it can be read again.
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	if err != nil {
		return ctx, err
	}
	return setValue(ctx, ContextKeyRequestBody, bodyBytes), nil
}

func setWriter(ctx context.Context, w http.ResponseWriter) context.Context {
	return setValue(ctx, ContextKeyResponseWriter, w)
}

// Init returns a context with the reader, writer, and other context variables set
func Init(w http.ResponseWriter, r *http.Request) context.Context {

	ctx := context.NewContext(r)
	ctx = setRequest(ctx, r)
	ctx = setWriter(ctx, w)

	// TODO Figure out how to handle errors here
	ctx, _ = setNamespace(ctx)
	ctx, _ = setVars(ctx)
	ctx, _ = setBody(ctx)
	return ctx
}

// New creates a new initialized router
func New() *mux.Router {
	r := mux.NewRouter()
	r.KeepContext = true
	return r
}
