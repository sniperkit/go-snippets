package main

import (
        "fmt"
        "net/http"

        "github.com/zenazn/goji"
        "github.com/zenazn/goji/web"
        "github.com/zenazn/goji/web/middleware"
)

func hello(c web.C, w http.ResponseWriter, r *http.Request) {
        //fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
        fmt.Fprintf(w, `{"message": "hello world", "status": %d}`, http.StatusOK)
}

func main() {
        goji.Abandon(middleware.Logger)
        goji.Get("/hello/:name", hello)
        goji.Serve()
}
