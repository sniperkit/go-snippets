package backend

import (
	"github.com/gorilla/mux"
	"net/http"
)

func init() {
	r := mux.NewRouter()

	r.Handle("/notes", mw.ThenFunc(HandleNotesGet)).Methods("GET")
	r.Handle("/notes", mw.ThenFunc(HandleNotesPut)).Methods("PUT")
	r.Handle("/notes", mw.ThenFunc(HandleNotesDelete)).Methods("DELETE")
	r.Handle("/notesall", mw.ThenFunc(HandleNotesGetAll)).Methods("GET") //TODO: /notesall is not looking like correct name. may be just /notes with a special query parameter like /notes?all=true

	http.Handle("/", r)
}
