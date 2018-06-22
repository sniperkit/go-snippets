package app

import (
	"encoding/json"
	"github.com/gophergala/GopherKombat/common/user"
	"github.com/gorilla/sessions"
	"html/template"
	"io"
	"log"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

const (
	FILE_SERVE_PATH = "/var/static"
)

func InitSession(r *http.Request) *sessions.Session {
	session, err := store.Get(r, "gopher-kombat")
	if err != nil {
		println("ovde")
	}
	return session
}

func GetCurrentUser(r *http.Request) (*user.User, bool) {
	session := InitSession(r)
	if session.Values["user"] == nil {
		return new(user.User), false
	} else {
		return session.Values["user"].(*user.User), true
	}
}

func render(w io.Writer, name string, data interface{}) {
	t, err := template.ParseFiles(FILE_SERVE_PATH + "/template/" + name + ".html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, data)
}

func renderJson(w http.ResponseWriter, r *http.Request, data interface{}) {
	json, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := GetCurrentUser(r)
	data := make(map[string]interface{})
	data["loggedIn"] = ok
	if ok {
		data["user"] = user
	}
	render(w, "home", data)

}
