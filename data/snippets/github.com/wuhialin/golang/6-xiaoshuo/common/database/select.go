package database

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var columnIndex = 1

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

func (t *Select) Columns(columns interface{}) *Select {
	if t.columns == nil {
		t.columns = make(map[string]string)
	}
	switch columns.(type) {
	case string:
		t.columns[strconv.Itoa(columnIndex)] = columns.(string)
		columnIndex++
	case map[string]string:
		if v, ok := columns.(map[string]string); ok {
			t.columns = v
		}
	}
	return t
}

func (t *Select) Order(field string, orderType ...string) *Select {
	if len(orderType) == 0 {
		orderType = append(orderType, "ASC")
	}
	t.order = append(t.order, strings.Join([]string{field, orderType[0]}, " "))
	return t
}

func (t *Select) renderColumns() string {
	var field string
	var columns []string
	if t.columns == nil {
		t.Columns("*")
	}
	exp, _ := regexp.Compile(`^[\d|\s]+$`)
	for k, v := range t.columns {
		field = v
		if k != "" && !exp.Match([]byte(k)) {
			field = strings.Join([]string{k, v}, " ")
		}
		columns = append(columns, field)
	}
	return strings.Join(columns, ", ")
}

func (t *Select) renderTable() string {
	tables := []string{"FROM", t.table}
	return strings.Join(tables, " ")
}

func (t *Select) renderWhere() string {
	if len(t.where) == 0 {
		return ""
	}
	var wheres []string
	for _, where := range t.where {
		wheres = append(wheres, fmt.Sprintf(`(%s)`, where))
	}
	return fmt.Sprintf(`WHERE %s`, strings.Join(wheres, " AND "))
}

func (t *Select) renderGroup() string {
	return ""
}

func (t *Select) renderHaving() string {
	return ""
}

func (t *Select) renderOrder() string {
	if len(t.order) == 0 {
		return ""
	}
	var orders []string
	for _, v := range t.order {
		orders = append(orders, v)
	}
	return fmt.Sprintf(`ORDER BY %s`, strings.Join(orders, " "))
}

func (t *Select) renderLimit() string {
	if t.setLimit {
		return fmt.Sprintf("LIMIT %d", t.limit)
	}
	return ""
}

func (t *Select) renderOffset() string {
	if t.setOffset {
		return fmt.Sprintf("OFFSET %d", t.offset)
	}
	return ""
}

func (t *Select) String() string {
	sqls := []string{"SELECT"}
	sqls = append(sqls, t.renderColumns())
	sqls = append(sqls, t.renderTable())

	where := t.renderWhere()
	if where != "" {
		sqls = append(sqls, where)
	}

	group := t.renderGroup()
	if group != "" {
		sqls = append(sqls, group)
	}

	having := t.renderHaving()
	if having != "" {
		sqls = append(sqls, having)
	}

	order := t.renderOrder()
	if order != "" {
		sqls = append(sqls, order)
	}

	limit := t.renderLimit()
	if limit != "" {
		sqls = append(sqls, limit)
	}

	offset := t.renderOffset()
	if offset != "" {
		sqls = append(sqls, offset)
	}
	return strings.Join(sqls, " ")
}
