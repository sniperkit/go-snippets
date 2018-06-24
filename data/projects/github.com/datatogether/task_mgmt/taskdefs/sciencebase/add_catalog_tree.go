package sciencebase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/datatogether/cdxj"
	"github.com/datatogether/core"
	sb "github.com/datatogether/linked_data/sciencebase"
	"github.com/datatogether/sql_datastore"
	"github.com/datatogether/task_mgmt/taskdefs/ipfs"
	"github.com/datatogether/task_mgmt/tasks"
	"github.com/ipfs/go-datastore"
	"time"
)

var IpfsApiServerUrl = ""
var count = 0

const defaultCrawlDelay = time.Second / 2
const pageSize = 100

// AddCatalogTree injests a collection to IPFS,
// it iterates through each setting hashes on collection urls
// and, eventually, generates a cdxj index of the archive
type AddCatalogTree struct {
	// title for collection if no collection present
	CollectionTitle string
	// url that points to catalog
	Url string `json:"url"`
	// how many items deep to crawl at most, -1 == no max
	MaxDepth int
	// how many fetching goroutines to spin up. max 5
	Parallelism int
	// skip items that already have a hash value
	SkipArchived bool
	// how long to sleep between requests in seconds(inside of parallel routines)
	CrawDelay time.Duration
	// url of IPFS api server, should be set internally
	ipfsApiServerUrl string
	// internal datastore pointer
	store datastore.Datastore
}

func NewAddCatalogTree() tasks.Taskable {
	return &AddCatalogTree{
		CrawDelay:        defaultCrawlDelay,
		Parallelism:      2,
		ipfsApiServerUrl: IpfsApiServerUrl,
	}
}

