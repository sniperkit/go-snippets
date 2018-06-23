package default_errors_service

import (
	"github.com/go-zero-boilerplate/ddd/interfaces/errors"
)

type service struct {
	httpClientErrorCreator errors.HttpClientErrorCreator
}

func (s *service) CreateClientError(statusCode int, statusText string) *errors.ClientError {
	return &errors.ClientError{statusCode, statusText}
}

func (s *service) HttpClientErrorCreator() errors.HttpClientErrorCreator {
	return s.httpClientErrorCreator
}

func New(httpClientErrorCreator errors.HttpClientErrorCreator) errors.Service {
	return &service{
		httpClientErrorCreator,
	}
}
