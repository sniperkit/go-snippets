package mux

import (
	"github.com/verdverm/tview"
)

type HandlerFunc func(*Request) (tview.Primitive, *Request, error)

type Handler interface {
	Serve(*Request) (tview.Primitive, *Request, error)
}

type DefaultHandler struct {
	handler HandlerFunc
}

func NewDefaultHandler(handle HandlerFunc) *DefaultHandler {
	return &DefaultHandler{
		handler: handle,
	}
}

func (H *DefaultHandler) Serve(req *Request) (tview.Primitive, *Request, error) {
	return H.handler(req)
}
