package router

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func GetRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", BaseHandler(Index))
	r.HandleFunc("/dashboard", BaseHandler(Dashboard))

	// Registrations routes

	r.HandleFunc("/registrations", BaseHandler(RegisterSession)).Methods("POST")

	// Sessions routes

	r.HandleFunc("/sessions", BaseHandler(CreateSession)).Methods("POST")
	r.HandleFunc("/sessions", BaseHandler(DestroySession)).Methods("DELETE")

	// Subscriptions routes

	r.HandleFunc("/subscriptions", BaseHandler(SubscriptionsIndex)).Methods("GET")
	r.HandleFunc("/subscriptions", BaseHandler(SubscriptionsCreate)).Methods("POST")
	r.HandleFunc("/subscriptions/{id:[0-9]+}", BaseHandler(SubscriptionsDestroy)).Methods("DELETE")

	// Serve static assets

	r.Handle("/public/javascripts/{rest}", http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
	r.Handle("/public/stylesheets/{rest}", http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
	r.Handle("/public/images/{rest}", http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))

	return r
}

func BaseHandler(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		log.Println("Started " + r.Method + " " + r.URL.Path + " from " + r.RemoteAddr)
		fn(w, r)
	}
}
