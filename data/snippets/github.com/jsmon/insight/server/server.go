package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/homebot/core/utils"
	"github.com/homebot/insight"
	insightV1 "github.com/homebot/protobuf/pkg/api/insight/v1"
)

// Insight is a homebot/api/insight/v1/log.proto:LogSink server
type Insight struct {
	store insight.Storage
}

// AddMessage adds a new log message
func (i *Insight) AddMessage(ctx context.Context, in *insightV1.LogMessage) (*empty.Empty, error) {
	if err := i.store.Save(in); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// ListMessages returns all log messages
func (i *Insight) ListMessages(ctx context.Context, in *insightV1.ListMessageRequest) (*insightV1.ListMessageResponse, error) {
	res, err := i.store.List()
	if err != nil {
		return nil, err
	}

	start, end, t, err := utils.Paginate(in.GetPageToken(), len(res), in.GetPageSize())
	if err != nil {
		return nil, err
	}

	return &insightV1.ListMessageResponse{
		Messages:      res[start:end],
		NextPageToken: t,
	}, nil
}

// DeleteMessage deletes a log message with a given id
func (i *Insight) DeleteMessage(ctx context.Context, in *insightV1.DeleteMessageRequest) (*empty.Empty, error) {
	if err := i.store.Delete(in.GetId()); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

// GetMessage returns the log message for the given ID
func (i *Insight) GetMessage(ctx context.Context, in *insightV1.GetMessageRequest) (*insightV1.LogMessage, error) {
	msg, err := i.store.Get(in.GetId())
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// SearchMessages searches for log messages
func (i *Insight) SearchMessages(ctx context.Context, in *insightV1.SearchMessagesRequest) (*insightV1.ListMessageResponse, error) {
	filter, err := insight.BuilderFromProto(in)
	if err != nil {
		return nil, err
	}

	res, err := i.store.Search(filter)
	if err != nil {
		return nil, err
	}

	start, end, token, err := utils.Paginate(in.GetPageToken(), len(res), in.GetPageSize())
	if err != nil {
		return nil, err
	}

	return &insightV1.ListMessageResponse{
		Messages:      res[start:end],
		NextPageToken: token,
	}, nil
}
