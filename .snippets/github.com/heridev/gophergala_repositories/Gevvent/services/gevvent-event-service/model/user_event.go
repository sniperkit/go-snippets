package model

import (
	readuserproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/readuser"
)

type UserEventStatus string

const (
	NotGoing UserEventStatus = "not_going"
	Going                    = "going"
	Invited                  = "invited"
	Saved                    = "saved"
)

func (s UserEventStatus) ToProto() readuserproto.Status {
	switch s {
	case Going:
		return readuserproto.Status_GOING
	case Invited:
		return readuserproto.Status_INVITED
	case Saved:
		return readuserproto.Status_SAVED
	default:
		return readuserproto.Status_NOT_GOING
	}
}

type UserEvent struct {
	ID      string          `gorethink:"id", json:"id"` // ID is the event ID + user ID concatenated
	EventID string          `gorethink:"event_id", json:"event_id"`
	UserID  string          `gorethink:"user_id", json:"user_id"`
	Status  UserEventStatus `gorethink:"status", json:"status"`
}
