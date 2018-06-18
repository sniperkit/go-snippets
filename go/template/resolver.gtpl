package resolvers

import (
	"context"
)

type {{.TypeName}} struct {
	rr *{{.RootResolver}}
	o  *models.{{.TypeName}}
}

func (rr *{{.RootResolver}}) typeName(ctx context.Context, args struct{ID graphql.ID}) (*{{.TypeName}}, error) {
	panic("TODO: initialize object for resolver")
	o := &models.{{.TypeName}}{} 
	return &{{.TypeName}}{rr: rr, o: o}, nil
}