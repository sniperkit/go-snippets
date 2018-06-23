package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

// TestAddLink_NoTitle_Broadcast tests that a newly added link with no title
// (and which therefore needs to be fetched) is broadcasted to peers after
// fetching completes.
func TestAddLink_NoTitle_Broadcast(t *testing.T) {
	// Start a test server that returns a page with a <title> tag, so we can
	// fetch locally.
	var fetched bool
	fakeTitleMux := http.NewServeMux()
	fakeTitleMux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`<title>Example</title>`))
		fetched = true
	})
	fakeTitleServer := httptest.NewServer(fakeTitleMux)
	defer fakeTitleServer.Close()

	link := `{"URL":"` + fakeTitleServer.URL + `"}`
	wantLink := `{"URL":"` + fakeTitleServer.URL + `","Title":"Example"}` + "\n"

	// Start a test server to receive the broadcasted link.
	var received bool
	fakePeerMux := http.NewServeMux()
	fakePeerMux.HandleFunc("/links", func(w http.ResponseWriter, r *http.Request) {
		received = true
		if body, _ := ioutil.ReadAll(r.Body); string(body) != wantLink {
			t.Errorf("got body %q, want %q", body, wantLink)
		}
	})
	fakePeerServer := httptest.NewServer(fakePeerMux)
	defer fakePeerServer.Close()

	// Add fake server to peers list.
	fakePeerURL, _ := url.Parse(fakePeerServer.URL)
	peers = map[string]struct{}{fakePeerURL.Host: struct{}{}}

	// Add the link to this server.
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/links", strings.NewReader(link))
	h.ServeHTTP(resp, req)
	testStatusCode(t, "after adding a link with no title (with peers to broadcast to)", resp.Code, http.StatusOK)

	// Wait until the fetch and broadcast have (probably) finished.
	time.Sleep(time.Millisecond * 10)

	// Test that adding the link to this server fetched it and broadcasted it to
	// its peer.
	if !fetched {
		t.Error("link was not fetched")
	}
	if !received {
		t.Error("fake server did not receive broadcasted link")
	}
}

// TestAddLink_WithTitle_Broadcast tests that a newly added link with a
// title (and which therefore does not need to be fetched) is immediately
// broadcasted to peers.
func TestAddLink_WithTitle_Broadcast(t *testing.T) {
	link := `{"URL":"http://example.com","Title":"Example"}` + "\n"

	// Start a test server to receive the broadcasted link.
	var received bool
	fakeMux := http.NewServeMux()
	fakeMux.HandleFunc("/links", func(w http.ResponseWriter, r *http.Request) {
		received = true
		if body, _ := ioutil.ReadAll(r.Body); string(body) != link {
			t.Errorf("got body %q, want %q", body, link)
		}
	})
	fakeServer := httptest.NewServer(fakeMux)
	defer fakeServer.Close()

	// Add fake server to peers list.
	fakeServerURL, _ := url.Parse(fakeServer.URL)
	peers = map[string]struct{}{fakeServerURL.Host: struct{}{}}

	// Add the link to this server.
	resp := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/links", strings.NewReader(link))
	h.ServeHTTP(resp, req)
	testStatusCode(t, "after adding a link with a title (with peers to broadcast to)", resp.Code, http.StatusOK)

	// Wait until the broadcast has (probably) finished.
	time.Sleep(10 * time.Millisecond)

	// Test that adding the link to this server broadcasted it to its peer.
	if !received {
		t.Error("fake server did not receive broadcasted link")
	}
}
