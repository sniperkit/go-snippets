package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"time"
)

type stateRequestBody struct {
	Name string `json:"name"`
}

func (js stateRequestBody) State() State {
	var s State
	s.Name = strings.ToUpper(js.Name)
	return s
}

func (js *stateRequestBody) validate() error {
	if js.Name <= "" {
		return errors.New("Name should not be empty")
	}

	return nil
}

type State struct {
	ID   int
	Name string
}

func (s *State) UnmarshalJSON(data []byte) error {
	var js stateRequestBody

	if err := json.Unmarshal(data, &js); err != nil {
		return err // panic?
	}
	if err := js.validate(); err != nil {
		return err
	}

	*s = js.State()
	return nil
}

type parkRequestBody struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	NearestCity string    `json:"nearestCity"`
	Visitors    int       `json:"visitors"`
	Established time.Time `json:"established"`
	StateID     int       `json:"stateId"`
}

func (jp *parkRequestBody) validate() error {
	fmt.Println("validating park")
	if jp.Name <= "" {
		return errors.New("Name should not be empty")
	}
	if jp.Description <= "" {
		return errors.New("Description should not be empty")
	}

	return nil
}

func (jp parkRequestBody) Park() Park {
	fmt.Println("mapping park json to struct")
	var p Park
	p.Name = strings.ToUpper(jp.Name)
	p.Description = jp.Description
	return p
}

type Park struct {
	ID          int
	Name        string
	Description string
	NearestCity string
	Visitors    int
	Established time.Time
	StateID     int
	State       State
}

func (p *Park) UnmarshalJSON(data []byte) error {
	var jp parkRequestBody

	if err := json.Unmarshal(data, &jp); err != nil {
		return err // panic?
	}
	if err := jp.validate(); err != nil {
		return err
	}

	*p = jp.Park()
	return nil
}

func parkHandler(w http.ResponseWriter, r *http.Request) {
	p := &Park{}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(p); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, p)
	}
}

func stateHandler(w http.ResponseWriter, r *http.Request) {
	s := &State{} // initialize state object
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(s); err != nil { // decode body to state object
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, s)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/parks", parkHandler).Methods("POST")
	r.HandleFunc("/api/states", stateHandler).Methods("POST")
	fmt.Println("Listening on port: 5000")
	log.Fatal(http.ListenAndServe(":5000", r))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
