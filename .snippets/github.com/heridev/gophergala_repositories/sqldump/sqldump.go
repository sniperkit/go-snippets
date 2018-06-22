package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
)

var base_url = "http://localhost"
var database = "information_schema"

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.StatusText(404)
}

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, loginPage)
}

func helpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// TODO help
	fmt.Fprintf(w, "Best viewed with cli-browser >= 6.0\n")
	fmt.Fprintf(w, "select links with mouse or finger\n")
}

func dumpIt(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	db := v.Get("db")
	t := v.Get("t")
	x := v.Get("x")
	user, _, host, port := getCredentials(r)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<p>")
	fmt.Fprint(w, href("/logout", "[X]"))
	fmt.Fprint(w, " &nbsp; ")
	fmt.Fprint(w, href("/help", "[?]"))
	fmt.Fprint(w, " &nbsp; ")
	fmt.Fprint(w, href("/", "[/]"))
	fmt.Fprint(w, " &nbsp; ")
	fmt.Fprint(w, user+"@"+host+":"+port)
	fmt.Fprint(w, " &nbsp; ")

	if db == "" {
		fmt.Fprintln(w, "</p>")
		fmt.Fprint(w, tableA)
		dumpHome(w, r)
		fmt.Fprint(w, tableO)
	} else if t == "" {
		fmt.Fprintln(w, db, "</p>")
		fmt.Fprint(w, tableA)
		dumpTables(w, r, db)
		fmt.Fprint(w, tableO)
	} else if x == "" {
		q := r.URL.Query()
		q.Add("action", "insert")
		linkinsert := q.Encode()
		q.Del("action")
		q.Add("action", "select")
		linkselect := q.Encode()
		q.Del("action")
		q.Del("t")
		linkescape := q.Encode()

		fmt.Fprint(w, db+"."+t)
		fmt.Fprint(w, " &nbsp; ")
		fmt.Fprint(w, " ["+href("?"+linkescape, ".")+"] ")
		fmt.Fprint(w, " ["+href("?"+linkinsert, "+")+"] ")
		fmt.Fprint(w, " ["+href("?"+linkselect, "?")+"] ")
		fmt.Fprintln(w, "</p>")

		fmt.Fprint(w, tableA)
		dumpRecords(w, r, db, t)
		fmt.Fprint(w, tableO)
	} else {
		xint, err := strconv.Atoi(x)
		checkY(err)
		xmax, err := strconv.Atoi(getCount(r, db, t))
		left := strconv.Itoa(maxI(xint-1, 1))
		right := strconv.Itoa(minI(xint+1, xmax))

		q := r.URL.Query()
		q.Set("x", left)
		linkleft := q.Encode()
		q.Set("x", right)
		linkright := q.Encode()
		q.Del("x")
		linkall := q.Encode()

		fmt.Fprint(w, db+"."+t)
		fmt.Fprint(w, " &nbsp; ")
		fmt.Fprint(w, " ["+href("?"+linkall, ".")+"] ")
		fmt.Fprint(w, " ["+href("?"+linkleft, "<")+"] ")
		fmt.Fprint(w, " ["+x+"] ")
		fmt.Fprint(w, " ["+href("?"+linkright, ">")+"] ")
		fmt.Fprintln(w, "</p>")

		fmt.Fprint(w, tableA)
		dumpFields(w, r, db, t, x)
		fmt.Fprint(w, tableO)
	}
}

func workload(w http.ResponseWriter, r *http.Request) {

	v := r.URL.Query()
	action := v.Get("action")
	db := v.Get("db")
	t := v.Get("t")

	if action == "select" && db != "" && t != "" {
		actionSelect(w, r, db, t)
	} else if action == "insert" && db != "" && t != "" {
		actionInsert(w, r, db, t)
	} else {
		dumpIt(w, r)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	pass := ""
	user, _, host, port := getCredentials(r)

	if user != "" {
		workload(w, r)
	} else {
		v := r.URL.Query()
		user = v.Get("user")
		pass = v.Get("pass")
		host = v.Get("host")
		port = v.Get("port")

		if user != "" && pass != "" {
			if host == "" {
				host = "localhost"
			}
			if port == "" {
				port = "3306"
			}
			setCredentials(w, user, pass, host, port)
			http.Redirect(w, r, r.URL.Host, 302)
		} else {
			loginPageHandler(w, r)
		}
	}
}

func main() {

	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/help", helpHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/insert", insertHandler)
	http.HandleFunc("/", indexHandler)

	fmt.Println("Listening at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
