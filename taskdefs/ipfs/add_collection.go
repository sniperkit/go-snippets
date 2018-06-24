package ipfs

import (
	"bytes"
	"fmt"
	"github.com/datatogether/cdxj"
	"github.com/datatogether/core"
	"github.com/datatogether/sql_datastore"
	"github.com/datatogether/task_mgmt/tasks"
	"github.com/ipfs/go-datastore"
	"path/filepath"
)

const pageSize = 100

// AddCollection injests a a collection to IPFS,
// it iterates through each setting hashes on collection urls
// and, eventually, generates a cdxj index of the archive
type AddCollection struct {
	CollectionId     string              `json:"collectionId"`     // url to resource to be added
	ipfsApiServerUrl string              `json:"ipfsApiServerUrl"` // url of IPFS api server
	store            datastore.Datastore // internal datastore pointer
}

func NewAddCollection() tasks.Taskable {
	return &AddCollection{
		ipfsApiServerUrl: IpfsApiServerUrl,
	}
}

// AddCollection task needs to talk to an underlying database
// it's expected that the task executor will call this method
// before calling Do
func (t *AddCollection) SetDatastore(store datastore.Datastore) {
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

func (t *AddCollection) Valid() error {
	if t.CollectionId == "" {
		return fmt.Errorf("collectionId is required")
	}
	if t.ipfsApiServerUrl == "" {
		return fmt.Errorf("no ipfs server url provided, please configure the ipfs tasks package")
	}
	return nil
}

func (t *AddCollection) Do(pch chan tasks.Progress) {
	p := tasks.Progress{Step: 1, Steps: 4, Status: "loading collection"}
	// 1. Get the Collection & Item Count
	pch <- p

	collection := &core.Collection{Id: t.CollectionId}
	if err := collection.Read(t.store); err != nil {
		p.Error = fmt.Errorf("Error reading collection: %s", err.Error())
		pch <- p
		return
	}
	p.Step++

	count, err := collection.ItemCount(t.store)
	if err != nil {
		p.Error = fmt.Errorf("Error reading collection item count: %s", err.Error())
		pch <- p
		return
	}

	pctAdd := 1.0 / float32(count)

	indexBuf := bytes.NewBuffer(nil)
	index := cdxj.NewWriter(indexBuf)

	pageCount := count / pageSize
	for i := 0; i <= pageCount; i++ {
		items, err := collection.ReadItems(t.store, "created DESC", pageSize, i*pageSize)
		if err != nil {
			p.Error = fmt.Errorf("Error reading items page index %d/%d: %s", i, pageCount, err.Error())
			pch <- p
			return
		}

		// TODO - parallelize a lil bit
		for j, item := range items {
			// TODO - parse this from schema
			urlstr := item.Url.Url

			p.Status = fmt.Sprintf("archiving url %d/%d", (i*pageSize)+j+1, count)
			p.Percent += pctAdd
			pch <- p

			// TODO - get the actual start time from header WARC Record
			// start := time.Now()
			headerHash, bodyHash, err := ArchiveUrl(t.store, t.ipfsApiServerUrl, &item.Url)
			if err != nil {
				p.Error = err
				pch <- p
				return
			}

			// TODO - demo content for now, this is going to need lots of refinement
			indexRec := &cdxj.Record{
				Uri:        urlstr,
				Timestamp:  *item.LastGet,
				RecordType: "", // TODO set record type?
				JSON: map[string]interface{}{
					"locator": fmt.Sprintf("urn:ipfs/%s/%s", headerHash, bodyHash),
				},
			}

			if err := index.Write(indexRec); err != nil {
				p.Error = fmt.Errorf("Error writing %s index record to ipfs: %s", filepath.Base(urlstr), err.Error())
				pch <- p
				return
			}

			// if err := item.Save(t.store); err != nil {
			// 	p.Error = fmt.Errorf("Error saving item: %s: %s", item.Id, err.Error())
			// 	pch <- p
			// 	return
			// }
		}
	}

	p.Step++
	p.Status = "writing index to IPFS"
	pch <- p
	// close & sort the index
	if err := index.Close(); err != nil {
		p.Error = fmt.Errorf("Error closing index %s", err.Error())
		pch <- p
		return
	}
	indexhash, err := WriteToIpfs(t.ipfsApiServerUrl, fmt.Sprintf("%s.cdxj", collection.Id), indexBuf.Bytes())
	if err != nil {
		p.Error = fmt.Errorf("Error writing index to ipfs: %s", err.Error())
		pch <- p
		return
	}
	fmt.Printf("collection %s index hash: %s\n", collection.Id, indexhash)

	p.Step++
	p.Status = "saving collection results"
	pch <- p
	if err := collection.Save(t.store); err != nil {
		p.Error = fmt.Errorf("Error saving collection: %s", err.Error())
		pch <- p
		return
	}

	p.Percent = 1.0
	p.Done = true
	pch <- p
	return
}
