package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestAddPeer_OK tests that adding a peer succeeds, and that subsequently the
// newly added peer is present in the set of peers.
func TestAddPeer_OK(t *testing.T) {
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/peers", strings.NewReader(`["example.com:1234"]`))
	h.ServeHTTP(resp, req)

	testStatusCode(t, "after adding peers", resp.Code, http.StatusOK)
	if _, present := peers["example.com:1234"]; !present {
		t.Errorf(`got peers %v, want member "example.com:1234"`, peers)
	}
}
