package templates

const QueryResult = `package resolvers

import (
	"context"

	"github.com/volatiletech/sqlboiler/queries/qm"
	"{{.Repo}}/models"
	"{{.Repo}}/modules/graphql/app"
)

type {{.TypeName}}ResultResolver struct {
	totalCount int64
	pageInfo   *app.PageInfoResolver
	items      []*{{.ResolverName}}
}

func (r *{{.TypeName}}ResultResolver) PageInfo() *app.PageInfoResolver {
	return r.pageInfo
}

type {{.TypeName}}Filter struct {
	// Add {{.TypeName}} filter here
}

type {{.TypeName}}esArgs struct {
	Page   *app.PageArgs
	Filter *{{.TypeName}}Filter
}

func (rr *{{.RootResolver}}) {{.ModelPlural}}(ctx context.Context, args *{{.TypeName}}esArgs) (*{{.TypeName}}ResultResolver, error) {
	var mods []qm.QueryMod

	if args.Filter != nil {
		whereMods, err := app.WhereMods([]app.MaybeMod{
			// Add mods for qm.Where here. See: app.MaybeMod
		})
		if err != nil {
			return nil, err
		}

		mods = append(mods, whereMods...)
	}

	count, err := models.{{.ModelPlural}}(rr.Db, mods...).Count()
	if err != nil {
		return nil, nil
	}

	pageInfo, err := args.Page.PageInfo(count)
	if err != nil {
		return nil, err
	}

	mods = append(mods, pageInfo.QM()...)

	o, err := models.{{.ModelPlural}}(rr.Db, mods...).All()
	if err != nil {
		return nil, err
	}

	result := &{{.TypeName}}ResultResolver{
		totalCount: count,
		pageInfo:   pageInfo.Resolver(),
	}
	for _, row := range o {
		result.items = append(result.items, &{{.ResolverName}}{rr: rr, o: row})
	}

	return result, nil
}

func (r *{{.TypeName}}ResultResolver) TotalCount() (int32, error) {
	return int32(r.totalCount), nil
}

func (r *{{.TypeName}}ResultResolver) Items() ([]*{{.ResolverName}}, error) {
	return r.items, nil
}
`
