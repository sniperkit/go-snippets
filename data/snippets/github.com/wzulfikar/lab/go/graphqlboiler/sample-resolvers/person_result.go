package resolvers

import (
	"context"

	"github.com/volatiletech/sqlboiler/queries/qm"
	"/models"
	"/modules/graphql/app"
)

type PersonResultResolver struct {
	totalCount int64
	pageInfo   *app.PageInfoResolver
	items      []*personResolver
}

func (r *PersonResultResolver) PageInfo() *app.PageInfoResolver {
	return r.pageInfo
}

type PersonFilter struct {
	// Add Person filter here
}

type PersonesArgs struct {
	Page   *app.PageArgs
	Filter *PersonFilter
}

func (rr *RootResolver) People(ctx context.Context, args *PersonesArgs) (*PersonResultResolver, error) {
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

	count, err := models.People(rr.Db, mods...).Count()
	if err != nil {
		return nil, nil
	}

	pageInfo, err := args.Page.PageInfo(count)
	if err != nil {
		return nil, err
	}

	mods = append(mods, pageInfo.QM()...)

	o, err := models.People(rr.Db, mods...).All()
	if err != nil {
		return nil, err
	}

	result := &PersonResultResolver{
		totalCount: count,
		pageInfo:   pageInfo.Resolver(),
	}
	for _, row := range o {
		result.items = append(result.items, &personResolver{rr: rr, o: row})
	}

	return result, nil
}

func (r *PersonResultResolver) TotalCount() (int32, error) {
	return int32(r.totalCount), nil
}

func (r *PersonResultResolver) Items() ([]*personResolver, error) {
	return r.items, nil
}
