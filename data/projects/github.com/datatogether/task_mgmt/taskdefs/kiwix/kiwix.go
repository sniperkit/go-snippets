package kiwix

import (
	"fmt"
	"github.com/datatogether/sql_datastore"
	"github.com/datatogether/task_mgmt/source"
	"github.com/datatogether/task_mgmt/tasks"
	"github.com/ipfs/go-datastore"
)

type TaskUpdateSources struct {
	// internal datastore pointer
	store datastore.Datastore
}

func NewTaskUpdateSources() tasks.Taskable {
	return &TaskUpdateSources{}
}

func (t *TaskUpdateSources) Valid() error {
	return nil
}

// TaskUpdateSources task needs to talk to an underlying database
// it's expected that the task executor will call this method
// before calling Do
func (t *TaskUpdateSources) SetDatastore(store datastore.Datastore) {
	if sqlds, ok := store.(sql_datastore.Datastore); ok {
		// if we're passed an sql datastore
		// make sure our "source" model is registered
		sqlds.Register(&source.Source{})
	}

	t.store = store
}

// Do performs the task
func (t *TaskUpdateSources) Do(updates chan tasks.Progress) {
	p := tasks.Progress{Percent: 0.0, Step: 1, Steps: 2, Status: "fetching zims list"}
	// make sure we have a database connection
	if t.store == nil {
		p.Error = fmt.Errorf("datastore is required")
		updates <- p
		return
	}

	updates <- p

	zims, err := FetchZims()
	if err != nil {
		p.Error = fmt.Errorf("Error fetching zims: %s", err.Error())
		updates <- p
		return
	}

	sources, err := source.ListSources(t.store, "created DESC", 1000, 0)
	if err != nil {
		p.Error = fmt.Errorf("error listing sources: %s", err.Error())
		updates <- p
		return
	}

	p.Status = "updating sources"
	p.Percent = 0.50
	p.Step++
	updates <- p
	for _, s := range sources {
		for _, z := range zims {
			if s.Url == z.Url {
				if err := z.FetchMd5(); err != nil {
					p.Error = fmt.Errorf("error fetching MD5 checksum for source '%s': %s", s.Url, err.Error())
					updates <- p
					return
				}

				if z.Md5 != s.Checksum {
					s.Title = z.Title()
					s.Checksum = z.Md5
					if err := s.Save(t.store); err != nil {
						p.Error = fmt.Errorf("error saving source '%s': %s", s.Url, err.Error())
						updates <- p
						return
					}
				}
			}
		}
	}

	p.Done = true
	p.Percent = 1.0
	updates <- p
}
