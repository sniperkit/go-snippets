package gist

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/datatogether/core"
	"github.com/datatogether/sql_datastore"
	"github.com/datatogether/task_mgmt/tasks"
	"github.com/ipfs/go-datastore"
)

const collectionInfoFilename = "collection.json"
const urlsFilename = "urls.txt"

// CollectionFromGist creates a collection from a gist of urls
// the gist requires a file called urls.txt be defined, and be
// one-url-per-line. an optional collection.json sets
// the title, description, and url properties of the collection
type CollectionFromGist struct {
	// title for collection if no collection present
	GistUrl string `json:"gistUrl"`
	// author of gist
	CreatorId string `json:"creatorId"`
	// internal datastore pointer
	store datastore.Datastore
}

func NewCollectionFromGist() tasks.Taskable {
	return &CollectionFromGist{}
}

// CollectionFromGist task needs to talk to an underlying database
// it's expected that the task executor will call this method
// before calling Do
func (t *CollectionFromGist) SetDatastore(store datastore.Datastore) {
	if sqlds, ok := store.(*sql_datastore.Datastore); ok {
		// if we're passed an sql datastore
		// make sure our collection model is registered
		sqlds.Register(
			&core.Url{},
			&core.Collection{},
			&core.CollectionItem{},
		)
	}

	t.store = store
}

func (t *CollectionFromGist) Valid() error {
	if t.GistUrl == "" {
		return fmt.Errorf("gistUrl is required")
	}
	return nil
}

func (t *CollectionFromGist) Do(pch chan tasks.Progress) {
	p := tasks.Progress{Step: 1, Steps: 1, Status: "creating collection"}
	pch <- p

	id, err := GistIdFromUrl(t.GistUrl)
	if err != nil {
		p.Error = err
		pch <- p
		return
	}
	fmt.Println(id)

	col, err := CollectionFromGistId(t.store, id, t.CreatorId)
	if err != nil {
		p.Error = err
		pch <- p
		return
	}

	p.Status = fmt.Sprintf("created collection: %s", col.Id)
	p.Percent = 1.0
	p.Done = true
	pch <- p
	return
}

// GistIdFromUrl extracts the url from a string that is either
// a url to the gist, or a raw id.
func GistIdFromUrl(rawurl string) (string, error) {
	if len(rawurl) == 32 && !strings.HasPrefix(rawurl, "http") {
		return rawurl, nil
	}
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", fmt.Errorf("invalid gist url: %s", err.Error())
	}

	id := filepath.Base(u.Path)
	if len(id) != 32 {
		return id, fmt.Errorf("invalid gist url or id: '%s'", id)
	}
	return id, nil
}

// CollectionFromGistId creates a core.Collection from a gist
func CollectionFromGistId(store datastore.Datastore, gistid, creatorId string) (*core.Collection, error) {
	col := &core.Collection{Creator: creatorId}
	res, err := http.Get(fmt.Sprintf("https://api.github.com/gists/%s", gistid))
	if err != nil {
		return col, fmt.Errorf("error fetching from github API: %s", err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return col, fmt.Errorf("api returned non-200 status: %d", res.StatusCode)
	}

	resp := struct {
		Description string
		Files       map[string]*struct {
			Size      int
			Raw_url   string
			Type      string
			Language  string
			Truncated bool
			Content   string
		}
	}{}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return col, err
	}

	if info := resp.Files[collectionInfoFilename]; info != nil {
		infop := struct {
			Name, Url, Description string
		}{}
		if err := json.Unmarshal([]byte(info.Content), &infop); err != nil {
			return nil, fmt.Errorf("error unmarshalling gist %s: %s", collectionInfoFilename, err.Error())
		}

		col.Title = infop.Name
		col.Url = infop.Url
		col.Description = infop.Description

	} else {
		col.Title = resp.Description
	}

	if err := EnsureCollection(store, col); err != nil {
		return col, err
	}

	items := []*core.CollectionItem{}

	if urls := resp.Files[urlsFilename]; urls != nil {
		s := bufio.NewScanner(strings.NewReader(resp.Files[urlsFilename].Content))
		for s.Scan() {
			if s.Err() != nil {
				if s.Err().Error() == "EOF" {
					break
				}
				return col, fmt.Errorf("error reading urls list: %s", err.Error())
			}

			u := &core.Url{Url: s.Text()}
			if err := u.Read(store); err != nil && err != datastore.ErrNotFound {
				return nil, fmt.Errorf("error reading url from datastore: %s, url: %s", err.Error(), s.Text())
			} else if err == datastore.ErrNotFound {
				if err = u.Save(store); err != nil {
					return col, fmt.Errorf("error saving url: %s, url: %s", err.Error())
				}
			}

			items = append(items, &core.CollectionItem{Url: *u})
		}

		if err := col.SaveItems(store, items); err != nil {
			return col, fmt.Errorf("error saving collection items: %s", err.Error())
		}
		return col, nil
	}

	return nil, fmt.Errorf("gist doesn't have a %s file", urlsFilename)
}

// EnsureCollection makes sure a collection exists, and if it doesn't creates one
// TODO - this is a candidate for moving into the core package
func EnsureCollection(store datastore.Datastore, col *core.Collection) error {
	title := col.Title
	description := col.Description
	url := col.Url
	creator := col.Creator

	if col.Url != "" {
		if err := col.Read(store); err != nil && err != datastore.ErrNotFound {
			return fmt.Errorf("Error reading collection: %s", err.Error())
		} else if err == datastore.ErrNotFound {
			col.Title = title
			col.Description = description
			col.Url = url
			col.Creator = creator

			if err = col.Save(store); err != nil {
				return fmt.Errorf("Error saving colleciton: %s", err.Error())
			}
		}
	}
	return nil
}
