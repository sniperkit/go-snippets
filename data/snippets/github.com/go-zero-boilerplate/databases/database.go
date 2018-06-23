package databases

type Database interface {
	IsErrNoRows(err error) bool

	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (ExecResult, error)
	QueryRow(query string, args ...interface{}) ResultRow
	Query(query string, args ...interface{}) (ResultRows, error)

	BeginTx() (Database, error)
	CommitTx() error
	RollbackTx() error

	//DeferredRollbackIfNotHandled must be used with the `defer` keyword to finish a transaction - see github repo readme for usage example
	DeferredRollbackIfNotHandled() error
}
