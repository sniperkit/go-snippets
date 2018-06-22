package main

import (
	"encoding/json"

	"github.com/google/go-github/github"
)

type ActivityPayload struct {
	Id         int64   `json:"id"`
	GithubId   string  `json:"github_id"`
	RawPayload RawJson `json:"-"`

	EventType string `json:"event_type"`

	//Events
	OpenPullRequest    bool `json:"open_pull_request"`
	ClosePullRequest   bool `json:"close_pull_request"`
	OpenIssue          bool `json:"open_issue"`
	CreateIssueComment bool `json:"create_issue_comment"`

	//Variables
	Title     string `json:"title"`
	Body      string `json:"body"`
	Number    int    `json:"number"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	ClosedAt  string `json:"closed_at"`
	MergedAt  string `json:"merged_at"`
}

func (p *ActivityPayload) Create() error {
	err := dbmap.Insert(p)
	return err
}

func ProcessPayload(e []github.Event, u User) ([]ActivityPayload, error) {
	activityPayloads := []ActivityPayload{}

	for _, event := range e {
		payload := ActivityPayload{}

		payload.GithubId = *event.ID

		rawPayload := make(map[string]interface{})
		if err := json.Unmarshal(*event.RawPayload, &rawPayload); err != nil {
			panic(err.Error())
		}

		//Save the raw payload
		payload.RawPayload = rawPayload

		switch *event.Type {
		case "PullRequestEvent":
			// Set Pull Request Event type (open or closed)
			if val, ok := rawPayload["action"]; ok && val == "opened" {
				payload.EventType = "OpenPullRequest"
				payload.OpenPullRequest = true
			} else {
				payload.EventType = "ClosePullRequest"
				payload.ClosePullRequest = true
			}

			if pr, ok := rawPayload["pull_request"].(map[string]interface{}); ok {
				if val, ok := pr["title"].(string); ok {
					payload.Title = val
				}
				if val, ok := pr["body"].(string); ok {
					payload.Body = val
				}
				if val, ok := pr["number"].(float64); ok {
					payload.Number = int(val)
				}

				//Set times
				if val, ok := pr["created_at"].(string); ok {
					payload.CreatedAt = val
				}
				if val, ok := pr["updated_at"].(string); ok {
					payload.UpdatedAt = val
				}
				if val, ok := pr["closed_at"].(string); ok {
					payload.ClosedAt = val
				}
				if val, ok := pr["merged_at"].(string); ok {
					payload.MergedAt = val
				}
			}

			// Save the paylod to the db events table
			err := payload.Create()
			if err != nil {
				//Assuming fail due to the unique constraint
			} else {
				activityPayloads = append(activityPayloads, payload)
			}

		case "IssueEvent":
			if val, ok := rawPayload["action"]; ok && val == "opened" {
				payload.EventType = "OpenIssue"
				payload.OpenIssue = true
			}

			if issue, ok := rawPayload["issue"].(map[string]interface{}); ok {
				if val, ok := issue["number"].(float64); ok {
					payload.Number = int(val)
				}
				if val, ok := issue["title"].(string); ok {
					payload.Title = val
				}
				if val, ok := issue["body"].(string); ok {
					payload.Body = val
				}

				//Set times
				if val, ok := issue["created_at"].(string); ok {
					payload.CreatedAt = val
				}
				if val, ok := issue["updated_at"].(string); ok {
					payload.UpdatedAt = val
				}
				if val, ok := issue["closed_at"].(string); ok {
					payload.ClosedAt = val
				}
			}
			activityPayloads = append(activityPayloads, payload)

		case "IssueCommentEvent":
			if val, ok := rawPayload["action"]; ok && val == "created" {
				payload.EventType = "CreateIssueComment"
				payload.CreateIssueComment = true
			}
			if issue, ok := rawPayload["issue"].(map[string]interface{}); ok {
				payload.Number = int(issue["number"].(float64))
				payload.Title = issue["title"].(string)
			}
			if comment, ok := rawPayload["comment"].(map[string]interface{}); ok {
				payload.Body = comment["body"].(string)
			}
			activityPayloads = append(activityPayloads, payload)
		}
	}

	return activityPayloads, nil
}
