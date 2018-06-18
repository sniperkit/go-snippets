package rest

import (
	"net/http"

	"github.com/bradberger/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// OnPanic when defined, will handle any panics from resulting funcs
var OnPanic AppHandler

// Namespace enables setting custom namespace
var Namespace func(ctx context.Context) (string, error) = func(ctx context.Context) (string, error) {
	return "", nil
}

// OnUnauthorized handles the response for 401 responses
var OnUnauthorized http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

// OnError is a custom error handler definition. By default it simple returns an http.Error()
// with the given code and status text.
var OnError func(ctx context.Context, code int, err error) = func(ctx context.Context, code int, err error) {
	Errorf(ctx, "error: %v", err)
	http.Error(ResponseWriter(ctx), err.Error(), code)
}

// GetErrorCode allows customizing of the http.StatusCode for any given error. If the code has already
// been set with SetCode() then it will not be overwritten by this function.
var GetErrorCode func(ctx context.Context, err error) int = func(ctx context.Context, err error) (code int) {

	exist := GetCode(ctx)
	if exist != http.StatusOK {
		return exist
	}

	switch err.(type) {
	case StatusCode:
		return err.(StatusCode).Code()
	}

	switch {
	case err == datastore.ErrNoSuchEntity:
		code = http.StatusNotFound
	case appengine.IsOverQuota(err):
		code = http.StatusTooManyRequests
	case appengine.IsTimeoutError(err):
		code = http.StatusGatewayTimeout
	default:
		code = http.StatusInternalServerError
	}
	return
}
