package httputil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type EncoderFunc func(interface{}) (io.Reader, error)
type DecoderFunc func(io.Reader, interface{}) error
type ErrorDecoderFunc func(*http.Response) error
type InterceptRequestFunc func(*http.Request) (interface{}, error)
type InterceptResponseFunc func(*http.Response, interface{}) error

func JSONEncoder(in interface{}) (io.Reader, error) {
	if data, err := json.Marshal(in); err == nil {
		return bytes.NewBuffer(data), nil
	} else {
		return nil, err
	}
}

func JSONDecoder(in io.Reader, out interface{}) error {
	return json.NewDecoder(in).Decode(out)
}
