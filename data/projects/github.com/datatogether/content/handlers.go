package main

import (
	"fmt"
	"github.com/datatogether/core"
	"io"
	"net/http"
	"strconv"
)

func reqParamInt(key string, r *http.Request) (int, error) {
	i, err := strconv.ParseInt(r.FormValue(key), 10, 0)
	return int(i), err
}

func reqParamBool(key string, r *http.Request) (bool, error) {
	return strconv.ParseBool(r.FormValue(key))
}

func DownloadUrlHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		EmptyOkHandler(w, r)
	case "GET":
		DownloadUrl(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

func DownloadUrl(w http.ResponseWriter, r *http.Request) {
	u := &core.Url{Id: r.URL.Path[len("/urls/"):]}
	if err := u.Read(store); err != nil {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, err.Error())
		return
	}

	f, err := u.File()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, err.Error())
		return
	}

	if err := f.GetS3(); err != nil {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, err.Error())
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, u.FileName))
	w.Write(f.Data)
}

// HealthCheckHandler is a basic "hey I'm fine" for load balancers & co
// TODO - add Database connection & proper configuration checks here for more accurate
// health reporting
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "status" : 200 }`))
}

// EmptyOkHandler is an empty 200 response, often used
// for OPTIONS requests that responds with headers set in addCorsHeaders
func EmptyOkHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// CertbotHandler pipes the certbot response for manual certificate generation
func CertbotHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, cfg.CertbotResponse)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{ "status" :  "not found" }`))
}
