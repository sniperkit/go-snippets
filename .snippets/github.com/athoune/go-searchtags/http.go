package main

import (
	/*"encoding/json"*/
	"fmt"
	"github.com/bmizerany/pat"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func toInt(s []string) []uint64 {
	ints := make([]uint64, len(s))
	for i, value := range s {
		var parsed int64
		parsed, _ = strconv.ParseInt(value, 10, 64) //[FIXME] handle error
		ints[i] = uint64(parsed)
	}
	return ints
}

func doSearch(w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()["q"]
	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Connection", "keep-alive")
	w.WriteHeader(200)
	w.Write([]byte("["))
	answer := make(chan []*docScore)
	searchQueue <- queryAnswer{
		NewDocument(toInt(strings.Split(q[0], ","))),
		answer}
	docScores := <-answer
	for _, ds := range docScores {
		w.Write([]byte(fmt.Sprintf("%d,", ds.doc)))
	}
	w.Write([]byte("]"))
}

func doSimilar(w http.ResponseWriter, req *http.Request) {
	id64, _ := strconv.ParseInt(req.URL.Query().Get("name"), 10, 32)
	id := uint32(id64)
	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Connection", "keep-alive")
	w.WriteHeader(200)
	w.Write([]byte("["))
	answer := make(chan []*docScore)
	searchQueue <- queryAnswer{
		docs.tags[id],
		answer}
	docScores := <-answer
	for _, ds := range docScores {
		if ds.doc != id {
			w.Write([]byte(fmt.Sprintf("%d,", ds.doc)))
		}
	}
	w.Write([]byte("]"))
}

func startHttp() {
	m := pat.New()
	m.Get("/search", http.HandlerFunc(doSearch))
	m.Get("/similar/:doc", http.HandlerFunc(doSimilar))
	http.Handle("/", m)
	log.Printf("About to start http://localhost:8000")
	err := http.ListenAndServe("localhost:8000", nil)
	if err != nil {
		panic(err)
	}
}
