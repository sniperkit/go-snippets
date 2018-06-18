package rest

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/bradberger/context"
)

var (
	ErrBodyIsEmpty        = errors.New("body is empty")
	ErrUnknownContentType = errors.New("could not determinte content type")
)

type TestRequest struct {
	Writer  *httptest.ResponseRecorder
	Request *http.Request
	Context context.Context
}

func (r *TestRequest) Decode(data interface{}) error {
	if r.Writer.Body == nil {
		return ErrBodyIsEmpty
	}

	switch r.Request.Header.Get("Content-Type") {
	case "application/json":
		return json.NewDecoder(r.Writer.Body).Decode(data)
	case "application/xml":
		return xml.NewDecoder(r.Writer.Body).Decode(data)
	default:
		return ErrUnknownContentType
	}
}

func (r *TestRequest) Do(h AppHandler) error {
	defer r.Writer.Flush()
	return h(r.Context)
}

// TestPostJSON bootstraps a POST request to the given urlStr with encoded JSON data as the request body
func TestPostJSON(urlStr string, data interface{}) (req *TestRequest, err error) {

	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(data); err != nil {
		return
	}

	return TestPost(urlStr, "application/json", &buf)
}

// TestPostXML bootstraps a POST request to the given urlStr with encoded XML data as the request body
func TestPostXML(urlStr string, data interface{}) (req *TestRequest, err error) {
	var buf bytes.Buffer
	if err = xml.NewEncoder(&buf).Encode(data); err != nil {
		return
	}

	return TestPost(urlStr, "application/xml", &buf)
}

// TestPostForm bootstraps a POST request to the given urlStr with the form as the request body
func TestPostForm(urlStr string, form url.Values) (req *TestRequest, err error) {
	return TestPost(urlStr, "application/x-www-form-urlencoded", bytes.NewBufferString(form.Encode()))
}

// TestPost bootstraps a POST request to the given urlStr with the given request body
func TestPost(urlStr string, contentType string, body io.Reader) (req *TestRequest, err error) {

	r, err := http.NewRequest("POST", urlStr, body)
	r.Header.Set("Content-Type", contentType)
	if err != nil {
		return
	}

	return NewTestRequest(r), nil
}

func NewTestRequest(req *http.Request) *TestRequest {
	w := httptest.NewRecorder()
	return &TestRequest{Writer: w, Context: Init(w, req), Request: req}
}
