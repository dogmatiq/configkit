package discovery

import (
	"context"

	"google.golang.org/grpc"
)

// Dialer connects to a gRPC target.
type Dialer func(context.Context, *Target) (*grpc.ClientConn, error)

// DefaultDialer is the default dialer used to connect to a gRPC target.
func DefaultDialer(ctx context.Context, t *Target) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, t.Name, t.Options...)
}
