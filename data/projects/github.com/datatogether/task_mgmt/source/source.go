// This package is a carry-over from a previous incarnation of task_mgmt
// it should be cleaned up & removed in favour of types from the github.com/datatogether/archive package
package source

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/datatogether/sql_datastore"
	"github.com/datatogether/sqlutil"
	"github.com/ipfs/go-datastore"
	"github.com/pborman/uuid"
	"time"
)

// Source is an http origin of data
// TODO - this should be folded into core.Url, feels redundant
// TODO - ugh, archive already has a "Source" type that we're only
// able to get around b/c task_mgmt & archive run against different dbs
// need to destroy this asap
type Source struct {
	// version 4 uuid
	Id string `json:"id"`
	// created date rounded to secounds
	Created time.Time `json:"created"`
	// updated date rounded to secounds
	Updated time.Time `json:"updated"`
	// Human-readable title for source
	Title string `json:"title"`
	// Url to source data
	Url string `json:"url"`
	// Checksum of url
	Checksum string `json:"checksum"`
	// any associated metadata
	Meta map[string]interface{} `json:"meta"`
}

func (s Source) DatastoreType() string {
	return "Source"
}

func (s *Source) GetId() string {
	return s.Id
}

func (s *Source) Key() datastore.Key {
	return datastore.NewKey(fmt.Sprintf("%s:%s", s.DatastoreType(), s.GetId()))
}

func (s *Source) Read(store datastore.Datastore) error {
	so, err := store.Get(s.Key())
	if err != nil {
		if err == datastore.ErrNotFound {
			return datastore.ErrNotFound
		}
		return err
	}

	got, ok := so.(*Source)
	if !ok {
		return fmt.Errorf("Invalid Response")
	}
	*s = *got
	return nil
}

func (s *Source) Save(store datastore.Datastore) (err error) {
	var exists bool

	if s.Id != "" {
		exists, err = store.Has(s.Key())
		if err != nil {
			return err
		}
	}

	if !exists {
		s.Id = uuid.New()
		s.Created = time.Now().In(time.UTC).Round(time.Millisecond)
		s.Updated = s.Created
	} else {
		s.Updated = time.Now().In(time.UTC).Round(time.Millisecond)
	}

	return store.Put(s.Key(), s)
}

func (s *Source) Delete(store datastore.Datastore) error {
	return store.Delete(s.Key())
}

func (s *Source) NewSQLModel(key datastore.Key) sql_datastore.Model {
	return &Source{Id: key.Name()}
}

func (s *Source) SQLQuery(cmd sql_datastore.Cmd) string {
	switch cmd {
	case sql_datastore.CmdCreateTable:
		return qSourceCreateTable
	case sql_datastore.CmdSelectOne:
		return qSourceReadById
	case sql_datastore.CmdExistsOne:
		return qSourceExists
	case sql_datastore.CmdInsertOne:
		return qSourceInsert
	case sql_datastore.CmdDeleteOne:
		return qSourceDelete
	case sql_datastore.CmdUpdateOne:
		return qSourceUpdate
	case sql_datastore.CmdList:
		return qSourcesList
	default:
		// returning empty string will fire an error which
		// we want if asked for a Cmd we don't understand
		return ""
	}
}

func (s *Source) UnmarshalSQL(row sqlutil.Scannable) error {
	var (
		id, url, checksum, title string
		created, updated         time.Time
		meta                     []byte
	)

	if err := row.Scan(&id, &created, &updated, &title, &url, &checksum, &meta); err != nil {
		if err == sql.ErrNoRows {
			return datastore.ErrNotFound
		}
		return err
	}

	*s = Source{
		Id:       id,
		Created:  created,
		Updated:  updated,
		Title:    title,
		Url:      url,
		Checksum: checksum,
	}

	if meta != nil {
		s.Meta = map[string]interface{}{}
		if err := json.Unmarshal(meta, &s.Meta); err != nil {
			return err
		}
	}

	return nil
}

func (s *Source) SQLParams(cmd sql_datastore.Cmd) []interface{} {
	var (
		meta []byte
		err  error
	)
	switch cmd {
	case sql_datastore.CmdSelectOne, sql_datastore.CmdExistsOne, sql_datastore.CmdDeleteOne:
		return []interface{}{s.Id}
	case sql_datastore.CmdList:
		return []interface{}{}
	default:
		if s.Meta != nil {
			if meta, err = json.Marshal(s.Meta); err != nil {
				meta = nil
			}
		}

		return []interface{}{
			s.Id,
			s.Created,
			s.Updated,
			s.Title,
			s.Url,
			s.Checksum,
			meta,
		}
	}
}
