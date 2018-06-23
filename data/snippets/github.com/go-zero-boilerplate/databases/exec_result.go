package databases

type ExecResult interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}
