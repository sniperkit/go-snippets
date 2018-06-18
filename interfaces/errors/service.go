package errors

type Service interface {
	CreateClientError(statusCode int, statusText string) *ClientError
	HttpClientErrorCreator() HttpClientErrorCreator
}
