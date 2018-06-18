package resolvers

import (
	"context"
	"errors"
	"strconv"

	graphql "github.com/neelance/graphql-go"
	
	"/models"
)

type personResolver struct {
	rr *RootResolver
	o  *models.Person
}

func (rr *RootResolver) Person(ctx context.Context, args struct{ID graphql.ID}) (*personResolver, error) {
	id, err := strconv.Atoi(string(args.ID))
	if err != nil {
		return nil, errors.New("Failed to get ID of Person")
	}

	o, err := models.FindPerson(rr.Db, uint(id))
	if err != nil {
		return nil, errors.New("Person not found")
	}

	return &personResolver{rr: rr, o: o}, nil
}

func (r *personResolver) Name() (string, error) {
	return r.o.Name, nil
}

func (r *personResolver) Age() (int32, error) {
	return int32(r.o.Age), nil
}
