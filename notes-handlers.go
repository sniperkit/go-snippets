package backend

import (
	"encoding/json"
	"google.golang.org/appengine"
	"net/http"
	"net/url"
)

func HandleNotesGet(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	notes := Notes{}

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	encodedURL := ""
	if val, exists := params["encodedurl"]; exists {
		encodedURL = val[0]
	} else {
		http.Error(w, "url parameter is missing in URI", http.StatusBadRequest)
		return
	}

	if err := notes.get(encodedURL, c); err != nil && err != ErrorNoMatch {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if err == ErrorNoMatch {
		http.Error(w, err.Error()+": "+encodedURL, http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func HandleNotesPut(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln(w, "hello qwerty")

	c := appengine.NewContext(r)
	notes := Notes{}

	if err := json.NewDecoder(r.Body).Decode(&notes); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := notes.put(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func HandleNotesDelete(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	notes := Notes{}

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if val, exists := params["encodedurl"]; exists {
		notes.URL = val[0]
	} else {
		http.Error(w, "url parameter is missing in URI", http.StatusBadRequest)
		return
	}

	err = notes.delete(c)
	if err == ErrorNoMatch {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func HandleNotesGetAll(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	notes := NotesAll{}

	if err := notes.get(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// TODO:
//// get/post api endpoint: if url (passed in request body) is already available in database, return that record, otherwise create a record for that url and return that record
//// put api endpoint: update/overwrite the record in the database for the given url
//// struct fields: url, notes (plain text for now, but will be html soon)

// Challenges:
