package main

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji"
)

func main() {
	goji.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "root page")
	})

	// set cookie
	// doc: http://godoc.org/net/http#SetCookie
	// visit http://127.0.0.1:8000/world will set cookie
	// cookie key is language, value is golang
	goji.Get("/world", func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{Name: "language", Value: "golang"}
		// expire := time.Now().AddDate(0, 0, 1)
		// cookie := &http.Cookie{Name: "language", Value: "golang", Expires: expire}

		http.SetCookie(w, cookie)
		fmt.Fprint(w, "world page")
	})

	// get cookiek
	// doc: http://godoc.org/net/http#Request.Cookie
	// visit http://127.0.0.1:8000/hello will get cookie
	// the response content is "hello page, cookie is: golang"
	goji.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		// get cookie value
		cookie, _ := r.Cookie("language")

		if cookie != nil {
			fmt.Fprintf(w, "hello page, cookie is: %s", cookie.Value)
		} else {
			fmt.Fprint(w, "hello page")
		}
	})

	goji.Serve()
}
