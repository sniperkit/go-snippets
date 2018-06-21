package middleware

import (
	"github.com/henrylee2cn/faygo"
	"regexp"
	"faygo-demo1/model"
)

//不需要开始SESSION的请求
var ignoreRegexp []*regexp.Regexp

//不需要检验是否登录的路径
var ignoreCheckRegexp []*regexp.Regexp

func init() {
	ignoreRegexp = []*regexp.Regexp{
		regexp.MustCompile(`^/(static|upload)/.*`),
	}
}

type CheckUser struct {
	//
}

func (t *CheckUser) Serve(ctx *faygo.Context) (err error) {
	for _, exp := range ignoreRegexp {
		if exp.MatchString(ctx.Path()) {
			return
		}
	}
	faygo.Debug("start session")
	if _, err := ctx.StartSession(); err != nil {
		faygo.Error("start session, failure:", err)
	}
	for _, exp := range ignoreCheckRegexp {
		if exp.MatchString(ctx.Path()) {
			return
		}
	}
	tmp := ctx.CurSession.Get("user")
	if tmp != nil {
		user := tmp.(model.User)
		faygo.RenderVar("user", user)
	}
	return
}
