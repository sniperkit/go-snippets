package middleware

import "github.com/henrylee2cn/faygo"

type Common struct {
	//
}

func (t *Common) Serve(ctx *faygo.Context) error {
	faygo.RenderVar("ctx", ctx)
	breadcrumbs := make(map[string]string)
	breadcrumbs["/"] = "首页"
	return nil
}
