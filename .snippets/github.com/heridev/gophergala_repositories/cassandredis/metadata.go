package cassandredis

import (
	"fmt"
	"sync"
)

func listsIndexesQuery() *Query {
	return &Query{
		Stmt: fmt.Sprintf("CREATE TABLE IF NOT EXISTS lists_indexes(table_name text, uuid_val uuid, list_index bigint, PRIMARY KEY ((table_name), list_index))"),
	}
}

// BootstrapMetadata takes care of creating metadata tables needed for Cassandredis
func (s *Server) BootstrapMetadata() error {
	q := listsIndexesQuery()

	if err := s.proxy.session.Query(q.Stmt, q.Args...).Exec(); err != nil {
		return err
	}

	return nil
}

type metadata struct {
	mu           sync.RWMutex
	listsIndexes map[string]int64
}

func newMetadata() *metadata {
	return &metadata{
		listsIndexes: make(map[string]int64),
	}
}

func (m *metadata) hasListIndex(key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.listsIndexes[key]
	return ok
}

func (m *metadata) setListIndex(key string, val int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.listsIndexes[key] = val
}

func (m *metadata) nextListIndex(key string) int64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	old, ok := m.listsIndexes[key]

	var newI int64
	if ok {
		newI = old + 1
	} else {
		newI = 0
	}

	m.listsIndexes[key] = newI

	return old
}
