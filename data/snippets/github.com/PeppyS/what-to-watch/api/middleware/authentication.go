package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Authentication middleware
func Authentication(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Retrieving metadata failed")
	}

	// Skip for now
	// auth, ok := md["authorization"]
	// if !ok {
	// 	return nil, status.Errorf(codes.InvalidArgument, "No auth details supplied")
	// }

	h, err := handler(ctx, req)

	return h, err
}
