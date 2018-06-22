package testing

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

type Transport struct {
	Message string
	Status  int
	Headers map[string][]string
}

func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp = &http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString(t.Message)),
		StatusCode: t.Status,
		Header:     http.Header(t.Headers),
	}
	return resp, nil
}

type ConditionalTransport struct {
	Transport
	CondFunc func(*http.Request) bool
}

func (t *ConditionalTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if !t.CondFunc(req) {
		return &http.Response{Body: nil, StatusCode: 500}, errors.New("condition failed")
	}
	return t.Transport.RoundTrip(req)
}
