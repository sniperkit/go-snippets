package main

import (
	"html/template"
	"log"
	"net/http"
)

func indexPage(w http.ResponseWriter, r *http.Request) {
	//this must add at begin of every session code
	c, err := r.Cookie("session")
	if err != nil || c.Value == "" {
		http.Error(w, "Session expired", 401)
		return
	}

	//handle GET/POST methods
	if r.Method == "GET" {
		//logout sequence
		if r.FormValue("logout") == "true" {
			cookie := http.Cookie{Name: "session", MaxAge: -1}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/", http.StatusFound)
			log.Println(c.Value + " logout")
			return
		}
	}

	//build page content
	b := `<pre>
	Wellcome to simple applications framework.
	
	This is a plugable environment with modules.
	Modules can be added and removed.
	
	You are now in index module connected as:
		` + c.Value + ` from ` + r.Host +
		`
	
	In the top side you have some demo modules
	
	Enjoy!`

	//finally show the page
	p := Page{
		Title:  "Index page",
		Status: c.Value,
		Body:   template.HTML(b),
	}

	t.ExecuteTemplate(w, "index.html", p)
}

func init() {
	http.HandleFunc("/index", indexPage)
}
