package elasticsearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	elastigo "github.com/mattbaird/elastigo/lib"
	"gopkg.in/errgo.v1"
)

var (
	AlreadyIndexedErr = errors.New("index of is already indexed")
	searchHostQuery   = `{"query": {"filtered": {"filter": {"term": {"host": "%s"}}}}}`
	countItemsQuery   = `{"query": {"match": {"index_of_id": "%s"}}}`
)

type IndexOf struct {
	Id        string    `json:"_id,omitempty"`
	Host      string    `json:"host"`
	Scheme    string    `json:"scheme"`
	Path      string    `json:"path"`
	Count     int       `json:"count"`
	CrawledAt time.Time `json:"crawled_at"`
}

func (i *IndexOf) Type() string {
	return "index_of"
}

func (i *IndexOf) GetId() string {
	return i.Id
}

func (i *IndexOf) SetId(id string) {
	i.Id = id
}

func (i *IndexOf) URL() string {
	return i.Scheme + "://" + i.Host + i.Path
}

func FindIndexOf(id string) (*IndexOf, error) {
	res, err := defaultConn.Get(defaultIndex, "index_of", id, nil)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	var indexOf *IndexOf
	err = json.Unmarshal(*res.Source, &indexOf)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	indexOf.Id = id
	return indexOf, nil
}

func (i *IndexOf) Index() error {
	var searchParams map[string]interface{}
	query := fmt.Sprintf(searchHostQuery, i.Host)
	err := json.Unmarshal([]byte(query), &searchParams)
	if err != nil {
		return errgo.Mask(err)
	}
	res, err := defaultConn.Search(defaultIndex, i.Type(), nil, searchParams)
	if err != nil {
		if err == elastigo.RecordNotFound {
			return Index(i)
		}
		return errgo.Mask(err)
	}
	if res.Hits.Len() > 0 {
		return AlreadyIndexedErr
	}
	return Index(i)
}

func ListIndexOf() ([]Document, error) {
	docs, err := ListDocuments((*IndexOf)(nil))
	if err != nil {
		return nil, errgo.Mask(err)
	}
	wg := &sync.WaitGroup{}
	for _, doc := range docs {
		wg.Add(1)
		go func(doc *IndexOf) {
			defer wg.Done()
			c, err := doc.CountItems()
			if err != nil {
				log.Println(err)
			}
			doc.Count = c
		}(doc.(*IndexOf))
	}
	wg.Wait()
	return docs, nil
}

func (i *IndexOf) CountItems() (int, error) {
	var countParams map[string]interface{}
	err := json.Unmarshal([]byte(fmt.Sprintf(countItemsQuery, i.Id)), &countParams)
	if err != nil {
		return -1, errgo.Mask(err)
	}
	res, err := defaultConn.Count(defaultIndex, "index_item", nil, countParams)
	return res.Count, nil
}
