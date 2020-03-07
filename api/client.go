package api

import (
	"context"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/api/internal/pb"
	"google.golang.org/grpc"
)

// Client is used to query a server about its application configurations.
type Client struct {
	client pb.ConfigClient
}

// NewClient returns a new configuration client for the given connection.
func NewClient(conn grpc.ClientConnInterface) *Client {
	return &Client{
		pb.NewConfigClient(conn),
	}
}

// ListApplicationIdentities returns the identities of applications hosted by
// the server.
func (c *Client) ListApplicationIdentities(
	ctx context.Context,
) (_ []configkit.Identity, err error) {
	req := &pb.ListApplicationIdentitiesRequest{}
	res, err := c.client.ListApplicationIdentities(ctx, req)
	if err != nil {
		return nil, err
	}

	var idents []configkit.Identity
	for _, in := range res.GetIdentities() {
		out, err := unmarshalIdentity(in)
		if err != nil {
			return nil, err
		}

		idents = append(idents, out)
	}

	return idents, nil
}

// ListApplications returns the configurations of the applications hosted by
// the server. The handler objects in the returned configuration are nil.
func (c *Client) ListApplications(
	ctx context.Context,
) ([]configkit.Application, error) {
	req := &pb.ListApplicationsRequest{}
	res, err := c.client.ListApplications(ctx, req)
	if err != nil {
		return nil, err
	}

	var configs []configkit.Application
	for _, in := range res.GetApplications() {
		out, err := unmarshalApplication(in)
		if err != nil {
			return nil, err
		}

		configs = append(configs, out)
	}

	return configs, nil
}
