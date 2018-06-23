package sqlimporter

import "strings"

var mysqlDialect = map[string]string{
	"create": "CREATE DATABASE ",
	"drop":   "DROP DATABASE ",
}

var postgresDialect = map[string]string{
	"create": "CREATE SCHEMA ",
	"drop":   "DROP SCHEMA ",
}

func getDialect(driver, process string) string {
	switch strings.ToLower(driver) {
	case "mysql":
		return mysqlDialect[process]
	case "postgres":
		return postgresDialect[process]
	default:
		return mysqlDialect[process]
	}
}
