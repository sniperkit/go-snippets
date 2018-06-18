package responses

import (
	"github.com/go-zero-boilerplate/facades/tasks"
)

//NewResponse creates a new Response
func NewResponse(status Status, body interface{}, formatter Formatter) *Response {
	return &Response{
		Status:    status,
		Body:      body,
		Formatter: formatter,
	}
}

//NewResponseFromTaskError creates a new Response from a tasks.Error
func NewResponseFromTaskError(err *tasks.Error, formatter Formatter) *Response {
	return &Response{
		Status:    Status(err.GetStatusCode()),
		Body:      NewError(err.GetError().Error()),
		Formatter: formatter,
	}
}

//Response is the generic response that Resource methods (Get, Post, etc) will return
type Response struct {
	Status    Status
	Body      interface{}
	Formatter Formatter
}

func (r *Response) Response() error {
	return r.Formatter(r.Status.Int(), r.Body)
}
