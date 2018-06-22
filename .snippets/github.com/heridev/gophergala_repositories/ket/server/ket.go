package server

import (
	"crypto/tls"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Server struct {
	Config *LConfig
	CA     *CertAuthority
	pac    string
	//CertFile, KeyFile string
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.handleInternal(w, r) {
		return
	}
	if s.isBlocked(r.URL) {
		fmt.Fprintf(w, "Blocked request: %q", html.EscapeString(r.URL.Path))
		return
	}
	// Proxy, otherwise.
	s.proxy(w, r)
}

func (s *Server) Init() error {
	pac, err := ioutil.ReadFile("./data/proxy.pac")
	if err != nil {
		return err
	}
	s.pac = string(pac)
	return nil
}

func (s *Server) Start(port string) error {
	log.Printf("|ket%s〉\n", port)
	server := &http.Server{
		Addr: port, Handler: s,
	}
	return server.ListenAndServe()
}

func (s *Server) StartTLS(port, certFile, keyFile string) error {
	log.Printf("|ket%s〉\n", port)
	server := &http.Server{
		Addr: port, Handler: s,
		/*TLSConfig: &tls.Config{
			GetCertificate:     s.CA.Get,
			InsecureSkipVerify: true,
		},*/
	}
	return server.ListenAndServeTLS(s.CA.CertFile, s.CA.KeyFile)
}

//------------------------------------------------------------------------------

func isInternalHost(host string) bool {
	return host == "" || host == "ket" || host == "ket:443"
}

func (s *Server) handleInternal(w http.ResponseWriter, r *http.Request) bool {
	url := r.URL
	if !isInternalHost(url.Host) {
		return false
	}
	if url.Path == "/proxy.pac" {
		log.Println("PAC requested!")
		if url.RawQuery != "txt" {
			w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig")
		}
		fmt.Fprint(w, s.pac)
		return true
	}

	// Try to serve content from local file system.
	data := s.Config.Data
	for _, dir := range data.Dirs {
		prefix := dir.Url
		if !strings.HasPrefix(url.Path, prefix) {
			continue
		}
		// TODO: fix annoying bug with missing '/' at the end of path!
		//return http.StripPrefix(prefix, http.FileServer(http.Dir(dir.FPath)))
		url.Path = strings.TrimPrefix(url.Path, prefix)
		handler := http.FileServer(http.Dir(dir.FPath))
		handler.ServeHTTP(w, r)
		return true
	}

	// TODO: make redirect to index url.
	fmt.Fprintf(w, "Hello, %q, %q", html.EscapeString(url.Path), url.RawQuery)
	return true
}

func (s *Server) isBlocked(url *url.URL) bool {
	data := s.Config.Data
	for _, block := range data.Block {
		// TODO: implement full url match!
		if strings.HasPrefix(url.Path, block) {
			return true
		}
	}
	return false
}

func (s *Server) proxy(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v\n", r.URL)

	// Based on: https://golang.org/src/net/http/httputil/reverseproxy.go
	out := new(http.Request)
	*out = *r

	out.Proto = "HTTP/1.1"
	out.ProtoMajor = 1
	out.ProtoMinor = 1
	out.Close = false

	transport := http.DefaultTransport
	if r.TLS != nil {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		out.URL.Scheme = "https"
	} else {
		transport = &http.Transport{}
		out.URL.Scheme = "http"
	}

	res, err := transport.RoundTrip(out)
	if err != nil {
		log.Println("Proxy error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	copyHeader(w.Header(), res.Header)
	w.WriteHeader(res.StatusCode)
	_, err = io.Copy(w, res.Body)
	if err != nil {
		log.Println("io.Copy:", err)
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
