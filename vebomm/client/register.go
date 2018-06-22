package client

import (
	"github.com/phaikawl/wade/page"

	"github.com/gophergala/vebomm/core"
	"gopkg.in/validator.v2"
)

var _vmr *RegisterVM

type RegisterVM struct {
	AppMain
	PasswordConfirm string
	User            core.User
	Gender          string
	Result          core.RegisterResult
	valErrs         error
}

func (vm *RegisterVM) HasError(key string) bool {
	return valError(vm.valErrs, key) != ""
}

func (vm *RegisterVM) DoRegister() {
	if vm.Gender == "M" {
		vm.User.Gender = 1
	}

	vm.valErrs = validator.Validate(&vm.User)
	if vm.valErrs != nil {
		jsAlert.Error("Please correct the form errors.")
		return
	}

	go func() {
		resp, err := vm.Http.POST("/api/register", vm.User)
		if err != nil || resp.Failed() {
			jsAlert.Error("Connection problem or server error.")
			return
		}

		err = resp.ParseJSON(&vm.Result)
		if err != nil || !vm.Result.ValOk {
			jsAlert.Error("Register failed.")
			return
		}

		update()

		if !vm.Result.DupEmail && !vm.Result.DupUsername {
			jsAlert.Success("Register successful.")
			vm.AppMain.saveCred(&vm.User)
			vm.PageMgr.GoToPage(HomePage)
		}
	}()
}

func (vm *RegisterVM) PasswordError() bool {
	return vm.PasswordConfirm != vm.User.Password
}

func (am AppMain) RegisterHandler(ctx page.Context) {
	_vmr = &RegisterVM{
		AppMain: am,
	}

	_vmr.User.BirthYear = 1995
}
