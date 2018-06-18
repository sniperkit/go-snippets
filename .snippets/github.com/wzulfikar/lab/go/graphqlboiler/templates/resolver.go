package templates

const Resolver = `package resolvers

import (
	"context"
	"errors"
	"strconv"

	graphql "github.com/neelance/graphql-go"
	{{ if .HasScalar }}
	"{{.Repo}}/modules/graphql/scalar"
	{{ end }}
	"{{.Repo}}/models"
)

type {{.ResolverName}} struct {
	rr *{{.RootResolver}}
	o  *models.{{.TypeName}}
}

func (rr *{{.RootResolver}}) {{.TypeName}}(ctx context.Context, args struct{ ID graphql.ID }) (*{{.ResolverName}}, error) {
	id, err := strconv.Atoi(string(args.ID))
	if err != nil {
		return nil, errors.New("Failed to get ID of {{.TypeName}}")
	}

	o, err := models.Find{{.TypeName}}(rr.Db, uint(id))
	if err != nil {
		return nil, errors.New("{{.TypeName}} not found")
	}

	return &{{.ResolverName}}{rr: rr, o: o}, nil
}
`
