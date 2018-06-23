package sql_statement

type BaseStatement interface {
	Execute() error
}
