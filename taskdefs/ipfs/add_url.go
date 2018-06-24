package ipfs

import (
	"fmt"
	"github.com/datatogether/core"
	"github.com/datatogether/sql_datastore"
	"github.com/datatogether/task_mgmt/tasks"
	"github.com/ipfs/go-datastore"
)

type TaskAdd struct {
	Url              string              `json:"url"`              // url to resource to be added
	Checksum         string              `json:"checksum"`         // optional checksum to check resp against
	ipfsApiServerUrl string              `json:"ipfsApiServerUrl"` // url of IPFS api server
	store            datastore.Datastore // internal datastore pointer
}

func NewTaskAdd() tasks.Taskable {
	return &TaskAdd{
		ipfsApiServerUrl: IpfsApiServerUrl,
	}
}

// AddCollection task needs to talk to an underlying database
// it's expected that the task executor will call this method
// before calling Do
func (t *TaskAdd) SetDatastore(store datastore.Datastore) {
	if sqlds, ok := store.(*sql_datastore.Datastore); ok {
		// if we're passed an sql datastore
		// make sure our collection model is registered
		sqlds.Register(
			&core.Url{},
		)
	}

	t.store = store
}

func (t *TaskAdd) Valid() error {
	if t.Url == "" {
		return fmt.Errorf("url param is required")
	}
	if t.ipfsApiServerUrl == "" {
		return fmt.Errorf("no ipfs server url provided, please configure the ipfs tasks package")
	}
	return nil
}

func (t *TaskAdd) Do(pch chan tasks.Progress) {
	p := tasks.Progress{Step: 1, Steps: 4, Status: "fetching resource"}

	u := &core.Url{
		Url: t.Url,
	}

	if err := u.Read(t.store); err != nil && err != datastore.ErrNotFound {
		p.Error = fmt.Errorf("Error reading url: %s", err.Error())
		pch <- p
		return
	}

	// TODO - unify these to use the same response from a given URL
	done := make(chan int, 0)
	go func() {
		if _, _, err := u.Get(t.store); err != nil {
			fmt.Printf("error getting url: %s\n", err.Error())
		}

		done <- 0
	}()
	go func() {
		_, _, err := ArchiveUrl(t.store, t.ipfsApiServerUrl, u)
		if err != nil {
			p.Error = err
			pch <- p
			return
		}

		done <- 0
	}()

	<-done
	<-done

	if err := u.Save(t.store); err != nil {
		p.Error = fmt.Errorf("error saving url: %s", err.Error())
		pch <- p
		return
	}

	p.Percent = 1.0
	p.Done = true
	p.Dest = fmt.Sprintf("/content/%s", u.Hash)
	pch <- p
	return
}
