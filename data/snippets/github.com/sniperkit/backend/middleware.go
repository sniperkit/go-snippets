package backend

import (
	"fmt"
	"github.com/justinas/alice"
	"net/http"
	"time"
)

var mw = alice.New(logger)

func logger(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Logging begin: url =", r.URL)
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Since(t1)
		fmt.Println("Logging: request duration", t2)
		fmt.Println("Logging end")
	})

}
