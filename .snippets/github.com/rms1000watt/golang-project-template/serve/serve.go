package serve

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/rms1000watt/belly"
	"github.com/rms1000watt/golang-project-template/handlers"
	log "github.com/sirupsen/logrus"
)

// CORSOrigins is the whitelisted list of CORS origins for middleware
var CORSOrigins = []string{
	"https://localhost",
	"https://127.0.0.1",
	"https://www.example.com",
}

// Serve starts the server
func Serve(cfg Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	srv := http.Server{
		Addr:              addr,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		Handler:           ServerHandler(),
	}

	log.Info("Starting HTTPS server at: ", addr)
	log.Fatal(srv.ListenAndServeTLS(filepath.Join(cfg.CertsPath, cfg.CertName), filepath.Join(cfg.CertsPath, cfg.KeyName)))
}

// ServerHandler maps all the paths to handlers via mux
func ServerHandler() http.Handler {
	belly.SetCORSOrigins(CORSOrigins)

	mux := http.NewServeMux()
	mux.HandleFunc("/person", belly.HandleMiddlewares(handlers.PersonHandler, belly.MiddlewareNoCache, belly.MiddlewareCORS, belly.MiddlewareLogger, belly.MiddlewareSecure))
	mux.HandleFunc("/", belly.HandleMiddlewares(handlers.RootHandler, belly.MiddlewareLogger))

	return mux
}
