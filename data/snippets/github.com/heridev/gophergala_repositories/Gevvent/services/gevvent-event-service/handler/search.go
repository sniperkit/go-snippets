package handler

import (
	"code.google.com/p/go.net/context"
	"code.google.com/p/goprotobuf/proto"
	log "github.com/cihub/seelog"

	"github.com/gophergala/Gevvent/services/gevvent-event-service/repo"

	searchproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/search"
)

type Search struct{}

func (e *Search) Call(ctx context.Context, req *searchproto.Request, rsp *searchproto.Response) error {
	var events []*repo.SearchResult
	var err error

	log.Debugf("Searching for events [query=%s, loc=%v]", req.Query, req.Where)

	if req.Where == nil {
		events, err = repo.Search(req.GetUserID(), req.Query, req.GetPage(), req.GetCount())
	} else {
		events, err = repo.SearchNearest(req.GetUserID(), req.Query, req.Where.Lat, req.Where.Lng, req.GetPage(), req.GetCount())
	}
	if err != nil {
		return err
	}

	for _, event := range events {
		rsp.Events = append(rsp.Events, searchproto.Event{
			ID:          event.ID,
			UserID:      event.UserID,
			Name:        event.Name,
			Description: event.Description,
			When: searchproto.TimeRange{
				From: event.When.From.Unix(),
				To:   event.When.To.Unix(),
			},
			Where: searchproto.Location{
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
