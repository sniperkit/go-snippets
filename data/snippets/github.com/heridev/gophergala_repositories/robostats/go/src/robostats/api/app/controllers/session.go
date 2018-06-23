package controllers

import (
	"errors"
	"fmt"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2/bson"
	"robostats/models/beat"
	"robostats/models/instance"
	"robostats/models/session"
	"robostats/models/user"
)

func init() {
	revel.InterceptFunc(addHeaderCORS, revel.AFTER, &Session{})
}

type sessionEnvelope struct {
	Session session.Session `json:"deviceSession"`
}

type sessionsEnvelope struct {
	Sessions []*session.Session `json:"deviceSessions"`
}

type Session struct {
	Common
}

func (c Session) Create() revel.Result {
	var err error
	var u *user.User
	var k sessionEnvelope
	var i *instance.Instance

	if u, err = c.requireAuthorization(); err != nil {
		return c.StatusUnauthorized()
	}

	if err = c.decodeBody(&k); err != nil {
		return c.StatusBadRequest()
	}

	if !k.Session.InstanceID.Valid() {
		return c.StatusBadRequest()
	}

	if i, err = instance.GetByID(k.Session.InstanceID); err != nil {
		return c.StatusBadRequest()
	}

	k.Session.UserID = u.ID
	k.Session.ClassID = i.ClassID
	k.Session.InstanceID = i.ID

	if err = k.Session.Create(); err != nil {
		return c.writeError(err)
	}

	return c.dataCreated(sessionEnvelope{k.Session})
}

// Index returns all sessions.
func (c Session) Index() revel.Result {
	var err error
	var u *user.User
	var sessions []*session.Session

	var deviceInstanceID string
	var deviceClassID string

	c.Params.Bind(&deviceInstanceID, "device_instance_id")
	c.Params.Bind(&deviceClassID, "device_class_id")

	if u, err = c.requireAuthorization(); err != nil {
		return c.StatusUnauthorized()
	}

	if deviceClassID != "" {
		if sessions, err = session.GetByClassID(bson.ObjectIdHex(deviceClassID)); err != nil {
			return c.writeError(err)
		}
	} else if deviceInstanceID != "" {
		if sessions, err = session.GetByInstanceID(bson.ObjectIdHex(deviceInstanceID)); err != nil {
			return c.writeError(err)
		}
	} else {
		if sessions, err = session.GetByUserID(u.ID); err != nil {
			return c.writeError(err)
		}
	}

	return c.dataGeneric(sessionsEnvelope{sessions})
}

// Get returns a specific session.
func (c Session) Get() revel.Result {
	var err error
	var k *session.Session

	id := c.Params.Get("id")

	if _, err = c.requireAuthorization(); err != nil {
		return c.StatusUnauthorized()
	}

	if k, err = session.GetByID(bson.ObjectIdHex(id)); err != nil {
		return c.writeError(err)
	}

	return c.dataGeneric(sessionEnvelope{*k})
}

// Remove deletes a session.
func (c Session) Remove() revel.Result {
	var err error
	var k *session.Session

	id := c.Params.Get("id")

	if _, err = c.requireAuthorization(); err != nil {
		return c.StatusUnauthorized()
	}

	if k, err = session.GetByID(bson.ObjectIdHex(id)); err != nil {
		return c.writeError(err)
	}

	if err = k.Remove(); err != nil {
		return c.statusNotFound()
	}

	return c.StatusOK()
}

func (c Session) TimeSeries() revel.Result {
	var sessionID string
	var key []string

	var err error

	if _, err = c.requireAuthorization(); err != nil {
		return c.StatusUnauthorized()
	}

	c.Params.Bind(&sessionID, "session_id")
	c.Params.Bind(&key, "key")

	if len(key) == 0 {
		return c.writeError(errors.New("Missing at least one key."))
	}
	fmt.Printf("key: %v\n", key)

	if sessionID == "" {
		return c.writeError(errors.New("Missing session_id."))
	}

	var bs []*beat.Beat
	if bs, err = beat.GetBySessionID(bson.ObjectIdHex(sessionID)); err != nil {
		return c.writeError(err)
	}

	ts := make(TimeEvents, 0, len(bs))

	for _, b := range bs {
		te := TimeEvent{
			LocalTime: b.LocalTime,
			Event:     make(map[string]interface{}),
		}
		for _, k := range key {
			te.Event[k], _ = b.GetKeyValue(k)
		}
		ts = append(ts, te)
	}

	return c.dataGeneric(ts.Envelope())
}
