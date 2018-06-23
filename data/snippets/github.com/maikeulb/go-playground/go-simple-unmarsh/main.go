package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type stateRequestBody struct {
	Name string `json:"name"`
}

func (js stateRequestBody) State() State {
	fmt.Println("initialize state from request body")
	var s State
	s.Name = js.Name
	return s
}

type State struct {
	ID   int    `json:"id"`
	Name string ` json:"name"`
}

func (s *State) UnmarshalJSON(b []byte) (err error) {
	fmt.Println("unmarshalling state")
	var js stateRequestBody
	if err := json.Unmarshal(b, &js); err != nil {
		return err
	}
	*s = js.State()
	return
}

func stateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("state handler")

	var s []State
	encodedJSON :=
		[]byte(`[{
        "name": "mick"
     },{
        "name": "mick"
    }]`)

	encodedStream := bytes.NewReader(encodedJSON)

	if err := json.NewDecoder(encodedStream).Decode(&s); err != nil {
		panic(err)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", stateHandler).Methods("POST")
	fmt.Println("Listening on port: 5000")
	log.Fatal(http.ListenAndServe(":5000", r))
}
