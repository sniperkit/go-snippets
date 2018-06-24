package main

import (
	"testing"
)

func TestReadRepos(t *testing.T) {
	repos, err := ReadRepos(appDB, "created DESC", 10, 0)
	if err != nil {
		t.Error(err)
		return
	}
	if len(repos) == 0 {
		t.Error("no repos returned")
	}
}
