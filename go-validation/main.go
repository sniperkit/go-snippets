package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
)

type articleRequestBody struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (a *articleRequestBody) validate() url.Values {
	errs := url.Values{}

	if a.Title == "" {
		errs.Add("title", "The title field is required!")
	}

	if len(a.Title) < 3 || len(a.Title) > 120 {
		errs.Add("title", "The title field must be between 3-120 chars!")
	}

	if a.Body == "" {
		errs.Add("body", "The body field is required!")
	}

	if len(a.Body) < 50 || len(a.Body) > 500 {
		errs.Add("body", "The title field must be between 50-500 chars!")
	}

	return errs
}

func handler(w http.ResponseWriter, r *http.Request) {
	article := &articleRequestBody{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(article); err != nil {
		panic(err)
	}

	if validErrs := article.validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.Header().Set("Content-type", "applciation/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("POST")
	fmt.Println("Listening on port: 5000")
	log.Fatal(http.ListenAndServe(":5000", r))
}
