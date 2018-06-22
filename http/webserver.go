package main

import (
  "encoding/json"
  "fmt"
  "net/http"
)

func main() {
  http.HandleFunc("/", handler)
  http.HandleFunc("/me/", me)
  fmt.Println("server is running on port 8080")
  http.ListenAndServe(":8080", nil)
}

type Message struct {
    Text string
}

func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("request url: " + r.URL.Path)
  fmt.Fprintf(w, "Welcome, %s!", r.URL.Path[1:])
}

func me (w http.ResponseWriter, r *http.Request) {
  m := Message{"last updated: just now"}
  b, err := json.Marshal(m)

  if err != nil {
      panic(err)
  }
  w.Write(b)
}
