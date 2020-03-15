package discovery

import (
	"context"

	"github.com/dogmatiq/configkit/api"
	"github.com/dogmatiq/configkit/api/internal/pb"
	"github.com/dogmatiq/linger/backoff"
	"google.golang.org/grpc"
)

// Connector connects to discovered targets and notifies a client observer if
// they implement the config API.
type Connector struct {
	// Observer is notified when a config API client connects to or disconnects
	// from a server. It must not be nil.
	Observer ClientObserver

	// Dial is the dialer used to connect to the discovered targets.
	// If it is nil, DefaultDialer is used.
	Dial Dialer

	// Ignore is a predicate function that returns true if the given target
	// should be ignored.
	Ignore func(*Target) bool

	// BackoffStrategy controls how long to wait between dialing retries.
	BackoffStrategy backoff.Strategy

	// IsFatal, if non-nil, is called when an error occurs connecting to a
	// target.
	//
	// If it returns true, Run() returns err immediately. If it is nil, or it
	// returns false, dialing is retried as per the backoff strategy.
	IsFatal func(err error) bool
}

// Run connects to a target in order to publish client connect/disconnect
// notifications to the observer.
//
// It retries until ctx is canceled.
func (c *Connector) Run(ctx context.Context, t *Target) error {
	if c.Ignore != nil && c.Ignore(t) {
		return nil
	}

	ctr := &backoff.Counter{
		Strategy: c.BackoffStrategy,
	}

	for {
		err := c.dial(ctx, ctr, t)

		if ctx.Err() != nil {
			return ctx.Err()
		}

		if c.IsFatal != nil && c.IsFatal(err) {
			return err
		}

		if err := ctr.Sleep(ctx, err); err != nil {
			return err
		}
	}
}

// dial attempts to dial the target.
func (c *Connector) dial(
	ctx context.Context,
	ctr *backoff.Counter,
	t *Target,
) error {
	dial := c.Dial
	if dial == nil {
		dial = DefaultDialer
	}

	conn, err := dial(ctx, t)
	if err != nil {
		return err
	}
	defer conn.Close()

	return c.watch(ctx, ctr, conn, t)
}

// watch checks if conn supports the config API and notifies the observer
// accordingly.
//
// It blocks until ctx is canceled or the connection is severed, at which point
// the observer is notified of the disconnection.
func (c *Connector) watch(
	ctx context.Context,
	ctr *backoff.Counter,
	conn *grpc.ClientConn,
	t *Target,
) error {
	// Make a context that closes the gRPC stream when this function exits, as
	// streams have no Close() method.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := pb.NewConfigClient(conn).Watch(ctx, &pb.WatchRequest{})
	if err != nil {
		return err
	}

	client := &Client{
		Client:     api.NewClient(conn),
		Target:     t,
		Connection: conn,
	}

	// The server sends exactly one response. The result of this call will be an
	// empty response if the server implements the config API, or an error if it
	// doesn't.
	if _, err := stream.Recv(); err != nil {
		return err
	}

	c.Observer.ClientConnected(client)
	defer c.Observer.ClientDisconnected(client)

	// We've successfully queried the server, so if there is an error now it's a
	// regular disconnection and we should retry immediately, therefore we reset
	// the failure counter.
	ctr.Reset()

	_, err = stream.Recv()
	return err
}
