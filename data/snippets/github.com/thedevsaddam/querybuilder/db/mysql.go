package db

import "strconv"

// Mysql ...
type Mysql struct {
	table string
}

// Table set a table name to perform action
func (m *Mysql) Table(name string) {
	m.table = name
}

// Query perform your raw query
func (m Mysql) Query(q string) string {
	return "Mysql>> Query on -> " + m.table + " " + q
}

// Find find a record using it's primary id
func (m Mysql) Find(id int) string {
	return "Mysql>> SELECT * FROM " + m.table + " WHERE id=" + strconv.Itoa(id)
}

// Get fetch all the records from the table
func (m Mysql) Get() string {
	return "Mysql>> SELECT * FROM " + m.table
}
