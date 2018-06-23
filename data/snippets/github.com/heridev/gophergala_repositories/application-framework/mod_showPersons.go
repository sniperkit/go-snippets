package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"net/http"
)

func showPersons(w http.ResponseWriter, r *http.Request) {
	//this must add at begin of every session code
	c, err := r.Cookie("session")
	if err != nil || c.Value == "" {
		http.Error(w, "Session expired", 401)
		return
	}

	//build page content
	b := `<p>This is a module that show a list of persons from a database`

	db, _ := sql.Open("sqlite3", "./foo.db")
	defer db.Close()

	rows, err := db.Query("select id, name, age, address from person")
	if err != nil {
		println(err.Error())
		return
	}
	defer rows.Close()
	var name, address string
	var id, age int

	b += `<table width="500">
	<tr><td>Id</td><td>Name</td><td>Age</td><td>Address</td></tr>
	<tr><td colspan=4><hr></td></tr>`
	for rows.Next() {
		rows.Scan(&id, &name, &age, &address)
		b += fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%d</td><td>%s</td></tr>", id, name, age, address)
	}
	b += `<tr><td colspan=4><hr></td></tr>
	</table>`

	b += `<a href="/index">OK</a>`

	//finally show the page
	p := Page{
		Title:  "Show persons page",
		Status: c.Value,
		Body:   template.HTML(b),
	}
	t.ExecuteTemplate(w, "index.html", p)
}

func init() {
	http.HandleFunc("/showPersons", showPersons)
}
