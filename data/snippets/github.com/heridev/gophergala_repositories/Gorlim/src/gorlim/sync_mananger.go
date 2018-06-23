package gorlim

import "time"
import "fmt"

type GitWebPair struct {
	repo               IssueRepositoryInterface
	uri                string // TBD corresponding object
	webUpdateTimestamp time.Time
}

type IssuesUpdate struct {
	Uri    string
	Issues []Issue
}

type SyncManager struct {
	idToReposMap  map[int]GitWebPair
	uriToReposMap map[string]GitWebPair
}

// TBD: first parameter should be web issues interface
func (sm *SyncManager) AddRepository(webIssuesUri string, repo IssueRepositoryInterface) {
	gwp := GitWebPair{repo: repo, uri: webIssuesUri, webUpdateTimestamp: time.Now()}
	sm.idToReposMap[repo.Id()] = gwp
	sm.uriToReposMap[webIssuesUri] = gwp
}

func Create() *SyncManager {
	return &SyncManager{
		idToReposMap:  make(map[int]GitWebPair),
		uriToReposMap: make(map[string]GitWebPair),
	}
}

// TBD: idea is that we don't nee third parameter is first paramter will be real WebIssue interface with getIssues method
func (sm *SyncManager) InitGitRepoFromIssues(webIssuesUri string, repo IssueRepositoryInterface, issues []Issue) {
	for _, issue := range issues {
		issue1 := issue
		issue1.Comments = []Comment{}
		issue1.Opened = true
		repo.Update("webimport Opened issue: " + issue.Title + " " + issue.Description, []Issue{issue1}, *issue.At, &issue.Creator)
		for i := 0; i < len(issue.Comments); i++ {
			issue1.Comments = issue.Comments[0:i] 
			repo.Update(fmt.Sprintf("webimport: #%v", issue.Comments[i].Text), []Issue{issue1}, *issue.Comments[i].At, &issue.Comments[i].Author)
		}
		if issue.Opened == false {
			if issue.Assignee == "" {
				repo.Update("webimport Closed issue: " + issue.Title, []Issue{issue}, *issue.ClosedAt, nil)
			} else {
				repo.Update("webimport Closed issue: " + issue.Title, []Issue{issue}, *issue.ClosedAt, &issue.Assignee)
			}
		}
	}
	gwp := sm.idToReposMap[repo.Id()]
	gwp.webUpdateTimestamp = time.Now()
	sm.idToReposMap[repo.Id()] = gwp
}

func (sm *SyncManager) SubscribeToPushEvent(pushevent <-chan int) {
	go func() {
		for push := range pushevent {
			// TBD here we can simply send current repo state to web interface
			gwp := sm.idToReposMap[push]
			repo := gwp.repo
			issues, timestamps := repo.GetIssues()
			currentTime := time.Now()
			for index, tm := range timestamps {
				// if modified later than last sync
				if time.Since(tm) < time.Since(gwp.webUpdateTimestamp) {
					issue := issues[index]
					fmt.Println("Pushed issue", issue)
					// TODO : send issue to web
				}
			}
			gwp.webUpdateTimestamp = currentTime
			sm.idToReposMap[push] = gwp
		}
	}()
}

// Simple implementation of web-to-git updater - do not care that commit may come from the user in the same time for starters
func (sm *SyncManager) SubscribeToWebUpdateEvent(webupdate <-chan IssuesUpdate) {
	go func() {
		for wupd := range webupdate {
			uri := wupd.Uri
			issues := wupd.Issues
			repo := sm.uriToReposMap[uri].repo
			fmt.Println(uri)
			for _, issue := range issues {
				repo.Update("webimport: "+issue.Title, []Issue{issue}, time.Now(), nil)
			}
		}
	}()
}
