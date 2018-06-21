package ctrl

import (
	"github.com/henrylee2cn/faygo"
	"net/http"
)

type Home struct {
	//
}

func (t *Home) Serve(ctx *faygo.Context) error {
	ctx.Render(http.StatusOK, "tpl/index.html", nil)
	return nil
}
