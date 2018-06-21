package middleware

import (
	"log"
	"time"
	"context"

	"google.golang.org/grpc"
)

// Logging middleware logs requests
func Logging(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	h, err := handler(ctx, req)

	log.Printf("Request - Method:%s\tDuration:%s\tError:%v", info.FullMethod, time.Since(start), err)

	return h, err
}
