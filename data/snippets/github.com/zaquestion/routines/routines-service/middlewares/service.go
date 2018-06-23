package middlewares

import (
	pb "github.com/zaquestion/routines"
)

func WrapService(in pb.RoutinesServer) pb.RoutinesServer {
	return in
}
