package rest

import (
	"errors"

	"github.com/bradberger/context"
)

var (
	// UserFunc is a used defined function which defines which user is associated with the current request
	UserFunc func(ctx context.Context, user interface{}) error

	ErrNoUserFunc = errors.New("no user func defined")
)

func User(ctx context.Context, user interface{}) error {
	if UserFunc == nil {
		return ErrNoUserFunc
	}
	return UserFunc(ctx, user)
}
