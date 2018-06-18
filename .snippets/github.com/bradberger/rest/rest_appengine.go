// +build appengine

package rest

import (
	"github.com/bradberger/context"

	"google.golang.org/appengine"
)

func setNamespace(ctx context.Context) (context.Context, error) {
	if Namespace == nil {
		return ctx, nil
	}
	ns, err := Namespace(ctx)
	if err != nil {
		return ctx, err
	}
	return appengine.Namespace(ctx, ns)
}
