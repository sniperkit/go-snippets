package echo

import (
	"net/url"

	"github.com/labstack/echo"

	"github.com/go-zero-boilerplate/api-resources/resources/inputs"
)

func newGetInput(echoCtx echo.Context) inputs.Get       { return &input{echoCtx: echoCtx} }
func newPostInput(echoCtx echo.Context) inputs.Post     { return &input{echoCtx: echoCtx} }
func newPatchInput(echoCtx echo.Context) inputs.Patch   { return &input{echoCtx: echoCtx} }
func newPutInput(echoCtx echo.Context) inputs.Put       { return &input{echoCtx: echoCtx} }
func newDeleteInput(echoCtx echo.Context) inputs.Delete { return &input{echoCtx: echoCtx} }

type input struct {
	echoCtx echo.Context
}

func (i *input) Query(key string) inputs.Value {
	valueUnescaped, err := url.QueryUnescape(i.echoCtx.QueryParam(key))
	if err != nil {
		return inputs.NewValueError(err, "Query", key)
	}
	return inputs.NewValue("Query", key, valueUnescaped)
}

func (i *input) Param(key string) inputs.Value {
	valueUnescaped, err := url.QueryUnescape(i.echoCtx.Param(key))
	if err != nil {
		return inputs.NewValueError(err, "Query", key)
	}
	return inputs.NewValue("Param", key, valueUnescaped)
}

func (i *input) BindBody(body interface{}) error {
	return i.echoCtx.Bind(body)
}
