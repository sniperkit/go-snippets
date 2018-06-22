/*
* @Author: souravray
* @Date:   2015-01-24 10:35:13
* @Last Modified by:   souravray
* @Last Modified time: 2015-01-25 22:38:41
 */

package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/sessions"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

type Config struct {
	DbURI  string `json:"db_uri"`
	DbName string `json:"db_name"`
}

var conf Config

var store = sessions.NewCookieStore([]byte("DellStudiow-Laptop-as-Desktop"))

var templates = template.Must(template.ParseFiles(
	"static/template/landing.html",
	"static/template/error.html",
	"static/template/login.html",
	"static/template/signup.html",
	"static/template/logout.html",
	"static/template/CPG.html",
	"static/template/CBG.html",
	"static/template/badge.html",
	"static/template/campaignregistration.html",
))

func render(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func dispatchError(w http.ResponseWriter, message string) {
	err := errors.New(message)
	http.Error(w, err.Error(), 500)
}

func dispatchJSON(w http.ResponseWriter, response interface{}) {
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(response)
}

func init() {
	configpath := os.Getenv("TE_CONF_PATH")
	if configpath == "" {
		configpath = "/tmp/config.json"
	}
	// read condig file
	content, err := ioutil.ReadFile(configpath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(content, &conf)
	if err != nil {
		panic(err)
	}
}

func IsAuth(req *http.Request) bool {
	websession, _ := store.Get(req, "pp-session")

	if websession.Values["id"] != nil {
		return true
	}
	return false
}

// controler for site root
func LandingPage(rw http.ResponseWriter, req *http.Request) {
	if IsAuth(req) {
		http.Redirect(rw, req, "/campaign/create", 301)
	}
	render(rw, "landing.html")
	return
}
