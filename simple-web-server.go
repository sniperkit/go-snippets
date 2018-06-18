package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", Cookie) // 1. http.HandleFunc

	http.ListenAndServe(":9090", nil) //2. http.ListenAndServe
}

func Cookie(w http.ResponseWriter, r *http.Request) { //3. the function in HandleFunc
	ck := &http.Cookie{
		Name:   "myCookie",
		Value:  "hello",
		Path:   "/",
		Domain: "localhost",
		MaxAge: 120,
	}

	http.SetCookie(w, ck)
	ck2, err := r.Cookie("uuid")

	if err != nil {
		io.WriteString(w, err.Error())
		fmt.Println("The cookie is not exit")
		return
	}

	io.WriteString(w, ck2.Value)
}
