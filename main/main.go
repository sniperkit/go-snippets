package main

import (
	"flag"
	"fmt"
	gh "github.com/barakb/create-branch/github"
	"github.com/barakb/create-branch/handlers"
	"github.com/barakb/create-branch/session"
	"github.com/vharitonsky/iniflags"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func init() {
	gs, err := session.NewManager("memory", "gosessionid", 10*1000*1000*1000)
	if err != nil {
		fmt.Printf("failed to create session manager error is: %#v\n", err)
		return
	}
	session.GlobalSessions = gs
}

var port = flag.Int("port", 4430, "Configure the server port")

var repos = flag.String("repos", defaultRepoFile(), "Configure the where the repositories file")

func main() {

	iniflags.Parse()
	repos, err := readRepos(*repos)
	if err != nil {
		fmt.Printf("Error:%s\n", err.Error())
		return
	}
	fmt.Printf("repos read from %s\n", repos)
	gh.ReposNames = repos
	go session.GlobalSessions.GC()
	websocketHandler := &handlers.WebSocketHandler{}
	http.Handle("/", handlers.MustAuth(handlers.MainHandler{File: "index.html"}))
	http.Handle("/events", handlers.MustAuth(websocket.Handler(websocketHandler.Handler)))
	http.Handle("/logout", handlers.MustAuth(handlers.LogoutHandler{}))
	http.Handle("/githuboa_cb", handlers.GithubLoginHandler{})
	http.Handle("/api/create_branch/", handlers.MustAuth(handlers.CreateBranchHandler{WS: websocketHandler}))
	http.Handle("/api/delete_branch/", handlers.MustAuth(handlers.DeleteBranchHandler{WS: websocketHandler}))
	http.Handle("/api/get_branches/", handlers.MustAuth(handlers.GetBranchsHandler{WS: websocketHandler, Repos: repos}))
	http.Handle("/web/", handlers.MustAuth(http.StripPrefix("/web/", http.FileServer(http.Dir("web")))))
	fmt.Printf("Server started on https://localhost:%d\n", *port)
	fmt.Println(http.ListenAndServeTLS(fmt.Sprintf(":%d", *port), "keys/server.pem", "keys/server.key", nil))
}

func defaultRepoFile() string {
	dir, err := os.Getwd()
	if err != nil {
		return "repos.txt"
	}
	return fmt.Sprintf("%s/repos.txt", dir)
}

func readRepos(repos string) ([]string, error) {
	bytes, err := ioutil.ReadFile(repos)
	if err != nil {
		return nil, err
	}
	str := strings.TrimSpace(string(bytes))
	fmt.Printf("repos are: %s\n", str)
	return regexp.MustCompile("\\s+").Split(str, -1), nil
}
