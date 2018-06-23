package github

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"log"
	"sort"
	"strings"
	"time"
	"context"
)

//{"ref":"yaelna_120","ref_type":"branch","pusher_type":"user"}
type EventData struct {
	Ref        string `json:"ref"`
	RefType    string `json:"ref_type"`
	PusherType string `json:"pusher_type"`
}

type BranchOwnershipEvent struct {
	Owner, Repo string
	Login       string
	CreatedAt   time.Time
	Branch      string
	Type        string
	Id          string
}

type SortByDateBranchOwnershipEvents []*BranchOwnershipEvent

func (s SortByDateBranchOwnershipEvents) Len() int {
	return len(s)
}
func (s SortByDateBranchOwnershipEvents) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortByDateBranchOwnershipEvents) Less(i, j int) bool {
	return s[i].CreatedAt.Before(s[j].CreatedAt)
}

func GetEvents(owner, repo string, lastEventId *string, client *github.Client) (firstEventId *string, eventsMap map[string]*BranchOwnershipEvent, deletedMap map[string]bool, err error) {
	var page int = 0
	eventsMap = make(map[string]*BranchOwnershipEvent)
	deletedMap = make(map[string]bool)
	ownershipEvents := make([]*BranchOwnershipEvent, 0)
	for {
		var events []*github.Event
		var resp *github.Response
		opt := &github.ListOptions{Page: page, PerPage: 100}
		events, resp, err = client.Activity.ListRepositoryEvents(context.Background(), owner, repo, opt)
		if err != nil {
			log.Printf("error geting repository events: %v\n", err)
			return
		}
		for _, event := range events {
			if (lastEventId != nil) && (*event.ID == *lastEventId) {
				page = resp.LastPage
				break
			}
			if firstEventId == nil {
				firstEventId = event.ID
			}
			if (*event.Type == "CreateEvent") || (*event.Type == "DeleteEvent") {
				var bytes []byte
				bytes, err = event.RawPayload.MarshalJSON()
				if err != nil {
					log.Printf("error geting parsing event data: %v\n", err)
					continue
				}
				data := EventData{}
				err = json.Unmarshal(bytes, &data)
				if err != nil {
					log.Printf("error geting parsing event data: %v\n", err)
					continue
				}
				if data.RefType == "branch" && data.PusherType == "user" {
					oe := &BranchOwnershipEvent{owner, repo, *event.Actor.Login, *event.CreatedAt, data.Ref, *event.Type, *event.ID}
					//log.Printf("[%v] Seeing event: %v\n", oe.CreatedAt, oe)
					ownershipEvents = append(ownershipEvents, oe)
				}
			}
		}
		if page < resp.LastPage {
			page = resp.NextPage
		} else {
			break
		}
	}

	// sort events
	sort.Sort(sort.Reverse(SortByDateBranchOwnershipEvents(ownershipEvents)))
	seenMap := make(map[string]bool)

	// loop over events
	for _, event := range ownershipEvents {
		//log.Printf("[%v] Handling event: %v\n", event.CreatedAt, event)
		if seenMap[event.Branch] {
			continue
		}
		seenMap[event.Branch] = true
		if event.Type == "CreateEvent" {
			eventsMap[event.Branch] = event
		} else {
			deletedMap[event.Branch] = true
		}
	}
	return
}

func GetEventsForRepos(client *github.Client) (created chan map[string]map[string]*BranchOwnershipEvent) {

	type repoBranchRes struct {
		branches map[string]*BranchOwnershipEvent
		repo     string
	}
	perRepoBranches := make(chan repoBranchRes, len(ReposNames))
	created = make(chan map[string]map[string]*BranchOwnershipEvent)
	for _, repo := range ReposNames {
		sli := strings.Split(repo, "/")
		owner, repo := sli[0], sli[1]
		go func(owner, repo string) {
			_, eventsMap, _, _ := GetEvents(owner, repo, nil, client)
			perRepoBranches <- repoBranchRes{eventsMap, fmt.Sprintf("%s/%s", owner, repo)}
		}(owner, repo)
	}
	returnedMap := make(map[string]map[string]*BranchOwnershipEvent)
	go func() {
		for range ReposNames {
			repoBranchs := <-perRepoBranches
			returnedMap[repoBranchs.repo] = repoBranchs.branches

		}
		created <- returnedMap
	}()
	return created
}
