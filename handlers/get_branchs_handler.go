package handlers

import (
	"encoding/json"
	gh "github.com/barakb/create-branch/github"
	"github.com/barakb/create-branch/session"
	"github.com/google/go-github/github"
	"log"
	"net/http"
)

type GetBranchsHandler struct {
	WS    *WebSocketHandler
	Repos []string
}

type Branches struct {
	Map   map[string]map[string]string `json:"branches"`
	Repos []string                     `json:"repos"`
	Login string                       `json:"login"`
}

func (h GetBranchsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sess := session.GlobalSessions.SessionStart(w, r)
	client := sess.Get("*github.client").(*github.Client)
	firstEventId, eventsMap, deletedMap, err := gh.GetEvents("xap", "xap", nil, client)
	if err == nil {
		log.Println("")
		log.Printf("firstEventId: %s\n", *firstEventId)
		log.Println("New branches")
		for _, value := range eventsMap {
			log.Printf("Created branch: %s by %s\n", value.Branch, value.Login)
		}
		log.Println("Events")
		for key, _ := range deletedMap {
			log.Printf("Deleted branch : %s\n", key)
		}
		log.Println("")
	}
	resChan := gh.GetEventsForRepos(client)
	branches := <-gh.ListAllRefsAsMap(client)
	//map[string]map[string]*BranchOwnershipEvent
	creationEventsByRepo := <-resChan
	log.Printf("Merging events %v\n", creationEventsByRepo)
	for branchName, branchRepos := range branches {
		for repoName, _ := range branchRepos {
			if repoEvents, ok := creationEventsByRepo[repoName]; ok {
				if event, ok := repoEvents[branchName]; ok {
					branches[branchName][repoName] = event.Login
				}
			}
		}
	}

	js, err := json.Marshal(Branches{branches, h.Repos, *sess.Get("user").(*github.User).Login})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
