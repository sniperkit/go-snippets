package errors

type ClientError struct {
	StatusCode int
	StatusText string
}
