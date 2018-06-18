package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {
	// 没有使用http.HandleFunc 而是自己实现了一个Handler

	// http.HandleFunc("/", sayHello)
	// http.HandleFunc("/bye", sayBye)
	fmt.Println("The simple server starts")
	server := http.Server{
		Addr:        ":8080",
		Handler:     &myHandler{},
		ReadTimeout: time.Second * 5,
	}

	mux = make(map[string]func(http.ResponseWriter, *http.Request))

	mux["/"] = sayHello
	mux["/bye"] = sayBye

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

type myHandler struct{} //定义这个结构的作用是为了能实现ServeHTTP 方法

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I am ServeHTTP")
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world, this is a simple server")
}

func sayBye(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Bye bye, this is a simple server")
}
