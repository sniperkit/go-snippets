package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
)

func checkUser(w http.ResponseWriter, r *http.Request) {
	//this must add at begin of every session code
	c, err := r.Cookie("session")
	if err != nil || c.Value == "" {
		http.Error(w, "Session expired", 401)
		return
	}

	//handle GET/POST methods
	var b string

	if r.Method == "POST" {
		db, _ := sql.Open("sqlite3", "./foo.db")
		defer db.Close()

		rows, _ := db.Query("select username from user")
		defer rows.Close()
		var username string
		b = "<p>Username " + r.FormValue("username") + " not exist."
		for rows.Next() {
			rows.Scan(&username)
			if username == r.FormValue("username") {
				b = "<p>Username " + username + " exist."
				break
			}
		}
	} else {

		//build page content
		b = `<pre>
This is a module that works with forms.

<form method="post" action="">
Search a name: <input type="text" size=16 name="username" value="george" autofocus>
               <input type="submit" name="submit" value="Submit" title="(tip: only john and george exist)">
</form>`
	}

	//finally show the page
	p := Page{
		Title:  "Check user page",
		Status: c.Value,
		Body:   template.HTML(b),
	}
	t.ExecuteTemplate(w, "index.html", p)
}

func init() {
	http.HandleFunc("/checkUser", checkUser)
}
