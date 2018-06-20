package db

import "strconv"

// MongoDB ...
type MongoDB struct {
	table string
}

// Table set a table name to perform action
func (m *MongoDB) Table(name string) {
	m.table = name
}

// Query perform your raw query
func (m MongoDB) Query(q string) string {
	return "MongoDB>> Query on -> " + m.table + " " + q
}

// Find find a record using it's primary id
func (m MongoDB) Find(id int) string {
	return "MongoDB>> SELECT * FROM " + m.table + " WHERE id=" + strconv.Itoa(id)
}

// Get fetch all the records from the table
func (m MongoDB) Get() string {
	return "MongoDB>> SELECT * FROM " + m.table
}
