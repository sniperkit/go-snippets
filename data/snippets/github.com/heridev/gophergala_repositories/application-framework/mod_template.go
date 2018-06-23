// this is a template module. do nothing. use to start your application.
package main

import (
	"html/template"
	"net/http"
)

func templateModule(w http.ResponseWriter, r *http.Request) {
	//this must add at begin of every session code
	c, err := r.Cookie("session")
	if err != nil || c.Value == "" {
		http.Error(w, "Session expired", 401)
		return
	}

	//build page content
	b := `<pre>This is a template module.`

	//finally show the page
	p := Page{
		Title:  "template",
		Status: c.Value,
		Body:   template.HTML(b),
	}
	t.ExecuteTemplate(w, "index.html", p)
}

func init() {
	http.HandleFunc("/templateModule", templateModule)
}
