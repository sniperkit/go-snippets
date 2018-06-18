package default_http_client_error_creator

import (
	"net/http"
	"strings"

	"github.com/go-zero-boilerplate/ddd/interfaces/errors"
)

type creator struct{}

func (c *creator) create(statusCode int, statusText string) *errors.ClientError {
	if strings.TrimSpace(statusText) == "" {
		return &errors.ClientError{statusCode, http.StatusText(statusCode)}
	}
	return &errors.ClientError{statusCode, statusText}
}

func (c *creator) BadRequest(statusText string) *errors.ClientError {
	return c.create(http.StatusBadRequest, statusText)
}

func (c *creator) Unauthorized(statusText string) *errors.ClientError {
	return c.create(http.StatusUnauthorized, statusText)
}

func (c *creator) PaymentRequired(statusText string) *errors.ClientError {
	return c.create(http.StatusPaymentRequired, statusText)
}

func (c *creator) Forbidden(statusText string) *errors.ClientError {
	return c.create(http.StatusForbidden, statusText)
}

func (c *creator) NotFound(statusText string) *errors.ClientError {
	return c.create(http.StatusNotFound, statusText)
}

func (c *creator) MethodNotAllowed(statusText string) *errors.ClientError {
	return c.create(http.StatusMethodNotAllowed, statusText)
}

func (c *creator) NotAcceptable(statusText string) *errors.ClientError {
	return c.create(http.StatusNotAcceptable, statusText)
}

func (c *creator) ProxyAuthRequired(statusText string) *errors.ClientError {
	return c.create(http.StatusProxyAuthRequired, statusText)
}

func (c *creator) RequestTimeout(statusText string) *errors.ClientError {
	return c.create(http.StatusRequestTimeout, statusText)
}

func (c *creator) Conflict(statusText string) *errors.ClientError {
	return c.create(http.StatusConflict, statusText)
}

func (c *creator) Gone(statusText string) *errors.ClientError {
	return c.create(http.StatusGone, statusText)
}

func (c *creator) LengthRequired(statusText string) *errors.ClientError {
	return c.create(http.StatusLengthRequired, statusText)
}

func (c *creator) PreconditionFailed(statusText string) *errors.ClientError {
	return c.create(http.StatusPreconditionFailed, statusText)
}

func (c *creator) RequestEntityTooLarge(statusText string) *errors.ClientError {
	return c.create(http.StatusRequestEntityTooLarge, statusText)
}

func (c *creator) RequestURITooLong(statusText string) *errors.ClientError {
	return c.create(http.StatusRequestURITooLong, statusText)
}

func (c *creator) UnsupportedMediaType(statusText string) *errors.ClientError {
	return c.create(http.StatusUnsupportedMediaType, statusText)
}

func (c *creator) RequestedRangeNotSatisfiable(statusText string) *errors.ClientError {
	return c.create(http.StatusRequestedRangeNotSatisfiable, statusText)
}

func (c *creator) ExpectationFailed(statusText string) *errors.ClientError {
	return c.create(http.StatusExpectationFailed, statusText)
}

func (c *creator) Teapot(statusText string) *errors.ClientError {
	return c.create(http.StatusTeapot, statusText)
}

func (c *creator) InternalServerError(statusText string) *errors.ClientError {
	return c.create(http.StatusInternalServerError, statusText)
}

func (c *creator) NotImplemented(statusText string) *errors.ClientError {
	return c.create(http.StatusNotImplemented, statusText)
}

func (c *creator) BadGateway(statusText string) *errors.ClientError {
	return c.create(http.StatusBadGateway, statusText)
}

func (c *creator) ServiceUnavailable(statusText string) *errors.ClientError {
	return c.create(http.StatusServiceUnavailable, statusText)
}

func (c *creator) GatewayTimeout(statusText string) *errors.ClientError {
	return c.create(http.StatusGatewayTimeout, statusText)
}

func (c *creator) HTTPVersionNotSupported(statusText string) *errors.ClientError {
	return c.create(http.StatusHTTPVersionNotSupported, statusText)
}

func New() errors.HttpClientErrorCreator {
	return &creator{}
}
