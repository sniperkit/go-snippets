package cassandredis

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

func (p *proxy) preLPUSH(req *Command, session *proxySession) {
	if len(req.Args) < 0 {
		session.lastErr = errors.New("syntax error on LPUSH")
		return
	}

	session.currentKey = newQueryKeyFromString(req.Args[0].(*respBulkStringValue).String())
	key := session.currentKey.metadataKey()

	if p.md.hasListIndex(key) {
		return
	}

	q := "SELECT list_index FROM lists_indexes WHERE table_name = ? ORDER BY list_index DESC LIMIT 1"
	it := p.session.Query(q, key).Iter()

	var maxIndex int64 = -1
	it.Scan(&maxIndex)

	if err := it.Close(); err != nil {
		session.lastErr = err
		return
	}

	if maxIndex >= 0 {
		p.md.setListIndex(key, maxIndex+1)
	}
}

func (p *proxy) processLPUSH(req *Command, session *proxySession) {
	queries, err := p.mapLPUSH(req, session)
	if err != nil {
		session.lastErr = err
		return
	}

	for _, query := range queries {
		log.Printf("executing DML query %s", query)

		q := p.session.Query(query.Stmt, query.Args...)
		if err := q.Exec(); err != nil {
			log.Println(err)
			sendErrorResponse(session, err)
			return
		}
	}

	sendIntegerResponse(session, session.scratch.(int64)+1)
}

func (p *proxy) ddlMapLPUSH(req *Command, session *proxySession) (Queries, error) {
	if len(req.Args) < 2 {
		return nil, errors.New("syntax error on LPUSH, need at least two arguments")
	}

	if len(session.currentKey.pk) > 0 {
		return ddlMapLPUSHWithPK(session.currentKey)
	}

	return Queries{&Query{
		Stmt: fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(bucket int, value blob, tuuid timeuuid, PRIMARY KEY ((bucket), tuuid))", session.currentKey.namespace),
	}}, nil
}

func ddlMapLPUSHWithPK(key queryKey) (Queries, error) {
	var buf bytes.Buffer
	buf.WriteString("CREATE TABLE IF NOT EXISTS ")
	buf.WriteString(key.namespace)
	buf.WriteRune('(')
	for i, _ := range key.pk {
		fmt.Fprintf(&buf, "f%d text, ", i)
	}
	buf.WriteString("value blob, tuuid timeuuid, PRIMARY KEY ((")
	for i, _ := range key.pk {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "f%d", i)
	}
	buf.WriteString("), tuuid));")

	return Queries{&Query{
		Stmt: buf.String(),
	}}, nil
}

func (p *proxy) mapLPUSH(req *Command, session *proxySession) (Queries, error) {
	var queries Queries

	values := []interface{}{req.Args[1].(*respBulkStringValue).value}

	if len(session.currentKey.pk) > 0 {
		return p.mapLPUSHWithPK(req, session, values)
	}

	uuid := gocql.TimeUUID()

	stmt := fmt.Sprintf("INSERT INTO %s(bucket, value, tuuid) VALUES(0, ?, ?)", session.currentKey.namespace)
	args := append(values, uuid)
	queries = append(queries, &Query{stmt, args})

	i := p.md.nextListIndex(session.currentKey.metadataKey())
	session.scratch = i

	stmt = "INSERT INTO lists_indexes(table_name, list_index, uuid_val) VALUES(?, ?, ?)"
	args = []interface{}{session.currentKey.metadataKey(), i, uuid}
	queries = append(queries, &Query{stmt, args})

	return queries, nil
}

func (p *proxy) mapLPUSHWithPK(req *Command, session *proxySession, values []interface{}) (Queries, error) {
	var buf bytes.Buffer

	buf.WriteString("INSERT INTO ")
	buf.WriteString(session.currentKey.namespace)
	buf.WriteRune('(')
	for i, _ := range session.currentKey.pk {
		fmt.Fprintf(&buf, "f%d, ", i)
	}
	buf.WriteString("value, tuuid) VALUES(")
	for i, v := range session.currentKey.pk {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteRune('\'')
		buf.WriteString(v)
		buf.WriteRune('\'')
	}
	buf.WriteString(", ?, ?)")

	return Queries{&Query{
		Stmt: buf.String(),
		Args: append(values, gocql.TimeUUID()),
	}}, nil
}
