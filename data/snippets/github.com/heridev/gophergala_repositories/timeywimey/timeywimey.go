package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	time.Local = time.UTC

	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/", DefaultHandler)
	r.HandleFunc("/auth", AuthenticationHandler)

	user := r.Path("/{username}").Subrouter()
	user.Methods("GET").HandlerFunc(UserIndexHandler)

	meetings := r.PathPrefix("/{username}/{project}/meetings").Subrouter()
	meetings.Methods("GET").Path("/calendar").HandlerFunc(MeetingCalendarHandler)
	meetings.Methods("GET").HandlerFunc(MeetingShowHandler)
	meetings.Methods("PUT", "POST").HandlerFunc(MeetingUpdateHandler)
	meetings.Methods("DELETE").HandlerFunc(MeetingDeleteHandler)

	issues := r.PathPrefix("/{username}/{project}/issues").Subrouter()
	issues.Methods("GET").Path("/calendar").HandlerFunc(IssueCalendarHandler)
	issues.Methods("GET").HandlerFunc(IssueShowHandler)
	issues.Methods("PUT", "POST").HandlerFunc(IssueUpdateHandler)
	issues.Methods("DELETE").HandlerFunc(IssueDeleteHandler)

	time := r.PathPrefix("/{username}/{project}/timesheet").Subrouter()
	time.Methods("GET").Path("/calendar").HandlerFunc(TimeCalendarHandler)
	time.Methods("GET").HandlerFunc(TimeShowHandler)
	time.Methods("PUT", "POST").HandlerFunc(TimeUpdateHandler)
	time.Methods("DELETE").HandlerFunc(TimeDeleteHandler)

	invoice := r.PathPrefix("/{username}/{project}/invoices").Subrouter()
	invoice.Methods("GET").Path("/calendar").HandlerFunc(InvoiceCalendarHandler)
	invoice.Methods("GET").HandlerFunc(InvoiceShowHandler)
	invoice.Methods("PUT", "POST").HandlerFunc(InvoiceUpdateHandler)
	invoice.Methods("DELETE").HandlerFunc(InvoiceDeleteHandler)

	project := r.PathPrefix("/{username}/{project}").Subrouter()
	project.Methods("GET").Path("/new").HandlerFunc(ProjectNewHandler)
	project.Methods("GET").HandlerFunc(ProjectIndexHandler)
	project.Methods("PUT", "POST").HandlerFunc(ProjectUpdateHandler)

	log.Fatal(http.ListenAndServe(address(), r))
}

// Retrieve the web server address from the environment variable TW_SERVER.
// If the environment variable is not set then default to `localhost:8080`.
func address() string {
	env := os.Getenv("TW_SERVER")
	if env == "" {
		return "localhost:8080"
	}
	return env
}
