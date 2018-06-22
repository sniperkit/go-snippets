package web

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"time"
	"goworkshop/persistence"
)

type RestServer struct {
	Port   int
	routes Routes
	router *mux.Router
	Store *persistence.DataStore
}

func (server *RestServer) StartServer() {
	server.initRoutes()
	server.router = mux.NewRouter()
	for _, route := range server.routes {
		server.router.Handle(route.Pattern, log(route.HandlerFunc)).Methods(route.Method)
	}

	fmt.Printf("Starting server on port: %d", server.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", server.Port), server.router); err != nil {
		panic(err)
	}
}

func log(routeFunc RouteFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request to: " + r.URL.Path)
		start := time.Now().UnixNano()
		err := routeFunc(w, r) // call original
		//handle the errors
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Errorf("error occurred while processing the request:%v", err)
		}
		end := time.Now().UnixNano()
		fmt.Printf("Request took: %d nano\n", end-start)
	})
}
