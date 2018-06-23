package handlers

import (
	"golang.org/x/net/context"

	pb "github.com/zaquestion/routines"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.RoutinesServer {
	return routinesService{}
}

type routinesService struct{}

// ScrapTrelloReset implements Service.
func (s routinesService) ScrapTrelloReset(ctx context.Context, in *pb.ScrapTrelloResetRequest) (*pb.ScrapTrelloResetReply, error) {
	var resp pb.ScrapTrelloResetReply
	resp = pb.ScrapTrelloResetReply{
	// Err:
	}
	return &resp, nil
}

// GetRoutines implements Service.
func (s routinesService) GetRoutines(ctx context.Context, in *pb.GetRoutinesRequest) (*pb.GetRoutinesReply, error) {
	var resp pb.GetRoutinesReply
	resp = pb.GetRoutinesReply{
	// Routines:
	// Err:
	}
	return &resp, nil
}
