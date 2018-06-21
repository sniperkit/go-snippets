package server

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

func Login(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, r.Method+" : "+r.RequestURI)
}

func Test(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%+v\n", r)
	if r.Method == "GET" {
		t, err := template.New("vue-v-for-filter.html").Delims("|||", "|||").ParseFiles("asset/web/html/vue-v-for-filter.html")
		if err != nil {
			log.Println(err)
			io.WriteString(w, err.Error())
			return
		}

		err = t.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
		}
		return

	} else {
		return
	}
	return
}

func Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Login)
	mux.HandleFunc("/test", Test)
	mux.Handle("/asset/web/", http.FileServer(http.Dir(".")))

	contextedMux := AddContextSupport(mux)
	log.Fatal(http.ListenAndServe(":6080", contextedMux))
}

func AddContextSupport(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/asset") {
			next.ServeHTTP(w, r)
			return
		} else {
			log.Println(r.Method, "-", r.RequestURI)
			cookie, _ := r.Cookie("SESSIONID")
			if cookie != nil {
				ctx := context.WithValue(r.Context(), "SESSIONID", cookie.Value)
				// WithContext returns a shallow copy of r with its context changed
				// to ctx. The provided ctx must be non-nil.
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			next.ServeHTTP(w, r)
		}
	})
}
