package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func loginPage(w http.ResponseWriter, r *http.Request) {

	//handle GET/POST methods
	if r.Method == "POST" {
		db, _ := sql.Open("sqlite3", "./foo.db")
		defer db.Close()

		rows, _ := db.Query("select username,password from user")
		defer rows.Close()
		for rows.Next() {
			var username, password string
			rows.Scan(&username, &password)
			if username == r.FormValue("username") && password == r.FormValue("password") {
				cookie := http.Cookie{Name: "session", Value: r.FormValue("username")}
				http.SetCookie(w, &cookie)
				http.Redirect(w, r, "/index", http.StatusFound)
				log.Println(r.FormValue("username") + " logged in")
				return
			}
		}
		log.Println("user " + r.FormValue("username") + " authentication failed")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	//finally show the page
	p := Page{
		Title: "Login",
	}
	t.ExecuteTemplate(w, "login.html", p)
}

func init() {
	http.HandleFunc("/", loginPage)
}
