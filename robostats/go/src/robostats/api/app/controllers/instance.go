package controllers

import (
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	"log"
	"robostats/models/instance"
	"robostats/models/user"
)

func init() {
	revel.InterceptFunc(addHeaderCORS, revel.AFTER, &Instance{})
}

type instanceEnvelope struct {
	Instance instance.Instance `json:"deviceInstance"`
}

type instancesEnvelope struct {
	Instances []*instance.Instance `json:"deviceInstances"`
}

type Instance struct {
	Common
}

func (c Instance) Create() revel.Result {
	var err error
	var u *user.User
	var k instanceEnvelope

	if u, err = c.requireAuthorization(); err != nil {
		return c.StatusUnauthorized()
	}

	log.Printf("decode...")

	if err = c.decodeBody(&k); err != nil {
		log.Printf("decode error...: %v", err)
		return c.StatusBadRequest()
	}

	if !k.Instance.ClassID.Valid() {
		log.Printf("NOT VALID\n")
		return c.StatusBadRequest()
	}

	k.Instance.UserID = u.ID

	if err = k.Instance.Create(); err != nil {
		return c.writeError(err)
	}

	return c.dataCreated(instanceEnvelope{k.Instance})
}

// Index returns all instances.
func (c Instance) Index() revel.Result {
	var err error
	var u *user.User
	var instances []*instance.Instance
	var deviceClassID string

	c.Params.Bind(&deviceClassID, "device_class_id")

	if u, err = c.requireAuthorization(); err != nil {
		return c.StatusUnauthorized()
	}

	if deviceClassID != "" {
		if instances, err = instance.GetByClassID(bson.ObjectIdHex(deviceClassID)); err != nil {
			return c.writeError(err)
		}
	} else {
		if instances, err = instance.GetByUserID(u.ID); err != nil {
			return c.writeError(err)
		}
	}

	return c.dataGeneric(instancesEnvelope{instances})
}

// Get returns a specific instance.
func (c Instance) Get() revel.Result {
	var err error
	var k *instance.Instance

	id := c.Params.Get("id")

	if _, err = c.requireAuthorization(); err != nil {
		return c.StatusUnauthorized()
	}

	if k, err = instance.GetByID(bson.ObjectIdHex(id)); err != nil {
		return c.writeError(err)
	}

	return c.dataGeneric(instanceEnvelope{*k})
}

// Remove deletes an instance.
func (c Instance) Remove() revel.Result {
	var err error
	var k *instance.Instance

	id := c.Params.Get("id")

	if _, err = c.requireAuthorization(); err != nil {
		return c.StatusUnauthorized()
	}

	if k, err = instance.GetByID(bson.ObjectIdHex(id)); err != nil {
		return c.writeError(err)
	}

	if err = k.Remove(); err != nil {
		return c.statusNotFound()
	}

	return c.StatusOK()
}
