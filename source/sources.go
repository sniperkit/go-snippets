package source

import (
	"database/sql"
	"fmt"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
)

// TODO - orderby is currently being ignored, need to fix with support for sorts in datastore queries
func ListSources(store datastore.Datastore, orderby string, limit, offset int) ([]*Source, error) {
	q := query.Query{
		Prefix: fmt.Sprintf("/%s", Source{}.DatastoreType()),
		Limit:  limit,
		Offset: offset,
	}

	res, err := store.Query(q)
	if err != nil {
		return nil, err
	}

	sources := make([]*Source, limit)
	i := 0
	for r := range res.Next() {
		if r.Error != nil {
			return nil, err
		}

		c, ok := r.Value.(*Source)
		if !ok {
			return nil, fmt.Errorf("Invalid Response")
		}

		sources[i] = c
		i++
	}

	return sources[:i], nil
}

func unmarshalSources(rows *sql.Rows, limit int) ([]*Source, error) {
	defer rows.Close()
	sources := make([]*Source, limit)
	i := 0
	for rows.Next() {
		s := &Source{}
		if err := s.UnmarshalSQL(rows); err != nil {
			return nil, err
		}

		sources[i] = s
		i++
	}

	sources = sources[:i]
	return sources, nil
}
