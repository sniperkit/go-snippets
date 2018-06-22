package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/GopherGala/i_love_indexes/elasticsearch"
	"github.com/Scalingo/go-workers"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"gopkg.in/errgo.v1"
)

func NewAPI() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/api/indices", AddIndexOf).Methods("POST")
	router.HandleFunc("/api/indices", ListIndexOf).Methods("GET")
	router.HandleFunc("/api/search", SearchIndexItems).Methods("GET")
	stack := negroni.New(negroni.NewLogger(), negroni.NewRecovery())
	stack.UseHandler(router)
	return stack
}

type BadRequestError struct {
	err string `json:"error"`
}

type UnprocessableEntityError struct {
	Errors map[string][]string `json:"errors"`
}

func NewBadRequestError(err string) *BadRequestError {
	return &BadRequestError{err}
}

func NewUnprocessableEntityError() *UnprocessableEntityError {
	err := &UnprocessableEntityError{}
	err.Errors = make(map[string][]string)
	return err
}

func (err *UnprocessableEntityError) Add(field string, errorMessage string) {
	err.Errors[field] = append(err.Errors[field], errorMessage)
}

type AddIndexOfParams struct {
	URL string `json:"url"`
}

func ListIndexOf(res http.ResponseWriter, req *http.Request) {
	indices, err := elasticsearch.ListIndexOf()
	if err != nil {
		res.WriteHeader(500)
		log.Println(errgo.Details(err))
		fmt.Fprintf(res, errgo.Details(err))
		return
	}

	res.WriteHeader(200)
	json.NewEncoder(res).Encode(&indices)
}

func AddIndexOf(res http.ResponseWriter, req *http.Request) {
	var params *AddIndexOfParams
	err := json.NewDecoder(req.Body).Decode(&params)
	if err != nil {
		res.WriteHeader(400)
		fmt.Fprintf(res, "invalid JSON: %v", err)
		return
	}

	if params.URL == "" {
		res.WriteHeader(422)
		err := NewUnprocessableEntityError()
		err.Add("url", "can't be blank")
		json.NewEncoder(res).Encode(&err)
		return
	}

	u, err := url.Parse(params.URL)
	if err != nil {
		res.WriteHeader(422)
		err := NewUnprocessableEntityError()
		err.Add("url", "is not an URL")
		json.NewEncoder(res).Encode(&err)
		return
	}

	if u.Path == "/" {
		u.Path = ""
	}
	index := &elasticsearch.IndexOf{Host: u.Host, Scheme: u.Scheme, Path: u.Path}
	err = index.Index()
	if err != nil {
		if err == elasticsearch.AlreadyIndexedErr {
			res.WriteHeader(422)
			err := NewUnprocessableEntityError()
			err.Add("url", "already taken")
			json.NewEncoder(res).Encode(&err)
			return
		}
		res.WriteHeader(500)
		log.Println(errgo.Details(err))
		fmt.Fprintf(res, errgo.Details(err))
		return
	}

	workers.Enqueue("index-crawler", "CrawlWorker", []string{index.Id, ""})
	res.WriteHeader(201)
	json.NewEncoder(res).Encode(&index)
}

func SearchIndexItems(res http.ResponseWriter, req *http.Request) {
	from := req.URL.Query().Get("from")
	query := req.URL.Query().Get("q")
	typ := req.URL.Query().Get("t")
	if typ == "" {
		typ = "any"
	}

	if query == "" {
		res.WriteHeader(400)
		err := NewBadRequestError("requires 'query' params")
		json.NewEncoder(res).Encode(&err)
		return
	}

	res.WriteHeader(200)
	indexItems := elasticsearch.SearchIndexItemsPerName(from, typ, query)
	json.NewEncoder(res).Encode(&indexItems)
}