// AddCatalogTree task needs to talk to an underlying database
// it's expected that the task executor will call this method
// before calling Do
func (t *AddCatalogTree) SetDatastore(store datastore.Datastore) {
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

func (t *AddCatalogTree) Valid() error {
	if t.Url == "" {
		return fmt.Errorf("url is required")
	}
	if t.Parallelism > 20 {
		t.Parallelism = 20
	}
	// super rough coercion to seconds if non-default crawldelay
	// is specified
	if t.CrawDelay != defaultCrawlDelay {
		if t.CrawDelay > 10 {
			t.CrawDelay = 10
		}
		t.CrawDelay = time.Second * t.CrawDelay
	}
	if t.ipfsApiServerUrl == "" {
		return fmt.Errorf("no ipfs server url provided, please configure the ipfs tasks package")
	}
	return nil
}

func (t *AddCatalogTree) Do(pch chan tasks.Progress) {
	p := tasks.Progress{Step: 1, Steps: 4, Status: "loading collection"}
	pch <- p

	collection := &core.Collection{Url: t.Url}
	if err := collection.Read(t.store); err != nil && err != datastore.ErrNotFound {
		p.Error = fmt.Errorf("Error reading collection: %s", err.Error())
		pch <- p
		return
	} else if err == datastore.ErrNotFound {
		collection.Title = t.CollectionTitle
		if err = collection.Save(t.store); err != nil {
			p.Error = fmt.Errorf("Error saving colleciton: %s", err.Error())
			pch <- p
			return
		}
	}
	p.Step++

	indexBuf := bytes.NewBuffer(nil)
	index := cdxj.NewWriter(indexBuf)

	// TODO - refactor done chan to report progress, possibly sending the number
	// of indexes *remaining* with each iteration

	if err := ArchiveCatalog(t.store, t.ipfsApiServerUrl, collection, index, t.Url, t.MaxDepth, t.Parallelism); err != nil {
		fmt.Println(err.Error())
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
	indexhash, err := ipfs.WriteToIpfs(t.ipfsApiServerUrl, fmt.Sprintf("%s.cdxj", collection.Id), indexBuf.Bytes())
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

func ArchiveCatalog(store datastore.Datastore, ipfsApiUrl string, col *core.Collection, index *cdxj.Writer, rootUrl string, maxDepth, parallelism int) error {
	visit := make(chan childItem, 100)
	visited := make(chan childItem, 100)
	tracks := make([]chan childItem, parallelism)
	wait := make(chan bool)
	start := time.Now()

	for i, _ := range tracks {
		tracks[i] = make(chan childItem, 100)
		go func(track, visit, visited chan childItem) {
			for child := range track {
				if err := ArchiveChild(store, ipfsApiUrl, col, index, child, visit, visited); err != nil {
					fmt.Println("error archiving url", err.Error())
					// TODO - collect errored urls, or flag as errored?
				}
			}
		}(tracks[i], visit, visited)
	}

	go func() {
		t := 0
		stack := 0
		for {
			select {
			case child, ok := <-visit:
				fmt.Println("visit", child.url)
				if ok && maxDepth == -1 || child.depth < maxDepth {
					tracks[t] <- child
					t++
					if t == len(tracks) {
						t = 0
					}
					stack++
				}
			case <-visited:
				stack--
				if stack <= 0 {
					wait <- false
				}
			}
		}
	}()

	// add root item to kick off
	visit <- childItem{0, rootUrl}

	<-wait
	fmt.Println(time.Since(start), count, "nodes")
	return nil
}

type childItem struct {
	depth int
	url   string
}

func ArchiveChild(store datastore.Datastore, ipfsApiUrl string, collection *core.Collection, index *cdxj.Writer, child childItem, visit, visited chan childItem) error {
	// core
	u := &core.Url{Url: child.url}
	hh, bh, err := ipfs.ArchiveUrl(store, ipfsApiUrl, u)
	if err != nil {
		return err
	}
	count++

	body, err := ipfs.ReadFile(ipfsApiUrl, bh)
	if err != nil {
		fmt.Println("error getting ipfs json body:", err.Error())
		return err
	}

	item := &sb.Item{}
	if err := json.NewDecoder(body).Decode(item); err != nil {
		fmt.Println("error decoding child json", u.Url, err.Error())
		return err
	}
	body.Close()

	// grab any children in a goroutine
	if item.HasChildren {
		u := &core.Url{Url: item.ChildrenJsonUrl()}
		hh, bh, err := ipfs.ArchiveUrl(store, ipfsApiUrl, u)
		if err != nil {
			fmt.Println("error archiving children catalog url: ", err.Error())
			return err
		}

		body, err := ipfs.ReadFile(ipfsApiUrl, bh)
		if err != nil {
			fmt.Println("error getting ipfs json body:", err.Error())
			return err
		}

		c := &sb.Catalog{}
		if err := json.NewDecoder(body).Decode(c); err != nil {
			fmt.Printf("error decoding chilren json: %s", err.Error())
			return err
		}
		body.Close()

		// iterate through "items", forward through channel
		for _, item := range c.Items {
			if item.Link != nil {
				visit <- childItem{child.depth + 1, item.Link.JsonUrl()}
			}
		}

		if err = collection.SaveItems(store, []*core.CollectionItem{
			&core.CollectionItem{Url: *u},
		}); err != nil {
			// p.Error = fmt.Errorf("error saving dataset %d dist %d to collection: %s", i, j, err.Error())
			// pch <- p
			return err
		}

		if u.LastGet == nil {
			now := time.Now()
			u.LastGet = &now
		}

		// TODO - demo content for now, this is going to need lots of refinement
		indexRec := &cdxj.Record{
			Uri:        u.Url,
			Timestamp:  *u.LastGet,
			RecordType: "", // TODO set record type?
			JSON: map[string]interface{}{
				"locator": fmt.Sprintf("urn:ipfs/%s/%s", hh, bh),
			},
		}

		if err := index.Write(indexRec); err != nil {
			// p.Error = fmt.Errorf("Error writing %s cdxj index to ipfs: %s", filepath.Base(u.Url), err.Error())
			// pch <- p
			return err
		}
	}

	// TODO - actually archive files
	// if err := ArchiveFiles(col, index, item); err != nil {
	// }

	if u.LastGet == nil {
		now := time.Now()
		u.LastGet = &now
	}

	// TODO - demo content for now, this is going to need lots of refinement
	indexRec := &cdxj.Record{
		Uri:        u.Url,
		Timestamp:  *u.LastGet,
		RecordType: "", // TODO set record type?
		JSON: map[string]interface{}{
			"locator": fmt.Sprintf("urn:ipfs/%s/%s", hh, bh),
		},
	}

	if err := index.Write(indexRec); err != nil {
		// p.Error = fmt.Errorf("Error writing %s cdxj index to ipfs: %s", filepath.Base(u.Url), err.Error())
		// pch <- p
		return err
	}

	if err = collection.SaveItems(store, []*core.CollectionItem{
		&core.CollectionItem{Url: *u},
	}); err != nil {
		// p.Error = fmt.Errorf("error saving dataset %d dist %d to collection: %s", i, j, err.Error())
		// pch <- p
		return err
	}

	visited <- child
	return nil
}

// func ArchiveFiles(store datastore.Datastore, ipfsApiUrl string, col *core.Collection, index *cdxj.Writer, item *sb.Item) error {
// 	for _, f := range item.Files {
// 		// TODO - core file
// 		u := &core.Url{Url: f.Url}
// 		hh, bh, err := ipfs.ArchiveUrl(store, ipfsApiUrl, u)

// 		if u.LastGet == nil {
// 			now := time.Now()
// 			u.LastGet = &now
// 		}

// 		indexRec := &cdxj.Record{
// 			Uri:        u.Url,
// 			Timestamp:  *u.LastGet,
// 			RecordType: "", // TODO set record type?
// 			JSON: map[string]interface{}{
// 				"locator": fmt.Sprintf("urn:ipfs/%s/%s", hh, bh),
// 			},
// 		}
// 	}
// 	return nil
// }
