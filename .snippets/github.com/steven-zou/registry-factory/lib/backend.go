package lib

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

//BackendServer ...
type BackendServer struct {
	server  *http.Server
	mux     *http.ServeMux
	host    string
	port    uint16
	running bool
}

//NewBackendServer ...
func NewBackendServer(host string, port uint16) *BackendServer {
	return &BackendServer{
		host: host,
		port: port,
		mux:  http.NewServeMux(),
	}
}

//Start the server
func (bs *BackendServer) Start() {
	if bs.running {
		return
	}

	if bs.server == nil {
		bs.server = &http.Server{
			Addr: fmt.Sprintf("%s:%d", bs.host, bs.port),
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				bs.mux.ServeHTTP(w, r)
			}),
		}
	}

	bs.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("Echo from backend: %d", time.Now().Unix())))
	})

	bs.server.ListenAndServe()
}

//Stop ...
func (bs *BackendServer) Stop(ctx context.Context) error {
	if !bs.running {
		return nil
	}

	if bs.server == nil {
		return errors.New("nil server")
	}

	return bs.server.Shutdown(ctx)
}
