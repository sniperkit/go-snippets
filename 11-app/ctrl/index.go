package ctrl

import (
	"fmt"
	"github.com/henrylee2cn/faygo"
	"net/http"
	"time"
)

type Index struct {
}

func (i *Index) Serve(ctx *faygo.Context) error {
	start := time.Now()
	return ctx.Render(http.StatusOK, `tpl/index.html`, faygo.Map{
		"runtime": fmt.Sprint(time.Since(start)),
	})
}
