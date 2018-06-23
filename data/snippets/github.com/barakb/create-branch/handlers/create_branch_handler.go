package handlers

import (
	"encoding/json"
	"fmt"
	gh "github.com/barakb/create-branch/github"
	"github.com/barakb/create-branch/session"
	"github.com/google/go-github/github"
	"net/http"
	"strings"
	"strconv"
)

type CreateBranchHandler struct {
	WS *WebSocketHandler
}

func (h CreateBranchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sess := session.GlobalSessions.SessionStart(w, r)
	branchName := strings.Trim(strings.Replace(r.URL.Path, "/api/create_branch/", "", 1), "/")
	from := r.URL.Query().Get("from")
	xapOnly := r.URL.Query().Get("isXAPOnly")
	fmt.Printf("creating branch %s from %s, xapOnly: %s\n", branchName, from, xapOnly)
	isXapOnly, err := strconv.ParseBool(xapOnly)
	if err != nil {
		isXapOnly = false
	}

	if strings.Contains(branchName, "/") {
		fmt.Printf("CreateBranchHandler: wrong request: %q\n", r.RequestURI)
		http.Error(w, fmt.Sprintf("CreateBranchHandler: wrong request: %q\n", r.RequestURI), http.StatusInternalServerError)
		return

	}

	client := sess.Get("*github.client").(*github.Client)

	progressChan, resChan := gh.CreateBranchsWithProgress(branchName, from, isXapOnly, client)

	go func() {
		for progress := range progressChan {
			js, err := json.Marshal(progress)
			if err == nil {
				fmt.Fprintf(h.WS.Conn(), "{'type' : 'create-branch-progress'\n 'branch' : %q,\n 'progress' : %s}\n", branchName, js)
			}
		}
	}()

	created := <-resChan
	js, err := json.Marshal(created)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
