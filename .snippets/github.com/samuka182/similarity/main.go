package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	similarity "similarity/search"
	"sort"

	"github.com/gorilla/mux"
)

var level = map[string]float64{
	"EXTRA_LOW":  0.0,
	"LOW":        0.2,
	"MEDIUM":     0.4,
	"HIGH":       0.6,
	"EXTRA_HIGH": 0.8,
}

// RequestPayload its a struct of request payload
type RequestPayload struct {
	Dictionary []string `json:"dictionary"`
	Input      string   `json:"input"`
	Level      string   `json:"level"`
}

// ResponsePayload its a struct of response payload
type ResponsePayload struct {
	ResultCode string      `json:"resultCode"`
	ResultData interface{} `json:"resultData"`
	Errors     []Error     `json:"errors"`
}

// Error ...
type Error struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// ResultsPayload its a struct of results data
type ResultsPayload struct {
	Results []string `json:"results"`
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.
		Methods("POST").
		Path("/similarity").
		HandlerFunc(post)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func post(w http.ResponseWriter, r *http.Request) {
	response := ResponsePayload{}
	jsonRequest, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		response.ResultCode = "ERROR"
		erro := Error{}
		erro.Code = "error"
		erro.Description = err.Error()
		response.Errors = append(response.Errors, erro)

		dispatch(w, response)
	} else {
		payload := RequestPayload{}
		err = json.Unmarshal(jsonRequest, &payload)

		retorno := similarity.Exec(payload.Dictionary, payload.Input, level[payload.Level])

		response.ResultCode = "SUCCESS"
		response.ResultData = getResults(retorno)

		dispatch(w, response)
	}
}

func dispatch(w http.ResponseWriter, resp ResponsePayload) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	err := enc.Encode(resp)
	if err != nil {
		panic(err)
	}
}

func getResults(matches []similarity.Match) ResultsPayload {
	results := ResultsPayload{}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Similarity > matches[j].Similarity
	})

	for _, m := range matches {
		results.Results = appendIfMissing(results.Results, m)
	}

	return results
}

func appendIfMissing(slice []string, i similarity.Match) []string {
	for _, ele := range slice {
		if ele == i.Value {
			return slice
		}
	}
	return append(slice, i.Value)
}
