package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// h is an alias (to save typing) to the default router with the handler we
// registered in server.go's init function.
var h = http.DefaultServeMux

// TestGetHome tests that the homepage returns successfully and includes the
// word "Links" somewhere.
func TestGetHome(t *testing.T) {
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	h.ServeHTTP(resp, req)

	testStatusCode(t, "homepage", resp.Code, http.StatusOK)
	if body := resp.Body.String(); !strings.Contains(body, "Links") {
		t.Errorf(`want the word "Links" to appear somewhere on homepage, got %q`, body)
	}
}

// TestAddLink_OK tests that adding a link succeeds, and that subsequently the
// homepage contains the newly added link.
func TestAddLink_OK(t *testing.T) {
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/links", strings.NewReader(`{"URL":"http://example.com"}`))
	h.ServeHTTP(resp, req)

	testStatusCode(t, "after adding a link", resp.Code, http.StatusOK)

	// Test that the link now appears in the list.
	resp = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/", nil)
	h.ServeHTTP(resp, req)
	if body := resp.Body.String(); !strings.Contains(body, "example.com") {
		t.Errorf(`want "example.com" to appear somewhere on homepage after adding a link, got %q`, body)
	}
}

func testStatusCode(t *testing.T, label string, got, want int) {
	if got != want {
		t.Errorf("%s: got HTTP %d, want HTTP %d", label, got, want)
	}
}
