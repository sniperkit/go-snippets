package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const DirName = "snippets"
const Placeholder = "{{%s}}"

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/questions", questionHandler)
	http.HandleFunc("/highscores", scoresHandler)
	http.HandleFunc("/submitQuestion", answerHandler)
	http.HandleFunc("/submitScore", scoreHandler)
	http.HandleFunc("/authorize", handleAuthorize)
	http.HandleFunc("/logout", handleLogout)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Listening on 8080")
	http.ListenAndServe(":8080", nil)
}

func mainHandler(rw http.ResponseWriter, req *http.Request) {
	f, _ := ioutil.ReadFile("templates/layout.html")
	m := make(map[string]interface{})
	user := CurrentUser(req)
	if user == nil {
		m["loginUrl"] = "/authorize"
	} else {
		m["username"] = user.Username
		m["login"] = true
		m["logoutUrl"] = "/logout"
	}
	b := templateToByte(f, m)
	rw.Write(b)
}

func faviconHandler(rw http.ResponseWriter, req *http.Request) {
	http.ServeFile(rw, req, "./static/favicon.ico")
}

func questionHandler(rw http.ResponseWriter, req *http.Request) {
	f, _ := os.Open(DirName)

	stat, _ := f.Stat()
	if !stat.IsDir() {
		serveErr(rw, req, errors.New("Invalid path"))
	}
	fileList := []string{}
	err := filepath.Walk(DirName, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == ".template" {
			fileList = append(fileList, path)
		}
		return nil
	})

	if err != nil {
		serveErr(rw, req, err)
		return
	}
	var list []interface{}
	for _, filePath := range fileList {
		snippet, _ := ioutil.ReadFile(filePath)
		question, _ := ioutil.ReadFile(filePath + ".question")
		m := make(map[string]string)
		m["snippet"] = strings.Replace(string(snippet), Placeholder, "[ snippet will be inserted here ]", 1)
		m["question"] = string(question)
		m["id"] = filePath
		list = append(list, m)
	}
	serveJson(rw, req, list)
}

func answerHandler(rw http.ResponseWriter, req *http.Request) {
	answer := req.FormValue("answer")
	snippetId := req.FormValue("id")
	if answer == "" || snippetId == "" {
		serveErrJson(rw, req, errors.New("Invalid Request"), "")
		return
	}
	output, err := process(snippetId, answer)
	if err != nil {
		serveErrJson(rw, req, err, output)
		return
	}
	m := map[string]interface{}{"message": "Correct code! Click next to continue.", "output": output}
	serveJson(rw, req, m)
}

func scoreHandler(rw http.ResponseWriter, req *http.Request) {
	score := req.FormValue("score")
	if req.Method != "POST" || score == "" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Bad Request"))
		return
	}
	saved, err := SaveHighScore(score, req)
	if err != nil {
		serveErrJson(rw, req, err, "")
		return
	}
	m := map[string]interface{}{"saved": saved}
	serveJson(rw, req, m)
}

func scoresHandler(rw http.ResponseWriter, req *http.Request) {
	scores := GetHighScores()
	serveJson(rw, req, scores)
}

func serveJson(rw http.ResponseWriter, req *http.Request, jsonBody interface{}) {
	m := make(map[string]interface{})
	m["status"] = "success"
	m["data"] = jsonBody
	b, _ := json.Marshal(m)
	rw.Header().Set("Content-type", "application/json")
	rw.Write(b)
}

func serveErrJson(rw http.ResponseWriter, req *http.Request, err error, output string) {
	m := make(map[string]interface{})
	m["status"] = "failure"
	m["message"] = err.Error()
	if err == WrongOutput {
		m["output"] = output
	}
	b, _ := json.Marshal(m)
	rw.Header().Set("Content-type", "application/json")
	rw.Write(b)
}

func serveErr(rw http.ResponseWriter, req *http.Request, err error) {
	rw.WriteHeader(500)
	rw.Write([]byte(err.Error()))
}

func templateToByte(contents []byte, data interface{}) []byte {
	buf := new(bytes.Buffer)
	t := template.Must(template.New("T").Parse(string(contents)))
	t.Execute(buf, data)
	return buf.Bytes()
}
