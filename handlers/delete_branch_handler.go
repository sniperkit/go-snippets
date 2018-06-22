package handlers

import (
	"encoding/json"
	"fmt"
	gh "github.com/barakb/create-branch/github"
	"github.com/barakb/create-branch/session"
	"github.com/google/go-github/github"
	"net/http"
	"strings"
)

type DeleteBranchHandler struct {
	WS *WebSocketHandler
}

func (h DeleteBranchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sess := session.GlobalSessions.SessionStart(w, r)
	branchName := strings.Trim(strings.Replace(r.RequestURI, "/api/delete_branch/", "", 1), "/")
	if strings.Contains(branchName, "/") || r.Method != "DELETE" {
		fmt.Printf("DeleteBranchHandler: wrong request: %s %q\n", r.Method, r.RequestURI)
		http.Error(w, fmt.Sprintf("DeleteBranchHandler: wrong request: %s %q\n", r.Method, r.RequestURI), http.StatusInternalServerError)
		return
	}

	client := sess.Get("*github.client").(*github.Client)
	progressChan, resChan := gh.DeleteBranchWithProgress(branchName, client)

	go func() {
		for progress := range progressChan {
			fmt.Printf("Reporting progress on deleted branch %#v\n", progress)
			js, err := json.Marshal(progress)
			//fmt.Printf("{'type' : 'delete-branch-progress'\n 'branch' : %q,\n 'progress' : %s}\n", branchName, js)
			if err == nil {
				fmt.Fprintf(h.WS.Conn(), "{'type' : 'delete-branch-progress'\n 'branch' : %q,\n 'progress' : %s}\n", branchName, js)
			}

		}
	}()

	res := <-resChan
	js, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
