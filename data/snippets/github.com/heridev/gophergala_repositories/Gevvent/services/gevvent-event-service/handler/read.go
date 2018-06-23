package handler

import (
	"code.google.com/p/go.net/context"
	"code.google.com/p/goprotobuf/proto"

	"github.com/gophergala/Gevvent/services/gevvent-event-service/repo"

	readproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/read"
)

type Read struct{}

func (e *Read) Call(ctx context.Context, req *readproto.Request, rsp *readproto.Response) error {
	event, err := repo.Get(req.ID)
	if err != nil {
		return err
	}
	if event == nil {
		return nil
	}

	rsp.Event = &readproto.Event{
		ID:          event.ID,
		UserID:      event.UserID,
		Name:        event.Name,
		Description: event.Description,
		When: readproto.TimeRange{
			From: event.When.From.Unix(),
			To:   event.When.To.Unix(),
		},
		Where: readproto.Location{
			Lat:     event.Where.LatLng.Lat,
			Lng:     event.Where.LatLng.Lon,
			Address: event.Where.Address,
		},
		Private:    proto.Bool(event.Private),
		PublicAddr: proto.Bool(event.PublicAddr),
	}

	return nil
}
