package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kkamdooong/go-restful-api-example/handler"
)

func main() {

	var (
		port = flag.String("port", "8080", "The WEB Server port number.")
	)
	flag.Parse()

	r := mux.NewRouter().StrictSlash(true)

	sub := r.PathPrefix("/api/v1").Subrouter()
	sub.Methods("GET").Path("/companies").HandlerFunc(handler.GetCompanies)
	sub.Methods("POST").Path("/companies").HandlerFunc(handler.SaveCompany)
	sub.Methods("GET").Path("/companies/{name}").HandlerFunc(handler.GetCompany)
	sub.Methods("PUT").Path("/companies/{name}").HandlerFunc(handler.UpdateCompany)
	sub.Methods("DELETE").Path("/companies/{name}").HandlerFunc(handler.DeleteCompany)

	// On the default page we will simply serve our static index page.
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./views/")))
	//	r.Handle("/", http.FileServer(http.Dir("./views/")))
	// We will setup our server so we can serve static assets like images, css from the /static/{file} route
	//	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Host", "User-Agent", "Accept", "Accept-Language", "Accept-Encoding", "Access-Control-Request-Method", "Access-Control-Request-Headers", "Origin", "Connection"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	fmt.Println("Staring WEB Server on port " + *port)

	log.Fatal(http.ListenAndServe(":"+*port, handlers.CORS(originsOk, headersOk, methodsOk)(r)))

}
