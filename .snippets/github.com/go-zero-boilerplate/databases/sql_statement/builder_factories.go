package sql_statement

import (
	"github.com/go-zero-boilerplate/databases"
)

type BuilderFactories struct {
	Insert InsertBuilderFactory
	Update UpdateBuilderFactory
	Delete DeleteBuilderFactory
	Select SelectBuilderFactory
}

func NewBuilderFactories(db databases.Database) *BuilderFactories {
	return &BuilderFactories{
		Insert: NewInsertBuilderFactory(db),
		Update: NewUpdateBuilderFactory(db),
		Delete: NewDeleteBuilderFactory(db),
		Select: NewSelectBuilderFactory(db),
	}
}
