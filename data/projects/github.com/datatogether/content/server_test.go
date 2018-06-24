package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	cfg = &config{}
}

func TestServerRoutes(t *testing.T) {
	cases := []struct {
		method, endpoint string
		body             []byte
		resStatus        int
	}{
		{"GET", "/", nil, 200},
	}

	client := &http.Client{}

	// s, err := New(func(opt *Config) {
	// 	opt.Online = false
	// 	opt.MemOnly = true
	// })
	// if err != nil {
	// 	t.Error(err.Error())
	// 	return
	// }

	server := httptest.NewServer(NewServerRoutes())

	for i, c := range cases {
		req, err := http.NewRequest(c.method, server.URL+c.endpoint, bytes.NewReader(c.body))
		if err != nil {
			t.Errorf("case %d error creating request: %s", i, err.Error())
			continue
		}

		res, err := client.Do(req)
		if err != nil {
			t.Errorf("case %d error performing request: %s", i, err.Error())
			continue
		}

		if res.StatusCode != c.resStatus {
			t.Errorf("case %d: %s - %s status code mismatch. expected: %d, got: %d", i, c.method, c.endpoint, c.resStatus, res.StatusCode)
			continue
		}
	}
}
