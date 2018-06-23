package handler

import (
	"code.google.com/p/go.net/context"
	"code.google.com/p/goprotobuf/proto"

	"github.com/gophergala/Gevvent/services/gevvent-event-service/repo"

	newestproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/newest"
)

type Newest struct{}

func (e *Newest) Call(ctx context.Context, req *newestproto.Request, rsp *newestproto.Response) error {
	events, err := repo.GetNewest(req.GetPage(), req.GetCount())
	if err != nil {
		return err
	}

	for _, event := range events {
		rsp.Events = append(rsp.Events, newestproto.Event{
			ID:          event.ID,
			UserID:      event.UserID,
			Name:        event.Name,
			Description: event.Description,
			When: newestproto.TimeRange{
				From: event.When.From.Unix(),
				To:   event.When.To.Unix(),
			},
			Where: newestproto.Location{
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
