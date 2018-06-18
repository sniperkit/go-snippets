package errors

type HttpClientErrorCreator interface {
	BadRequest(statusText string) *ClientError
	Unauthorized(statusText string) *ClientError
	PaymentRequired(statusText string) *ClientError
	Forbidden(statusText string) *ClientError
	NotFound(statusText string) *ClientError
	MethodNotAllowed(statusText string) *ClientError
	NotAcceptable(statusText string) *ClientError
	ProxyAuthRequired(statusText string) *ClientError
	RequestTimeout(statusText string) *ClientError
	Conflict(statusText string) *ClientError
	Gone(statusText string) *ClientError
	LengthRequired(statusText string) *ClientError
	PreconditionFailed(statusText string) *ClientError
	RequestEntityTooLarge(statusText string) *ClientError
	RequestURITooLong(statusText string) *ClientError
	UnsupportedMediaType(statusText string) *ClientError
	RequestedRangeNotSatisfiable(statusText string) *ClientError
	ExpectationFailed(statusText string) *ClientError
	Teapot(statusText string) *ClientError
	InternalServerError(statusText string) *ClientError
	NotImplemented(statusText string) *ClientError
	BadGateway(statusText string) *ClientError
	ServiceUnavailable(statusText string) *ClientError
	GatewayTimeout(statusText string) *ClientError
	HTTPVersionNotSupported(statusText string) *ClientError
}
