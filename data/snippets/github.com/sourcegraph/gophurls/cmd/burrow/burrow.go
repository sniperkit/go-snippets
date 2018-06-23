package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"text/template"
	"time"
)

var httpAddr = flag.String("http", "localhost:7002", "externally addressable host:port to host stats counter server")
var peerAddr = flag.String("peer", "localhost:7001", "externally addressable peer listen host:port to use (when adding self as peer to servers)")
var serversStr = flag.String("servers", "", "comma-separated list of servers (ex: 'example.com:7000,foo.com:1234')")
var numLinks = flag.Int("links", 10, "number of links to add per server")
var verbose = flag.Bool("v", false, "show verbose output")

var servers []string

func main() {
	flag.Parse()

	if *serversStr == "" {
		log.Fatal("Error: -servers must not be empty.")
	}
	if *numLinks < 1 {
		log.Fatal("Error: -links must be a positive number.")
	}

	// Start a fake server.
	stats := struct {
		fetched map[string]int
		readded map[string]int
		mu      sync.Mutex
	}{fetched: make(map[string]int), readded: make(map[string]int)}
	fakeServerMux := http.NewServeMux()
	fakeServerMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(120)))
		id := strings.TrimPrefix(r.URL.Path, "/")
		stats.mu.Lock()
		defer stats.mu.Unlock()
		stats.fetched[id]++
		fmt.Fprintf(w, "<title>fetched-%s</title>", template.HTMLEscapeString(id))
	})
	go func() {
		err := http.ListenAndServe(*httpAddr, fakeServerMux)
		if err != nil {
			log.Fatal(err)
		}
	}()
	fakePeerMux := http.NewServeMux()
	fakePeerMux.HandleFunc("/links", func(w http.ResponseWriter, r *http.Request) {
		var link *link
		err := json.NewDecoder(r.Body).Decode(&link)
		if err != nil {
			http.Error(w, fmt.Sprintf("bad JSON: %s", err), http.StatusBadRequest)
			return
		}
		// Validate the URL.
		if link.URL == "" {
			http.Error(w, "no url", http.StatusBadRequest)
			return
		}
		url, err := url.Parse(link.URL)
		if err != nil {
			http.Error(w, "bad url", http.StatusBadRequest)
			return
		}
		stats.mu.Lock()
		defer stats.mu.Unlock()
		stats.readded[strings.TrimPrefix(url.Path, "/")]++
	})
	fakePeerMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	go func() {
		err := http.ListenAndServe(*peerAddr, fakePeerMux)
		if err != nil {
			log.Fatal(err)
		}
	}()

	servers = strings.Split(*serversStr, ",")

	// Register servers as each other's peers.
	for _, s := range servers {
		// Make a list of all servers except for this one (to avoid self-loops).
		allPeers := make([]string, len(servers))
		for i, s2 := range servers {
			if s == s2 {
				// Don't add a server as its own peer. Instead, use this slot
				// for the burrow server.
				allPeers[i] = *peerAddr
				continue
			}
			allPeers[i] = s2
		}
		allPeersJSON, err := json.Marshal(allPeers)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := http.Post(fmt.Sprintf("http://%s/peers", s), "application/json", bytes.NewReader(allPeersJSON))
		if err != nil {
			log.Printf("Error setting peers for %q: %s", s, err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("Error setting peers for %q: HTTP status %d", s, resp.StatusCode)
		}
		if resp.Body != nil {
			resp.Body.Close()
		}
		if *verbose {
			log.Printf("Set %d peers for %q.", len(allPeers), s)
		}
	}

	// Add links.
	if *verbose {
		log.Printf("Adding %d links to each server...", *numLinks)
	}
	for i := 0; i < *numLinks; i++ {
		if i%2 == 0 {
			sort.Strings(servers)
		} else {
			sort.Sort(sort.Reverse(sort.StringSlice(servers)))
		}
		for _, s := range servers {
			// Add a link title for about half of all links.
			link := &link{
				URL: fmt.Sprintf("http://%s/%s?i=%d", *httpAddr, strings.Replace(s, ":", "-", -1), i),
			}
			if i < *numLinks/2 {
				link.Title = fmt.Sprintf("server %s link %d", s, i)
			}

			if *verbose {
				log.Printf("Adding link %v to %q...", link, s)
			}

			go func() {
				err := addLink(s, link)
				if err != nil {
					log.Printf("Error adding link %v: %s", link, err)
				}

				if *verbose {
					log.Printf("Added link %v.", link)
				}
			}()
		}
	}

	if *verbose {
		log.Printf("Done.")
	}
	time.Sleep(time.Millisecond * 5000)
	stats.mu.Lock()
	defer stats.mu.Unlock()

	fmt.Println("Fetches:")
	for s, n := range stats.fetched {
		fmt.Printf("%s\t%d\n", s, n)
	}
	fmt.Println()

	fmt.Println("Re-adds:")
	for s, n := range stats.readded {
		fmt.Printf("%s\t%d\n", s, n)
	}
}

type link struct {
	URL   string
	Title string `json:",omitempty"`
}

func addLink(host string, link *link) error {
	linkJSON, err := json.Marshal(link)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/links", host), "application/json", bytes.NewReader(linkJSON))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP status %d", resp.StatusCode)
	}
	return nil
}
