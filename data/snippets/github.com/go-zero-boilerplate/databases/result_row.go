package databases

type ResultRow interface {
	Scan(dest ...interface{}) error
}
