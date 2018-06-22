package cassandredis

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"github.com/gocql/gocql"
)

func (p *proxy) preLRANGE(req *Command, session *proxySession) {
	if len(req.Args) < 3 {
		session.lastErr = errors.New("syntax error on LRANGE, need at least 3 arguments")
	}

	session.currentKey = newQueryKeyFromString(req.Args[0].(*respBulkStringValue).String())
}

func (p *proxy) processLRANGE(req *Command, session *proxySession) {
	start, stop, err := getRangeStartStop(req.Args[1:])
	if err != nil {
		session.lastErr = err
		return
	}

	first, last, err := p.getListsMatchingUUIDs(session.currentKey, start, stop)
	if err != nil {
		session.lastErr = err
		return
	}

	if len(session.currentKey.pk) > 0 {
		p.processLRANGEWithPk(req, session, first, last)
		return
	}

	q := makeLRANGEQuery(session.currentKey, first, last)
	p.readAndSendListValues(q, session)
}

func (p *proxy) getListsMatchingUUIDs(key queryKey, start, stop int64) (first gocql.UUID, last gocql.UUID, err error) {
	q := "SELECT uuid_val FROM lists_indexes WHERE table_name = ? AND list_index = ?"

	it := p.session.Query(q, key.metadataKey(), start).Iter()
	it.Scan(&first)
	if err = it.Close(); err != nil {
		return
	}

	it = p.session.Query(q, key.metadataKey(), stop).Iter()
	it.Scan(&last)
	err = it.Close()

	return
}

func (p *proxy) processLRANGEWithPk(req *Command, session *proxySession, first, last gocql.UUID) {
	q := makeLRANGEWithPkQuery(session.currentKey, first, last)
	p.readAndSendListValues(q, session)
}

func (p *proxy) readAndSendListValues(q *Query, session *proxySession) {
	it := p.session.Query(q.Stmt, q.Args...).Iter()

	var values [][]byte
	var value []byte
	for it.Scan(&value) {
		values = append(values, value)
	}

	if err := it.Close(); err != nil {
		session.lastErr = err
		return
	}

	sendArrayResponse(session, values)
}

func makeLRANGEQuery(key queryKey, first, last gocql.UUID) *Query {
	return &Query{
		Stmt: fmt.Sprintf("SELECT value FROM %s WHERE bucket = 0 AND tuuid >= ? AND tuuid <= ?", key.namespace),
		Args: []interface{}{first, last},
	}
}

func makeLRANGEWithPkQuery(key queryKey, first, last gocql.UUID) *Query {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "SELECT value FROM %s", key.namespace)
	for i, _ := range key.pk {
		if i == 0 {
			buf.WriteString(" WHERE ")
		} else {
			buf.WriteString(" AND ")
		}
		fmt.Fprintf(&buf, "f%d = ?", i)
	}
	buf.WriteString(" AND tuuid >= ? AND tuuid < ?")

	var args []interface{}
	for _, v := range key.pk {
		args = append(args, v)
	}
	args = append(args, first, last)

	return &Query{Stmt: buf.String(), Args: args}
}

func getRangeStartStop(args []interface{}) (int64, int64, error) {
	v1 := args[0].(*respBulkStringValue).String()
	v2 := args[1].(*respBulkStringValue).String()

	i1, err := strconv.ParseInt(v1, 10, 64)
	if err != nil {
		return -1, -1, err
	}

	i2, err := strconv.ParseInt(v2, 10, 64)
	if err != nil {
		return -1, -1, err
	}

	return i1, i2, nil
}
