package handlers

import (
	"fmt"
	gh "github.com/barakb/create-branch/github"
	"github.com/barakb/create-branch/session"
	"github.com/google/go-github/github"
	"net/http"
)

type LogoutHandler struct {
}

func (h LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LogoutHandler")
	sess := session.GlobalSessions.SessionStart(w, r)
	client := sess.Get("*github.client").(*github.Client)
	//user := sess.Get("user").(*github.User)

	//session.GlobalSessions.SessionDestroy(w, r)
	//http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	var refs []*github.Reference
	refs, err := gh.ListRefs(client, "GigaSpaces", "xap")
	if err != nil {
		fmt.Printf("failed to get heads, error is:%#v\n", err)
	} else {
		for _, ref := range refs {
			fmt.Printf("ref:%s, url:%s, type:%s \n", *ref.Ref, *ref.URL, *ref.Object.Type)
		}
		fmt.Printf("got heads: %#v\n", refs)
	}
	w.WriteHeader(http.StatusOK)
}
