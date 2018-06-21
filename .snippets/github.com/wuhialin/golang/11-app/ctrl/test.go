package ctrl

import (
	"github.com/henrylee2cn/faygo"
	"mime/multipart"
)

/*
Test test struct handler
*/
type Test struct {
	Token     string                `param:"<in:query><required><desc:token必填>"`
	Name      string                `param:"<in:formData><required><len:1:10><desc:your name (1~10 words)>"`
	Age       uint8                 `param:"<in:formData><range:1:100><desc:your age (1~100)>"`
	Avatar    *multipart.FileHeader `param:"<in:formData><desc:your avatar>"`
	AvatarUrl string                `param:"-"`
}

// Serve impletes Handler.
func (t *Test) Serve(ctx *faygo.Context) error {
	info, err := ctx.SaveFile("avatar", false)
	if err != nil {
		return ctx.JSON(412, faygo.Map{"error": err.Error()}, true)
	}
	t.AvatarUrl = info.Url
	return ctx.JSON(200, t, true)
}

// Doc returns the API's note, result or parameters information.
func (t *Test) Doc() faygo.Doc {
	return faygo.Doc{
		Note:   "test struct handler",
		Return: "// JSON\n{}",
	}
}
