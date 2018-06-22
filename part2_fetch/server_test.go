package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

// TestAddLink_NoTitle tests that adding a link (with no title specified)
// succeeds, that the title is fetched, and that subsequently the
// homepage contains the title of the newly added link.
func TestAddLink_NoTitle(t *testing.T) {
	// Start a test server that returns a page with a <title> tag, so we can
	// fetch locally.
	fakeMux := http.NewServeMux()
	fakeMux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`<title>Example</title>`))
	})
	fakeServer := httptest.NewServer(fakeMux)
	defer fakeServer.Close()

	// Add the URL of the fake server.
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/links", strings.NewReader(`{"URL":"`+fakeServer.URL+`"}`))
	h.ServeHTTP(resp, req)

	testStatusCode(t, "after adding a link without a title", resp.Code, http.StatusOK)

	// Wait until the fetch has (probably) finished.
	time.Sleep(time.Millisecond * 10)

	// Test that the link now appears in the list.
	resp = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/", nil)
	h.ServeHTTP(resp, req)
	if body := resp.Body.String(); !strings.Contains(body, "Example") {
		t.Errorf(`want "Example" to appear somewhere on homepage after adding a link, got %q`, body)
	}
}

// TestAddLink_WithTitle tests that adding a link (with a title) succeeds, and
// that subsequently the homepage contains the title of the newly added link.
func TestAddLink_WithTitle(t *testing.T) {
	// Add a link with the title field set (which means we don't need to fetch
	// the link to determine the title).
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/links", strings.NewReader(`{"URL":"http://example.com","Title":"Example"}`))
	h.ServeHTTP(resp, req)

	testStatusCode(t, "after adding a link with a title", resp.Code, http.StatusOK)

	// Test that the link now appears in the list.
	resp = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/", nil)
	h.ServeHTTP(resp, req)
	if body := resp.Body.String(); !strings.Contains(body, "Example") {
		t.Errorf(`want "Example" in response body, got %q`, body)
	}
}

func testStatusCode(t *testing.T, label string, got, want int) {
	if got != want {
		t.Errorf("%s: got HTTP %d, want HTTP %d", label, got, want)
	}
}
