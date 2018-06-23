package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
	DB     *mgo.Database
}

func (a *App) Initialize(server, database string) {
	session, err := mgo.Dial(server)
	if err != nil {
		log.Fatal(err)
	}
	a.DB = session.DB(database)

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/movies", a.getMovies).Methods("GET")
	a.Router.HandleFunc("/api/movies", a.createMovie).Methods("POST")
	a.Router.HandleFunc("/api/movies", a.updateMovie).Methods("PUT")
	a.Router.HandleFunc("/api/movies", a.deleteMovie).Methods("DELETE")
	a.Router.HandleFunc("/api/movies/{id}", a.getMovie).Methods("GET")
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":5000", a.Router))
}

func (a *App) getMovies(w http.ResponseWriter, r *http.Request) {
	m, err := getMovies(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, m)
}

func (a *App) getMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	m, err := getMovie(a.DB, params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Movie ID")
		return
	}
	respondWithJson(w, http.StatusOK, m)
}

func (a *App) createMovie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var m movie
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	m.ID = bson.NewObjectId()
	if err := m.createMovie(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, m)
}

func (a *App) updateMovie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var m movie
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := m.updateMovie(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) deleteMovie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var m movie
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := m.deleteMovie(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
