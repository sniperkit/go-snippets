package handler

import (
	"code.google.com/p/go.net/context"
	"code.google.com/p/goprotobuf/proto"

	"github.com/gophergala/Gevvent/services/gevvent-event-service/model"
	"github.com/gophergala/Gevvent/services/gevvent-event-service/repo"

	listproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/list"
)

type List struct{}

func (e *List) Call(ctx context.Context, req *listproto.Request, rsp *listproto.Response) error {
	var events []*model.Event
	var err error

	switch req.GetViewType() {
	case listproto.ViewType_UPCOMING:
		events, err = repo.GetUpcoming(req.UserID, req.GetPage(), req.GetCount())
	case listproto.ViewType_INVITATIONS:
		events, err = repo.GetInvitations(req.UserID, req.GetPage(), req.GetCount())
	case listproto.ViewType_HOSTING:
		events, err = repo.GetHosting(req.UserID, req.GetPage(), req.GetCount())
	case listproto.ViewType_PAST:
		events, err = repo.GetPast(req.UserID, req.GetPage(), req.GetCount())
	}
	if err != nil {
		return err
	}

	for _, event := range events {
		rsp.Events = append(rsp.Events, listproto.Event{
			ID:          event.ID,
			UserID:      event.UserID,
			Name:        event.Name,
			Description: event.Description,
			When: listproto.TimeRange{
				From: event.When.From.Unix(),
				To:   event.When.To.Unix(),
			},
			Where: listproto.Location{
				Lat:     event.Where.LatLng.Lat,
				Lng:     event.Where.LatLng.Lon,
				Address: event.Where.Address,
			},
			Private:    proto.Bool(event.Private),
			PublicAddr: proto.Bool(event.PublicAddr),
		})
	}

	return nil
}
