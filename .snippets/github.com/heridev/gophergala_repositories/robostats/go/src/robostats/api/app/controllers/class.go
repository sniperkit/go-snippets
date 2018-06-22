package controllers

import (
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	"robostats/models/class"
	"robostats/models/user"
)

func init() {
	revel.InterceptFunc(addHeaderCORS, revel.AFTER, &Class{})
}

type classEnvelope struct {
	Class class.Class `json:"deviceClass"`
}

type classesEnvelope struct {
	Classes []*class.Class `json:"deviceClasses"`
}

type Class struct {
	Common
}

func (c Class) Create() revel.Result {
	var err error
	var u *user.User
	var k classEnvelope

	if u, err = c.requireAuthorization(); err != nil {
		return c.StatusUnauthorized()
	}

	if err = c.decodeBody(&k); err != nil {
		return c.StatusBadRequest()
	}

	k.Class.UserID = u.ID

	if err = k.Class.Create(); err != nil {
		return c.writeError(err)
	}

	return c.dataCreated(classEnvelope{k.Class})
}

func (c Class) Index() revel.Result {
	var err error
	var u *user.User
	var classes []*class.Class

	if u, err = c.requireAuthorization(); err != nil {
		return c.StatusUnauthorized()
	}

	if classes, err = class.GetByUserID(u.ID); err != nil {
		return c.writeError(err)
	}

	return c.dataGeneric(classesEnvelope{classes})
}

func (c Class) Get() revel.Result {
	var err error
	var k *class.Class

	id := c.Params.Get("id")

	if _, err = c.requireAuthorization(); err != nil {
		return c.StatusUnauthorized()
	}

	if k, err = class.GetByID(bson.ObjectIdHex(id)); err != nil {
		return c.writeError(err)
	}

	return c.dataGeneric(classEnvelope{*k})
}

func (c Class) Remove() revel.Result {
	var err error
	var k *class.Class

	id := c.Params.Get("id")

	if _, err = c.requireAuthorization(); err != nil {
		return c.StatusUnauthorized()
	}

	if k, err = class.GetByID(bson.ObjectIdHex(id)); err != nil {
		return c.writeError(err)
	}

	if err = k.Remove(); err != nil {
		return c.statusNotFound()
	}

	return c.StatusOK()
}
