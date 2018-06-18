package databases

type ResultRows interface {
	Next() bool
	StructScan(dest interface{}) error
}
