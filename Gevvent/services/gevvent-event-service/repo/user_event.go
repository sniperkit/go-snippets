package repo

import (
	"fmt"
	r "github.com/dancannon/gorethink"

	"github.com/gophergala/Gevvent/services/gevvent-event-service/model"
	"github.com/gophergala/Gevvent/services/gevvent-lib/rethinkdb"
)

func UpdateUserEvent(event *model.UserEvent) (*model.UserEvent, error) {
	// Set ID
	event.ID = fmt.Sprintf("%s|%s", event.EventID, event.UserID)

	s, err := rethinkdb.Session()
	if err != nil {
		return nil, err
	}

	err = r.Table("user_events").Insert(event, r.InsertOpts{
		Conflict: "replace",
	}).Exec(s)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func DeleteUserEvent(eventID, userID string) error {
	s, err := rethinkdb.Session()
	if err != nil {
		return err
	}

	return r.Table("user_events").Get(userID + "|" + eventID).Delete().Exec(s)
}

func GetUserEvent(eventID, userID string) (*model.UserEvent, error) {
	var event *model.UserEvent
	id := fmt.Sprintf("%s|%s", eventID, userID)

	s, err := rethinkdb.Session()
	if err != nil {
		return event, err
	}

	rows, err := r.Table("user_events").Get(id).Run(s)
	if err != nil {
		return event, err
	}

	err = rows.One(&event)
	if err != nil {
		if err == r.ErrEmptyResult {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return event, nil
}

func GetAttendees(eventID string) ([]*model.UserEvent, error) {
	var events []*model.UserEvent

	s, err := rethinkdb.Session()
	if err != nil {
		return events, err
	}

	rows, err := r.Table("user_events").GetAllByIndex("event", eventID).Filter(func(row r.Term) interface{} {
		return row.Field("status").Eq("going").Or(row.Field("status").Eq("invited"))
	}).Run(s)
	if err != nil {
		return events, err
	}

	err = rows.All(&events)
	if err != nil {
		return events, err
	}

	return events, err
}

type ListResult struct {
	*model.Event

	Status string `gorethink:"status", json:"status"`
}

func GetUpcoming(userID string, page, count int64) ([]*model.Event, error) {
	var events []*model.Event

	s, err := rethinkdb.Session()
	if err != nil {
		return events, err
	}

	rows, err := r.Table("user_events").Filter(func(row r.Term) interface{} {
		return row.Field("status").Eq("going").And(row.Field("user_id").Eq(userID))
	}).EqJoin("event_id", r.Table("events")).Zip().Filter(func(row r.Term) interface{} {
		return row.Field("when").Field("from").Gt(r.Now())
	}).Skip((page - 1) * count).Limit(count).Run(s)
	if err != nil {
		return events, err
	}

	err = rows.All(&events)
	if err != nil {
		return events, err
	}

	return events, err
}

func GetInvitations(userID string, page, count int64) ([]*model.Event, error) {
	var events []*model.Event

	s, err := rethinkdb.Session()
	if err != nil {
		return events, err
	}

	rows, err := r.Table("user_events").Filter(func(row r.Term) interface{} {
		return row.Field("status").Eq("invited").And(row.Field("user_id").Eq(userID))
	}).EqJoin("event_id", r.Table("events")).Zip().Filter(func(row r.Term) interface{} {
		return row.Field("when").Field("from").Gt(r.Now())
	}).Skip((page - 1) * count).Limit(count).Run(s)
	if err != nil {
		return events, err
	}

	err = rows.All(&events)
	if err != nil {
		return events, err
	}

	return events, err
}

func GetPast(userID string, page, count int64) ([]*model.Event, error) {
	var events []*model.Event

	s, err := rethinkdb.Session()
	if err != nil {
		return events, err
	}

	rows, err := r.Table("user_events").Filter(func(row r.Term) interface{} {
		return row.Field("status").Eq("going").And(row.Field("user_id").Eq(userID))
	}).EqJoin("event_id", r.Table("events")).Zip().Filter(func(row r.Term) interface{} {
		return row.Field("when").Field("from").Le(r.Now())
	}).Skip((page - 1) * count).Limit(count).Run(s)
	if err != nil {
		return events, err
	}

	err = rows.All(&events)
	if err != nil {
		return events, err
	}

	return events, err
}
