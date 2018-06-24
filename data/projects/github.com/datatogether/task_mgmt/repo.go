package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ipfs/go-datastore"
	"github.com/pborman/uuid"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// a Repo represents a github repository that can be applied to sources
type Repo struct {
	// uuid identifier for task
	Id string `json:"id"`
	// created date rounded to secounds
	Created time.Time `json:"created"`
	// updated date rounded to secounds
	Updated time.Time `json:"updated"`
	// location of code. This should be a github url
	Url string
	// Branch to read From
	// For now this should always just be "master", url should match any
	// branching scheme
	Branch string `json:"branch"`
	// version control commit to execute code from
	// currently should only look at the master branch
	LatestCommit string `json:"codeCommit"`
}

func (r *Repo) Owner() string {
	url, err := url.Parse(r.Url)
	if err != nil {
		return ""
	}
	spl := strings.Split(url.Path, "/")
	if len(spl) >= 1 {
		return spl[1]
	}

	return ""
}

func (r *Repo) Repo() string {
	url, err := url.Parse(r.Url)
	if err != nil {
		return ""
	}
	spl := strings.Split(url.Path, "/")
	if len(spl) >= 2 {
		return spl[2]
	}

	return ""
}

func (r *Repo) FetchLatestCommit() (string, error) {
	u := fmt.Sprintf("http://api.github.com/repos/%s/%s/branches/%s", r.Owner(), r.Repo(), r.Branch)
	res, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid Github API response code while fetching latest commit: %d. url: %s", res.StatusCode, u)
	}

	body := map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return "", err
	}

	if commit, ok := body["commit"].(map[string]interface{}); ok {
		if commitSha, ok := commit["sha"].(string); ok {
			return commitSha, nil
		}
	}

	return "", fmt.Errorf("malformed github response: %s", body)
}

func (r *Repo) Read(db sqlQueryable) error {
	if r.Id == "" {
		return datastore.ErrNotFound
	}
	return r.UnmarshalSQL(db.QueryRow(qRepoReadById, r.Id))
}

func (r *Repo) Save(db sqlQueryExecable) error {
	prev := &Repo{Id: r.Id}
	if err := prev.Read(db); err == datastore.ErrNotFound {
		r.Id = uuid.New()
		r.Created = time.Now().Round(time.Millisecond).In(time.UTC)
		r.Updated = r.Created
		_, err := db.Exec(qRepoInsert, r.sqlArgs()...)
		return err
	} else if err != nil {
		return err
	} else {
		r.Updated = time.Now().Round(time.Millisecond).In(time.UTC)
		_, err := db.Exec(qRepoUpdate, r.sqlArgs()...)
		return err
	}

	return nil
}

func (r *Repo) Delete(db sqlQueryExecable) error {
	_, err := db.Exec(qRepoDelete, r.Id)
	return err
}

func (t *Repo) UnmarshalSQL(row sqlScannable) error {
	var (
		id, url, branch, commit string
		created, updated        time.Time
	)

	err := row.Scan(&id, &created, &updated, &url, &branch, &commit)
	if err != nil {
		if err == sql.ErrNoRows {
			return datastore.ErrNotFound
		}
		return err
	}

	*t = Repo{
		Id:           id,
		Created:      created,
		Updated:      updated,
		Url:          url,
		Branch:       branch,
		LatestCommit: commit,
	}

	return nil
}

func (t *Repo) sqlArgs() []interface{} {
	return []interface{}{
		t.Id,
		t.Created,
		t.Updated,
		t.Url,
		t.Branch,
		t.LatestCommit,
	}
}
