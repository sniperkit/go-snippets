package authorize

import (
	"net/http"

	"github.com/go-zero-boilerplate/facades/tasks"
)

//Generic just calls the given check func and if an error is returned will return a StatusUnauthorized error
func Generic(fn func() error) tasks.LoadFunc {
	return func(taskCtx *tasks.Context) {
		if err := fn(); err != nil {
			taskCtx.Error = tasks.NewError(err, http.StatusUnauthorized)
			return
		}
	}
}
