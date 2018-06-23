// Package apitoken is a little shim for making per-RPC credentials work easier.
package apitoken

import (
	"context"

	"google.golang.org/grpc/credentials"
)

type apitoken struct {
	token string
}

// NewFromToken creates a new API token from a given string.
func NewFromToken(token string) credentials.PerRPCCredentials {
	return apitoken{token: token}
}

func (at apitoken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": at.token,
	}, nil
}

func (at apitoken) RequireTransportSecurity() bool {
	return false
}
