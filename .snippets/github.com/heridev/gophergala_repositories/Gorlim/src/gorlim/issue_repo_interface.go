package gorlim

import "time"

type IssueRepositoryInterface interface {
	GetIssue(id int) (Issue, bool) // not-implemented yet
	GetIssues() ([]Issue, []time.Time)
	Update(string, []Issue, time.Time, *string) 
	Id() int
	Path() string
}

func CreateRepo (repoPath string) IssueRepositoryInterface {
  repo := issueRepository{}
  repo.initialize(repoPath)
  return &repo
}