package gorlim

import "time"

type Comment struct {
	Author string
	Text   string
	At     *time.Time
}

type Issue struct {
	Id          int
	Opened      bool
	Creator     string
	At          *time.Time
	ClosedAt    *time.Time
	Assignee    string
	Milestone   string
	Title       string
	Description string
	PullRequest string
	Labels      []string
	Comments    []Comment
}
