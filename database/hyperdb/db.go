package hyperdb

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"os"
	"sync"
	"time"
)

// ErrInvalidConfigType tells you that you specified a configuration type which is
// invalid. Helpful, I know. At the time of this writing all that is supported
// is json, and so you probably intend (or have) to use that.
var ErrInvalidConfigType = errors.New("Invalid configuration type")

func init() {
	rand.Seed(time.Now().UnixNano())
}

// DB consumes Server structs for configuration, and uses the contained
// information to provide a lookup mechanism for getting *sql.DB pointers that
// allows for multi datacenter, weighted read, and write server support
type DB struct {
	CacheReadServerFor time.Duration
	m                  sync.Mutex
	ds                 datasetList
	datacenter         string
}

// AddServer tells HyperDB to consume a server and add it to the lookup lists
func (d *DB) AddServer(server Server) {
	d.m.Lock()
	defer d.m.Unlock()
	server.hyperdb = d
	ds := d.dataset(server.Dataset)
	ds.addServer(server)
}

// Look up a read server, and hand back a pointer to a query-able interface.
func (d *DB) Read(dataset string) *sql.DB {
	d.m.Lock()
	defer d.m.Unlock()
	var s *Server
	ds := d.dataset(dataset)
	if ds.lastReadServer != nil {
		if time.Now().Sub(ds.lastReadServerCached) > d.CacheReadServerFor {
			ds.lastReadServer = nil
		} else {
			if err := ds.lastReadServer.db.Ping(); err == nil {
				return ds.lastReadServer.db
			} else {
				ds.lastReadServer = nil
			}
		}
	}
	if sl, ok := ds.read[d.datacenter]; ok {
		s = sl.get()
	}
	if s == nil {
		s = ds.allReadServers.get()
	}
	if s != nil {
		ds.lastReadServerCached = time.Now()
		ds.lastReadServer = s
		return s.db
	}
	return nil
}

// Look up the write server, and hand back a pointer to a query-able interface
func (d *DB) Write(dataset string) *sql.DB {
	d.m.Lock()
	defer d.m.Unlock()
	ds := d.dataset(dataset)
	s := ds.write
	if s == nil {
		return nil
	}
	return s.getDB()
}

// dataset finds, or creates, and returns a dataset based on its name
func (d *DB) dataset(name string) *dataset {
	if ds, ok := d.ds[name]; ok {
		return ds
	}
	var ds = &dataset{
		read:           datacenterList{},
		allReadServers: serverList{},
	}
	d.ds[name] = ds
	return ds
}

// ReadFile accepts a path and optional kind in additional to the stantard New
// function and its datacenter parameter.  The path should be a path to a real
// file that is readable and contains the proper data in the proper format. The
// only "kind" that is supported right now is "json".
//
// For json the config file should be formatted as follows:
//  { ds: "users", dc: "dfw" [...] }
//  { ds: "users", dc: "dfw" [...] }
//	[...]
func ReadFile(datacenter, path string, kind ...string) (*DB, error) {
	if len(kind) == 0 {
		kind = []string{"json"}
	}
	switch kind[0] {
	case "json", "js":
		fp, err := os.Open(path)
		defer fp.Close()
		if err != nil {
			return nil, err
		}
		var db = New(datacenter)
		dec := json.NewDecoder(fp)
		for {
			var server = Server{}
			if err := dec.Decode(&server); err != nil {
				if err == io.EOF {
					return db, nil
				}
				return nil, err
			}
			db.AddServer(server)
		}
	default:
		return nil, ErrInvalidConfigType
	}
}

// New returns a new HyperDB object
func New(datacenter string) *DB {
	return &DB{
		datacenter: datacenter,
		ds:         datasetList{},
	}
}
