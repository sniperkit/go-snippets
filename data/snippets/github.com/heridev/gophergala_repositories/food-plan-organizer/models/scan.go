package models

type StructScanner interface {
	Scan(dest ...interface{}) error
}
