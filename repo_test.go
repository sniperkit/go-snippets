package main

import (
	"fmt"
	"testing"
)

func TestRepoStorage(t *testing.T) {
	defer resetTestData(appDB, "repos", "repo_sources")

	s := &Repo{
		Url:          "test_repo_url",
		Branch:       "test_repo_branch",
		LatestCommit: "test_repo_commit",
	}
	if err := s.Save(appDB); err != nil {
		t.Error(err.Error())
		return
	}

	s2 := &Repo{Id: s.Id}
	if err := s2.Read(appDB); err != nil {
		t.Error(err.Error())
		return
	}

	if err := CompareRepos(s, s2); err != nil {
		t.Error(err)
		return
	}

	if err := s.Delete(appDB); err != nil {
		t.Error(err.Error())
		return
	}
}

func CompareRepos(a, b *Repo) error {
	if a.Id != b.Id {
		return fmt.Errorf("Id mismatch: %s != %s", a.Id, b.Id)
	}
	if !a.Created.Equal(b.Created) {
		return fmt.Errorf("Created mismatch: %s != %s", a.Created, b.Created)
	}
	if !a.Updated.Equal(b.Updated) {
		return fmt.Errorf("Updated mismatch: %s != %s", a.Updated, b.Updated)
	}

	if a.Url != b.Url {
		return fmt.Errorf("Url mismatch: %s != %s", a.Url, b.Url)
	}
	if a.Branch != b.Branch {
		return fmt.Errorf("Branch mismatch: %s != %s", a.Branch, b.Branch)
	}
	if a.LatestCommit != b.LatestCommit {
		return fmt.Errorf("LatestCommit mismatch: %s != %s", a.LatestCommit, b.LatestCommit)
	}

	return nil
}
