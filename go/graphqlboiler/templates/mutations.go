package templates

const Mutations = `package resolvers

import (
	"context"

	graphql "github.com/neelance/graphql-go"
	"{{.Repo}}/models"
)

func (rr *{{.RootResolver}}) Create{{.TypeName}}(ctx context.Context, args struct{ ID graphql.ID }) (*{{.ResolverName}}, error) {
	
	panic("TODO: handle Create{{.TypeName}} mutation")
	
	var o *models.{{.TypeName}}
	// Sample code:
	// o, err := NewBusiness(*args.Request)
	// if err != nil {
	// 	return nil, app.Error(err)
	// }

	return &{{.ResolverName}}{rr: rr, o: o}, nil
}
`
