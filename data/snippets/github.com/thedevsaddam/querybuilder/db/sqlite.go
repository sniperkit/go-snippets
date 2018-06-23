package db

import "strconv"

// Sqlite ...
type Sqlite struct {
	table string
}

// Table set a table name to perform action
func (s *Sqlite) Table(name string) {
	s.table = name
}

// Query perform your raw query
func (s Sqlite) Query(q string) string {
	return "Sqlite>> Query on -> " + s.table + " " + q
}

// Find find a record using it's primary id
func (s Sqlite) Find(id int) string {
	return "Sqlite>> SELECT * FROM " + s.table + " WHERE id=" + strconv.Itoa(id)
}

//Get fetch all the records from the table
func (s Sqlite) Get() string {
	return "Sqlite>> SELECT * FROM " + s.table
}
