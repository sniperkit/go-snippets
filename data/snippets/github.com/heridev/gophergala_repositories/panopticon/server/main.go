package server

// Panopticon, Copyright © 2015, Huck Ridge Software LLC
// All rights reserved.

// vim:sw=4:ts=4

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/gophergala/panopticon/entry"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

const numEntries = 30

func init() {
	http.HandleFunc("/", Root)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/api/v1/add_entry", apiAddEntry)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	c.Infof("Logout")
	var u *user.User
	var url string
	var err error
	if u = user.Current(c); u == nil {
		url, err = user.LoginURL(c, r.URL.String())
	} else {
		url, err = user.LogoutURL(c, "/")
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
}

func userKey(token string, c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "User", token, 0, nil)
}

func Root(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	var u *user.User
	if u = user.Current(c); u == nil {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}

	// Get last numEntries entries & display them
	q := datastore.NewQuery("Entry").Ancestor(userKey(u.Email, c)).Order("-Time").Limit(numEntries)
	entries := make([]entry.Entry, 0, numEntries)
	if _, err := q.GetAll(c, &entries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for i := 0; i < len(entries)/2; i++ {
		entries[i], entries[len(entries)-i-1] = entries[len(entries)-i-1], entries[i]
	}
	for i, _ := range entries {
		entries[i].Time = entries[i].Time.Round(10 * time.Millisecond)
	}
	if err := homeTemplate.Execute(w, map[string]interface{}{
		"User":    u.Email,
		"Entries": entries}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var homeTemplate = template.Must(template.New("entry").Parse(`
<html>
  <head>
    <title>Gopher Gala Panopticon</title>
	<style>
	table, th, td {
		border: 1px solid black;
		border-collapse: collapse;
	}
	</style>
  </head>
  <body>
    <div>
		You are logged in as {{.User}}
	</div>
  	<table>
	  <tr>
	    <th>Time</th>
		<th>WasIdle</th>
		<th>Idle Time (ms)</th>
		<th>Window title</th>
	  </tr>
    {{range .Entries}}
	  <tr>
	  	<td>{{.Time}}</td>
	  	<td>{{.WasIdle}}</td>
	  	<td>{{.Idle}}</td>
	  	<td>{{.Title}}</td>
	  </tr>
    {{end}}
	</table>
	<a href="/logout">
		<button name="logout">Logout</button>
	</a>
	<div>
		Panopticon is a project for the 2015 Gopher Gala by
		<a href="https://twitter.com/readcodesing">Larry Clapp</a>
	</div>
  </body>
</html>
`))

func AddEntry(c appengine.Context, user string, e *entry.Entry) (*datastore.Key, error) {
	userKey := userKey(user, c)
	entryPartialKey := datastore.NewKey(c, "Entry", "", 0, userKey)
	if newEntryKey, err := datastore.Put(c, entryPartialKey, e); err != nil {
		return nil, errors.New("Could not put the entry")
	} else {
		return newEntryKey, nil
	}
}

func apiAddEntry(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	token := r.Header.Get("X-Panopticon-Token")
	if token == "" {
		http.Error(w, "Missing token header", http.StatusUnauthorized)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var entry entry.Entry
	err := decoder.Decode(&entry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// log.Printf("Adding %v", entry)
	newKey, err := AddEntry(c, token, &entry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Location", newKey.Encode())
	return
}
