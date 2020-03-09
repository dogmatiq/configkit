package fixtures

import (
	"context"
	"sync"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/api"
	"github.com/dogmatiq/configkit/api/discovery"
)

var _ api.Client = (*Client)(nil)

// Client is a mock of the api.Client interface.
type Client struct {
	// ListApplicationIdentities returns the identities of applications hosted
	// by the server.
	ListApplicationIdentitiesFunc func(ctx context.Context) ([]configkit.Identity, error)

	// ListApplications returns the configurations of the applications hosted by
	// the server. The handler objects in the returned configuration are nil.
	ListApplicationsFunc func(ctx context.Context) ([]configkit.Application, error)
}

// ListApplicationIdentities returns the identities of applications hosted
// by the server.
//
// If h.ListApplicationIdentitiesFunc is nil, it returns (nil, nil),
// otherwise it calls h.ListApplicationIdentitiesFunc(ctx).
func (c *Client) ListApplicationIdentities(ctx context.Context) ([]configkit.Identity, error) {
	if c.ListApplicationIdentitiesFunc == nil {
		return nil, nil
	}

	return c.ListApplicationIdentitiesFunc(ctx)
}

// ListApplications returns the configurations of the applications hosted by the
// server. The handler objects in the returned configuration are nil.
//
// If h.ListApplicationsFunc is nil, it returns (nil, nil), otherwise it calls
// h.ListApplicationsFunc(ctx).
func (c *Client) ListApplications(ctx context.Context) ([]configkit.Application, error) {
	if c.ListApplicationsFunc == nil {
		return nil, nil
	}

	return c.ListApplicationsFunc(ctx)
}

var _ discovery.ClientObserver = (*ClientObserver)(nil)

// ClientObserver is a mock of the discovery.ClientObserver interface.
type ClientObserver struct {
	m                      sync.Mutex
	ClientConnectedFunc    func(*discovery.Client)
	ClientDisconnectedFunc func(*discovery.Client)
}

// ClientConnected calls o.ClientConnectedFunc(c) if it is non-nil.
func (o *ClientObserver) ClientConnected(c *discovery.Client) {
	if o.ClientConnectedFunc != nil {
		o.m.Lock()
		defer o.m.Unlock()
		o.ClientConnectedFunc(c)
	}
}

// ClientDisconnected calls o.ClientDisconnectedFunc(c) if it is non-nil.
func (o *ClientObserver) ClientDisconnected(c *discovery.Client) {
	if o.ClientDisconnectedFunc != nil {
		o.m.Lock()
		defer o.m.Unlock()
		o.ClientDisconnectedFunc(c)
	}
}
