package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type HandlerTestObj struct {
	Method  string
	Path    string
	In      interface{}
	Out     interface{}
	Handler http.HandlerFunc
}

var baseURL = "https://127.0.0.1/"

func runHandlerTest(t *testing.T, handlerTest HandlerTestObj) (w *httptest.ResponseRecorder, expectedOut, actualOut interface{}) {
	inBytes, err := json.Marshal(handlerTest.In)
	if err != nil {
		t.Fatal("Failed marshalling handlerTest.In:", err)
	}

	req := httptest.NewRequest(handlerTest.Method, baseURL+handlerTest.Path, bytes.NewBuffer(inBytes))
	w = httptest.NewRecorder()
	handlerTest.Handler(w, req)

	if handlerTest.Out == nil {
		return
	}

	// Handle expected output
	expectedBytes, err := json.Marshal(handlerTest.Out)
	if err != nil {
		t.Fatal("Failed marshalling handlerTest.Out:", err)
	}
	expectedBuf := bytes.NewBuffer(expectedBytes)
	expectedValue := reflect.New(reflect.TypeOf(handlerTest.Out))
	json.NewDecoder(expectedBuf).Decode(expectedValue.Interface())
	expectedOut = expectedValue.Interface()

	// Handle actual output
	actualValue := reflect.New(reflect.TypeOf(handlerTest.Out))
	json.NewDecoder(w.Body).Decode(actualValue.Interface())
	actualOut = actualValue.Interface()

	return
}
