package api

import (
	"context"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/enginekit/grpc/configgrpc"
	"google.golang.org/grpc"
)

// Client wraps a [configgrpc.ConfigAPIClient] to unmarshal the server's
// responses into types that implement the core configkit
// [configkit.Application] and [configkit.Handler] interfaces.
type Client struct {
	Client configgrpc.ConfigAPIClient
}

// NewClient returns a new configuration client for the given connection.
func NewClient(conn grpc.ClientConnInterface) *Client {
	return &Client{
		configgrpc.NewConfigAPIClient(conn),
	}
}

// ListApplications returns the configurations of the applications hosted by
// the server. The handler objects in the returned configuration are nil.
func (c *Client) ListApplications(
	ctx context.Context,
) ([]configkit.Application, error) {
	req := &configgrpc.ListApplicationsRequest{}
	res, err := c.Client.ListApplications(ctx, req)
	if err != nil {
		return nil, err
	}

	var configs []configkit.Application
	for _, in := range res.GetApplications() {
		out, err := configkit.FromProto(in)
		if err != nil {
			return nil, err
		}

		configs = append(configs, out)
	}

	return configs, nil
}
