package api

import (
	"github.com/gorilla/mux"
	"github.com/vsouza/pug/config"
)

// Router is an exported struct to store router, config and resources of
// service, to be used later.
type Router struct {
	router    *mux.Router
	config    *config.Config
	resources *Resources
}

// NewRouter is an exported method to start an mux router an mount APIRouter
// struct, and returns it.
func NewRouter(c *config.Config, r *Resources) *Router {
	return &Router{
		router:    mux.NewRouter(),
		config:    c,
		resources: r,
	}
}

// GetRouter is an exported method to return router pointer.
func (ar *Router) GetRouter() *mux.Router {
	return ar.router
}

// InitRoutes is an exported method to start all handlers with mux.router pointer
func (ar *Router) InitRoutes() {
	ar.router.HandleFunc("/health", ar.healthHandler).Methods("GET")
}
