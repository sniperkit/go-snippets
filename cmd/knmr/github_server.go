package main

import (
	"net/http"

	"github.com/Ken2mer/knmr/logger"
	"github.com/google/go-github/github"
)

type gitHubEventMonitor struct {
	webhookSecretKey []byte
}

// cf. https://godoc.org/github.com/google/go-github/github#ParseWebHook
func (s *gitHubEventMonitor) serveHTTP(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, s.webhookSecretKey)
	if err != nil {
		logger.Errorf("error: %s", err)
	}
	var event interface{}
	event, err = github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		logger.Errorf("error: %s", err)
	}
	switch event := event.(type) {
	case *github.CreateEvent:
		logger.Debugf("create: %s\n", event)
	}
}

func ghServe() error {
	// var webhookSecretKey string = "XXX"
	s := gitHubEventMonitor{
		webhookSecretKey: []byte(webhookSecretKey),
	}
	http.HandleFunc("/payload", s.serveHTTP)
	return http.ListenAndServe(":12345", nil)
}
