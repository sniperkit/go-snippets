package requests

import (
	"net/http"

	"github.com/go-zero-boilerplate/facades/tasks"
)

//GenericCheck just calls the given check func and if an error is returned will return a StatusBadRequest error
func GenericCheck(fn func() error) tasks.LoadFunc {
	return func(taskCtx *tasks.Context) {
		if err := fn(); err != nil {
			taskCtx.Error = tasks.NewError(err, http.StatusBadRequest)
			return
		}
	}
}
