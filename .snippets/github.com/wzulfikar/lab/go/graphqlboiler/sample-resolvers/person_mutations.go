package resolvers

import (
	"context"

	graphql "github.com/neelance/graphql-go"
	"/models"
)

func (rr *RootResolver) CreatePerson(ctx context.Context, args struct{ID graphql.ID}) (*personResolver, error) {
	
	panic("TODO: handle CreatePerson mutation")
	
	var o *models.Person
	// Sample code:
	// o, err := NewBusiness(*args.Request)
	// if err != nil {
	// 	return nil, app.Error(err)
	// }

	return &personResolver{rr: rr, o: o}, nil
}
