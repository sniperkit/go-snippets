package store

import (
	"strings"

	"github.com/acsellers/inflections"
)

type SQLConfig interface {
	SQLTable(string) string
	SQLColumn(string, string) string
}

type NameMap map[string]string

type AppConfig struct {
	SpecialTables  NameMap
	SpecialColumns map[string]NameMap

	Normal SQLConfig
}

func NewAppConfig(driverName string) AppConfig {
	return AppConfig{
		SpecialTables:  NameMap{},
		SpecialColumns: map[string]NameMap{},
		Normal:         LowerConfig{},
	}
}

func (c AppConfig) SQLTable(table string) string {
	if name, ok := c.SpecialTables[table]; ok {
		return name
	}
	if c.Normal == nil {
		return table
	}
	return c.Normal.SQLTable(table)
}

func (c AppConfig) SQLColumn(table, column string) string {
	if cols, ok := c.SpecialColumns[table]; ok {
		if name, ok := cols[column]; ok {
			return name
		}
	}
	if c.Normal == nil {
		return column
	}
	return c.Normal.SQLColumn(table, column)
}

type LowerConfig struct{}

func (LowerConfig) SQLTable(table string) string {
	return strings.ToLower(table)
}

func (LowerConfig) SQLColumn(table, column string) string {
	return strings.ToLower(column)
}

type PrefixConfig struct {
	TablePrefix  string
	ColumnPrefix string
}

func (pc PrefixConfig) SQLTable(table string) string {
	return pc.TablePrefix + table
}

func (pc PrefixConfig) SQLColumn(table, column string) string {
	return pc.ColumnPrefix + column
}

type RailsConfig struct{}

func (RailsConfig) SQLTable(table string) string {
	return strings.ToLower(inflections.Pluralize(inflections.Underscore(table)))
}

func (RailsConfig) SQLColumn(table, column string) string {
	return strings.ToLower(inflections.Underscore(column))
}
