package controllers

import (
	"github.com/revel/revel"
	"robostats/models/user"
)

func init() {
	revel.InterceptFunc(addHeaderCORS, revel.AFTER, &User{})
}

type userEnvelope struct {
	User user.User `json:"user"`
}

type User struct {
	Common
}

func (c User) Me() revel.Result {
	var err error
	var u *user.User

	if u, err = c.requireAuthorization(); err != nil {
		return c.StatusUnauthorized()
	}

	return c.dataGeneric(userEnvelope{*u})
}

func (c User) Login() revel.Result {
	var email string
	var password string

	c.Params.Bind(&email, "email")
	c.Params.Bind(&password, "password")

	var err error
	var u *user.User

	if u, err = user.Login(email, password); err != nil {
		return c.StatusUnauthorized()
	}

	data := map[string]string{
		"access_token": u.Session.Token,
		"user_id":      u.ID.Hex(),
		"token_type":   "bearer",
	}

	return c.dataGeneric(data)
}

func (c User) Create() revel.Result {
	var u userEnvelope
	var err error

	if err = c.decodeBody(&u); err != nil {
		c.StatusBadRequest()
	}

	if err = u.User.Create(); err != nil {
		return c.writeError(err)
	}

	return c.dataCreated(userEnvelope{u.User})
}
