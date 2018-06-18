package main

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/param"
	"github.com/zenazn/goji/web"
)

func main() {

	goji.Get("/", getForm)
	goji.Post("/post_form", postForm)
	goji.Get("/user/:name", func(c web.C, w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
	})

	goji.Serve()
}

// this handler response a html form page to browser,you can visit in
// http://127.0.0.1:8000/
func getForm(c web.C, w http.ResponseWriter, r *http.Request) {
	form := `
	<form action="/post_form" method="post">
	  First name: <input type="text" name="fname"><br>
	  Last name: <input type="text" name="lname"><br>
	  <input type="submit" value="Submit">
	</form>
	`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, form)
}

// get post params and redirect to other page 
func postForm(c web.C, w http.ResponseWriter, r *http.Request) {
	var user User
	r.ParseForm()
	err := param.Parse(r.Form, &user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url := "/user/" + user.Fname + user.Lname
	http.Redirect(w, r, url, http.StatusFound)
	fmt.Println("this is test")
}


// this is the form & struct mapping
// line 41, params.Parse(r.Form, &user)
type User struct {
	Fname string `param:"fname"`
	Lname string `param:"lname"`
}
