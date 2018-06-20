package resputil

import (
	"errors"
	"strings"
	"testing"

	"encoding/json"
	"io/ioutil"
	"net/http/httptest"

	e2 "github.com/pkg/errors"
)

func TestTests(t *testing.T) {
	testResponse := func(output interface{}) string {
		w := httptest.NewRecorder()
		JSON(w, output)
		body, _ := ioutil.ReadAll(w.Result().Body)
		return string(body)
	}

	testCase := func(name string, output interface{}, expected string) {
		got := strings.TrimSpace(testResponse(output))
		expected = strings.TrimSpace(expected)
		if got != expected {
			t.Errorf("test '%s', got %#v, expected %#v", name, got, expected)
		}
	}

	testCase("nil", nil, `{"response":false}`)
	testCase("bool true", true, `{"response":true}`)
	testCase("bool false", false, `{"response":false}`)
	testCase("string empty", "", `{"response":false}`)
	testCase("string", "string", `{"response":"string"}`)
	testCase("int zero", 0, `{"response":0}`)
	testCase("int non-zero", 1337, `{"response":1337}`)
	testCase("int sub-zero", -1, `{"response":-1}`)

	testCase("error nil", func() error {
		return nil
	}, `{"response":false}`)

	testCase("error", func() error {
		return errors.New("error response")
	}, `{"error":{"message":"error response"}}`)

	testCase("value + error", func() (interface{}, error) {
		return "string response", errors.New("error response")
	}, `{"error":{"message":"error response"}}`)

	testCase("empty value + error", func() (interface{}, error) {
		return "", errors.New("error response")
	}, `{"error":{"message":"error response"}}`)

	testCase("value + empty error", func() (interface{}, error) {
		return "string response", nil
	}, `{"response":"string response"}`)

	testCase("success default", Success(), `{"success":{"message":"OK"}}`)
	testCase("success default", OK(), `{"success":{"message":"OK"}}`)
	testCase("success custom", Success("string"), `{"success":{"message":"string"}}`)

	testCase("error ptr nil", new(error), `{"response":false}`)
	{
		err := errors.New("string")
		testCase("error ptr ok", &err, `{"error":{"message":"string"}}`)
	}
	testCase("error stdlib", errors.New("string"), `{"error":{"message":"string"}}`)
	{
		testCase("error stdlib nil", func() interface{} { return func() error { return nil }() }(), `{"response":false}`)
	}

	testCase("func json nil", func() ([]byte, error) { return json.Marshal(nil) }, `null`)
	testCase("func json false", func() ([]byte, error) { return json.Marshal(false) }, `false`)
	testCase("func json 0", func() ([]byte, error) { return json.Marshal(0) }, `0`)
	testCase("func json empty string", func() ([]byte, error) { return json.Marshal("") }, `""`)

	testCase("custom struct", struct {
		Name string `json:"name"`
	}{"Tit Petric"}, `{"response":{"name":"Tit Petric"}}`)

	SetConfig(Options{
		Pretty: true,
	})

	// Test pretty printing
	testCase("custom struct", struct {
		Name string `json:"name"`
	}{"Tit Petric"}, `{
	"response": {
		"name": "Tit Petric"
	}
}`)
	SetConfig(Options{
		Pretty: true,
		Trace:  true,
	})

	// Test pretty printing with tracing (no trace because stdlib error)
	testCase("custom struct", struct {
		Name string `json:"name"`
	}{"Tit Petric"}, `{
	"response": {
		"name": "Tit Petric"
	}
}`)

	// Test internal encoding error (this is really your own fault for passing chan's to json encode)
	{
		resp := testResponse(make(chan int))
		payload := struct {
			Error struct {
				Message string
				Trace   string
			}
		}{}
		err := json.Unmarshal([]byte(resp), &payload)
		if err != nil {
			t.Errorf("Unexpected error while decoding trace payload: %s", err)
		}
		if payload.Error.Message == "" {
			t.Errorf("Expected non-empty message on internal error")
		}
		if payload.Error.Trace == "" || payload.Error.Trace == E_EMPTY_TRACE {
			t.Logf("internal error: %+v\n", payload.Error)
			t.Errorf("Expected non-empty trace on internal error, got '%s'", payload.Error.Trace)
		}
	}

	// Test pretty printing with tracing (trace because pkg/errors)
	{
		resp := testResponse(e2.New("traced error"))
		payload := struct {
			Error struct {
				Message string
				Trace   string
			}
		}{}
		err := json.Unmarshal([]byte(resp), &payload)
		if err != nil {
			t.Errorf("Unexpected error while decoding trace payload: %s", err)
		}
		if payload.Error.Message != "traced error" {
			t.Errorf("Traced error decoding error, message doesn't match: %s != %s", "traced error", payload.Error.Message)
		}
		if payload.Error.Trace == "" || payload.Error.Trace == E_EMPTY_TRACE {
			t.Logf("traced error: %+v\n", payload.Error)
			t.Errorf("Expected non-empty trace on traced error, got '%s'", payload.Error.Trace)
		}
	}

	if getTrace() != E_EMPTY_TRACE {
		t.Errorf("expected empty trace from no errors given")
	}
}
