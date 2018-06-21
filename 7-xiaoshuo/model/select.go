package model

import (
	"fmt"
	"strconv"
	"strings"
)

type Select struct {
	table               string
	where, order        []string
	limit, offset       uint
	setLimit, setOffset bool
	columns             map[string]string
	params              []interface{}
}

func (t *Select) From(from string) *Select {
	t.table = from
	return t
}

func (t *Select) Where(where string, params ...interface{}) *Select {
	for _, v := range params {
		m := ""
		switch v.(type) {
		case string:
			where = strings.Replace(where, "?", fmt.Sprintf(`'%s'`, v.(string)), 1)
			m = v.(string)
		case int, int8, int16, int32, int64:
			m = fmt.Sprintf("%d", v)
		case uint, uint8, uint16, uint32, uint64:
			m = fmt.Sprintf("%d", v)
		case float32, float64:
			m = fmt.Sprintf("%f", v)
		case []int:
			var values []string
			for _, n := range v.([]int) {
				values = append(values, strconv.Itoa(n))
			}
			m = strings.Join(values, ", ")
		case []string:
			var values []string
			for _, s := range v.([]string) {
				values = append(values, fmt.Sprintf(`"%s"`, s))
			}
			m = strings.Join(values, ", ")
		}
		if m != "" {
			where = strings.Replace(where, "?", m, 1)
		}
	}
	t.where = append(t.where, where)
	return t
}

func (t *Select) FilterWhere(where string, params ...interface{}) *Select {
	if len(params) != 0 {
		t.Where(where, params...)
	}
	return t
}

func (t *Select) Limit(limit uint) *Select {
	t.limit = limit
	t.setLimit = true
	return t
}

func (t *Select) Offset(offset uint) *Select {
	t.offset = offset
	t.setOffset = true
	return t
}
