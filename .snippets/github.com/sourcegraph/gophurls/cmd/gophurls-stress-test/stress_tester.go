package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"sync"
	"text/template"
	"time"
)

var defaultCmdPath string

func init() {
	defaultCmdPath, _ = exec.LookPath("gophurls")
}

var cmdPath = flag.String("cmd", defaultCmdPath, "path to the gophurls program to test")
var numServers = flag.Int("servers", 1, "number of gophurls servers to spawn")
var numLinks = flag.Int("links", 10, "number of links to add")
var verbose = flag.Bool("v", false, "show verbose output")

type server struct {
	host string
	cmd  *exec.Cmd
}

var servers []*server

func main() {
	flag.Parse()

	if *cmdPath == "" {
		log.Fatal("Error: must specify -cmd. Run with -h for instructions.")
	}
	if *numServers < 1 {
		log.Fatal("Error: -servers must be at least 1.")
	}
	if *numLinks < 1 {
		log.Fatal("Error: -links must be a positive number.")
	}

	// Start servers.
	err := startServers()
	if err != nil {
		log.Fatalf("Error starting servers: %s", err)
	}
	defer killServers()

	// Start a fake server.
	var fakeServerRequests struct {
		n  int
		mu sync.Mutex
	}
	fakeMux := http.NewServeMux()
	fakeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		title := strings.TrimPrefix(r.URL.Path, "/")
		fmt.Fprintf(w, "<title>fetched-%s</title>", template.HTMLEscapeString(title))
		fakeServerRequests.mu.Lock()
		defer fakeServerRequests.mu.Unlock()
		fakeServerRequests.n++
	})
	fakeServer := httptest.NewServer(fakeMux)
	defer fakeServer.Close()

	t0 := time.Now()

	// Register servers as each other's peers.
	for _, s := range servers {
		// Make a list of all servers except for this one (to avoid self-loops).
		allPeers := make([]string, len(servers)-1)
		i := 0
		for _, s2 := range servers {
			if s.host == s2.host {
				// Don't set self as peer.
				continue
			}
			allPeers[i] = s2.host
			i++
		}
		allPeersJSON, err := json.Marshal(allPeers)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := http.Post(fmt.Sprintf("http://%s/peers", s.host), "application/json", bytes.NewReader(allPeersJSON))
		if err != nil {
			log.Fatalf("Error setting peers for %q: %s", s.host, err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Error setting peers for %q: HTTP status %d", s.host, resp.StatusCode)
		}
		resp.Body.Close()
		if *verbose {
			log.Printf("Set %d peers for %q.", len(allPeers), s.host)
		}
	}

	// Add links.
	if *verbose {
		log.Printf("Adding %d links...", *numLinks)
	}
	for i := 0; i < *numLinks; i++ {
		// Cycle through the servers.
		s := servers[i%len(servers)]

		// Add a link title for about half of all links.
		link := &link{
			URL: fmt.Sprintf("%s/page-%d", fakeServer.URL, i),
		}
		if i < *numLinks/2 {
			link.Title = fmt.Sprintf("unfetched-%d", i)
		}

		if *verbose {
			log.Printf("Adding link %v to %q...", link, s.host)
		}

		err := addLink(s.host, link)
		if err != nil {
			log.Printf("Error adding link %v: %s", link, err)
		}

		if *verbose {
			log.Printf("Added link %v.", link)
		}
	}

	if *verbose {
		log.Printf("Done.")
	}
	killServers()
	time.Sleep(time.Millisecond * 10)
	fmt.Printf("Total time elapsed: %s\n", time.Since(t0))
	fmt.Printf("# link fetches: %d\n", fakeServerRequests.n)
}

func startServers() error {
	servers = make([]*server, *numServers)
	for i := 0; i < *numServers; i++ {
		// Find an open port to listen on.
		l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4zero, Port: 0})
		if err != nil {
			return err
		}
		err = l.Close()
		if err != nil {
			return err
		}

		host := l.Addr().String()
		s := &server{
			host: host,
			cmd:  exec.Command(*cmdPath, fmt.Sprintf("-http=%s", host)),
		}
		s.cmd.Stdout, s.cmd.Stderr = os.Stdout, os.Stderr
		err = s.cmd.Start()
		if err != nil {
			return err
		}
		servers[i] = s
		if *verbose {
			log.Printf("Started server on %s.", s.host)
		}
	}
	time.Sleep(time.Millisecond * 50)
	return nil
}

func killServers() {
	for _, s := range servers {
		err := s.cmd.Process.Kill()
		if err != nil {
			log.Printf("Failed to kill server process for %q.", s.host)
		}
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
