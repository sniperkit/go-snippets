package webserver

import (
	"fmt"
	"net/http"

	"github.com/gophergala/gallivanting_gophers/data"
	"github.com/julienschmidt/httprouter"
)

// Server encapsulates the webserver aspect of the service. This will handle the
// routing and registration of routes (as well as security when implemented).
type Server struct {
	router *httprouter.Router
	sm     *SessionManager
}

// NewServer creates a webserver and sets up the router.
func NewServer(db *data.DB) *Server {
	s := &Server{
		router: httprouter.New(),
		sm:     NewSessionManager(db),
	}

	return s
}

// Start turns on the webserver and begins accepting connections.
func (s *Server) Start() {
	fmt.Println("Listening on :8888")
	go http.ListenAndServe(":8888", corsHandler(s.router))
}

// RegisterGet registers a get request handler with the webserver.
func (s *Server) RegisterGet(url string, fn func(w http.ResponseWriter, r *http.Request, p httprouter.Params)) {
	s.router.GET(url, fn)
}

// RegisterPost registers a post request handler with the webserver.
func (s *Server) RegisterPost(url string, fn func(w http.ResponseWriter, r *http.Request, p httprouter.Params)) {
	s.router.POST(url, fn)
}

// Used to enable cors support
func corsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
