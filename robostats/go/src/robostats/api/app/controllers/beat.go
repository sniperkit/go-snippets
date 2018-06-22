package controllers

import (
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	"robostats/models/beat"
	"robostats/models/session"
	"robostats/models/user"
)

func init() {
	revel.InterceptFunc(addHeaderCORS, revel.AFTER, &Beat{})
}

type beatEnvelope struct {
	Beat beat.Beat `json:"deviceEvent"`
}

type beatsEnvelope struct {
	Beats []*beat.Beat `json:"deviceEvents"`
}

type Beat struct {
	Common
}

func (c Beat) Create() revel.Result {
	var err error
	var u *user.User
	var k beatEnvelope
	var s *session.Session

	if u, err = c.requireAuthorization(); err != nil {
		return c.StatusUnauthorized()
	}

	if err = c.decodeBody(&k); err != nil {
		return c.StatusBadRequest()
	}

	if !k.Beat.SessionID.Valid() {
		return c.StatusBadRequest()
	}

	if s, err = session.GetByID(k.Beat.SessionID); err != nil {
		return c.StatusBadRequest()
	}

	k.Beat.UserID = u.ID
	k.Beat.ClassID = s.ClassID
	k.Beat.InstanceID = s.InstanceID
	k.Beat.SessionID = s.ID

	if err = k.Beat.Create(); err != nil {
		return c.writeError(err)
	}

	return c.dataCreated(beatEnvelope{k.Beat})
}

// Index returns all beats.
func (c Beat) Index() revel.Result {
	var err error
	var u *user.User
	var beats []*beat.Beat

	var deviceClassID string
	var deviceInstanceID string
	var deviceSessionID string

	if u, err = c.requireAuthorization(); err != nil {
		return c.StatusUnauthorized()
	}

	if deviceClassID != "" {
		if beats, err = beat.GetByClassID(bson.ObjectIdHex(deviceClassID)); err != nil {
			return c.writeError(err)
		}
	} else if deviceInstanceID != "" {
		if beats, err = beat.GetByInstanceID(bson.ObjectIdHex(deviceInstanceID)); err != nil {
			return c.writeError(err)
		}
	} else if deviceSessionID != "" {
		if beats, err = beat.GetBySessionID(bson.ObjectIdHex(deviceSessionID)); err != nil {
			return c.writeError(err)
		}
	} else {
		if beats, err = beat.GetByUserID(u.ID); err != nil {
			return c.writeError(err)
		}
	}

	return c.dataGeneric(beatsEnvelope{beats})
}

// Get returns a specific beat.
func (c Beat) Get() revel.Result {
	var err error
	var k *beat.Beat

	id := c.Params.Get("id")

	if _, err = c.requireAuthorization(); err != nil {
		return c.StatusUnauthorized()
	}

	if k, err = beat.GetByID(bson.ObjectIdHex(id)); err != nil {
		return c.writeError(err)
	}

	return c.dataGeneric(beatEnvelope{*k})
}

// Remove deletes a beat.
func (c Beat) Remove() revel.Result {
	var err error
	var k *beat.Beat

	id := c.Params.Get("id")

	if _, err = c.requireAuthorization(); err != nil {
		return c.StatusUnauthorized()
	}

	if k, err = beat.GetByID(bson.ObjectIdHex(id)); err != nil {
		return c.writeError(err)
	}

	if err = k.Remove(); err != nil {
		return c.statusNotFound()
	}

	return c.StatusOK()
}
