package cassandredis

import "fmt"

// Query is a CQL statement with its args
type Query struct {
	Stmt string
	Args []interface{}
}

func (q *Query) String() string {
	return fmt.Sprintf("stmt: %s args: %v", q.Stmt, q.Args)
}

type Queries []*Query
