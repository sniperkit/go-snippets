package main

import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
    "log"
)

func main() {
  router := httprouter.New()

  router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Hello!\n")
  })

  log.Fatal(http.ListenAndServe(":8080", router))
}
