// +build !appengine,go1.7

package rest

import (
	"github.com/bradberger/context"
)

func setNamespace(ctx context.Context) (context.Context, error) {
	return ctx, nil
}
