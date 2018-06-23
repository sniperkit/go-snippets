package ctrl

import (
	"github.com/henrylee2cn/faygo"
	"net/http"
	"strings"
	"faygo-demo1/model"
)

type Login struct {
}

func (t *Login) Serve(ctx *faygo.Context) (err error) {
	ctx.Render(http.StatusOK, "tpl/login.html", faygo.Map{
		"title": "用户登录",
	})
	return
}

type LoginPost struct {
	Username string `param:"<in:formData><required><len:5:50>"`
	Password string `param:"<in:formData><required><len:5:20>"`
}

func (t *LoginPost) Serve(ctx *faygo.Context) (err error) {
	t.Username = strings.TrimSpace(t.Username)
	user := new(model.User)
	user.DB().Where(model.User{Username: t.Username}).Find(user)
	if user.ID == 0 {
		ctx.JSONMsg(http.StatusOK, 0, "该用户不存在")
	}
	if !user.CheckPassword(t.Password) {
		ctx.JSONMsg(http.StatusOK, 0, "密码错误")
	}
	return
}
