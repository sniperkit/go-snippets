package urls

import (
	"log"
	"net/http"
	"time"
)

func Middleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(w, r)
	log.Println("url:", r.RequestURI, "execute time:", time.Since(start))
}
