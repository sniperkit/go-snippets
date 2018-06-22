package handler

import (
	"fmt"
	"time"

	"code.google.com/p/go.net/context"
	rtypes "github.com/dancannon/gorethink/types"

	"github.com/gophergala/Gevvent/services/gevvent-event-service/model"
	"github.com/gophergala/Gevvent/services/gevvent-event-service/repo"

	createproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/create"
)

type Create struct{}

func (e *Create) Call(ctx context.Context, req *createproto.Request, rsp *createproto.Response) error {
	if req.Name == "" {
		return fmt.Errorf("Must specify an event name")
	}

	whenFrom := time.Unix(req.When.From, 0)
	whenTo := time.Unix(req.When.To, 0)

	event := &model.Event{
		UserID:      req.UserID,
		Name:        req.Name,
		Description: req.Description,
		When: model.TimeRange{
			From: whenFrom,
			To:   whenTo,
		},
		Where: model.Location{
			Address: req.Where.Address,
			LatLng: rtypes.Point{
				Lat: req.Where.Lat,
				Lon: req.Where.Lng,
			},
		},
		Private:    req.GetPrivate(),
		PublicAddr: req.GetPublicAddr(),
	}

	var err error
	event, err = repo.CreateEvent(event)
	if err != nil {
		return err
	}

	rsp.ID = event.ID

	return nil
}
