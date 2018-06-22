package gorgonzola

import (
	"encoding/json"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
)

// Datastore implements the Storage interface targeting GAE Datastore
type Datastore struct{}

// NewDatastore creates new Datastore
func NewDatastore() *Datastore {
	return &Datastore{}
}

func (ds *Datastore) saveJob(c appengine.Context, job *Job) error {
	key := datastore.NewKey(c, "Job", job.Hash, 0, nil)
	_, err := datastore.Put(c, key, job)
	return err
}

func (ds *Datastore) saveJSONJobs(c appengine.Context, key string, jj *JSONJobs) error {
	for _, job := range jj.getJobs(key) {
		job.Hash = job.getHash()
		if err := ds.saveJob(c, job); err != nil {
			return err
		}
	}
	return nil
}

func (ds *Datastore) disableOldJobs(c appengine.Context, key *datastore.Key) error {
	q := datastore.NewQuery("Job").Filter("LinkKey = ", key.Encode())
	var jobs []Job
	if _, err := q.GetAll(c, &jobs); err != nil {
		return err
	}
	for _, job := range jobs {
		job.Active = false
		if err := ds.saveJob(c, &job); err != nil {
			return err
		}
	}
	return nil
}

func (ds *Datastore) updateJobs(c appengine.Context, keyString string) error {
	key, err := datastore.DecodeKey(keyString)
	if err != nil {
		return err
	}
	var link Link
	if err := datastore.Get(c, key, &link); err != nil {
		return err
	}
	if err := ds.disableOldJobs(c, key); err != nil {
		return err
	}
	var jjraw []byte
	if jjraw, err = getJSONJobsDoc(c, link.URL); err != nil {
		return err
	}
	if err := validateDoc(string(jjraw)); err != nil {
		return err
	}
	var jj JSONJobs
	json.Unmarshal(jjraw, &jj)
	if err := ds.saveJSONJobs(c, key.Encode(), &jj); err != nil {
		return err
	}
	link.Fetched = time.Now()
	if _, err := datastore.Put(c, key, &link); err != nil {
		return err
	}
	return nil
}

// AddURL adds new Json-job url to the job board
func (ds *Datastore) AddURL(r *http.Request, url string) error {
	c := appengine.NewContext(r)
	link := &Link{
		URL:     url,
		Created: time.Now(),
	}
	key := datastore.NewKey(c, "Link", url, 0, nil)
	if _, err := datastore.Put(c, key, link); err != nil {
		return err
	}
	laterFunc.Call(c, key.Encode())
	return nil
}

// GetJobs returns maximum `limit` jobs ordered by creation date
func (ds *Datastore) GetJobs(r *http.Request, limit int) ([]Job, error) {
	var jobs []Job
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Job").Order("-Created").Limit(limit)
	_, err := q.GetAll(c, &jobs)
	return jobs, err
}

// GetJob returns single job by hash value
func (ds *Datastore) GetJob(r *http.Request, hash string) (*Job, error) {
	var job Job
	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "Job", hash, 0, nil)
	err := datastore.Get(c, key, &job)
	if err == datastore.ErrNoSuchEntity {
		return nil, HTTPError{
			err,
			"Job not found",
			http.StatusNotFound,
		}
	}
	return &job, err
}

// Update updates the job offers for single Link
func (ds *Datastore) Update(r *http.Request) error {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Link").Order("-Fetched").Limit(1).KeysOnly()
	keys, err := q.GetAll(c, nil)
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		ds.updateJobs(c, keys[0].Encode())
	}
	return nil
}
