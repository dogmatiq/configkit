package discovery

import (
	"context"
	"sync"

	"github.com/dogmatiq/configkit/api"
	"github.com/dogmatiq/configkit/api/internal/pb"
	"github.com/dogmatiq/dodeca/logging"
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

	// BackoffStrategy controls how long to wait between failures to dial a
	// discovered target.
	BackoffStrategy backoff.Strategy

	// Logger is the target for log messages about dialing failures.
	Logger logging.Logger

	m       sync.Mutex
	cancels map[*Target]context.CancelFunc
}

// TargetAvailable is called when a target is becomes available.
func (c *Connector) TargetAvailable(t *Target) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.cancels == nil {
		c.cancels = map[*Target]context.CancelFunc{}
	} else if _, ok := c.cancels[t]; ok {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.cancels[t] = cancel

	w := &watcher{
		Target:   t,
		Observer: c.Observer,
		Dial:     c.Dial,
		Failures: backoff.Counter{
			Strategy: c.BackoffStrategy,
		},
		Logger: c.Logger,
	}

	go w.Run(ctx)
}

// TargetUnavailable is called when a target becomes unavailable.
func (c *Connector) TargetUnavailable(t *Target) {
	c.m.Lock()
	defer c.m.Unlock()

	if cancel, ok := c.cancels[t]; ok {
		delete(c.cancels, t)
		cancel()
	}
}

// watcher connects to a target in order to publish client connect/disconnect
// notifications to an observer.
//
// The observer is only notified if the target implements the config API.
type watcher struct {
	// Target is the target to dial. It must not be nil.
	Target *Target

	// Observer is notified when a config API client connects to or disconnects
	// from the target, if the target implements the config API. It must not be
	// nil.
	Observer ClientObserver

	// Dial is the dialer used to connect to the target. If it is nil,
	// DefaultDialer is used.
	Dial Dialer

	// Failures keeps track of the dial failures.
	Failures backoff.Counter

	// Logger is the target for dial failure messages. If it is nil,
	// logging.DefaultLogger is used.
	Logger logging.Logger
}

// Run dials the target and publishes client connection/disconnection
// notifications to the observer.
//
// It runs until ctx is canceled.
func (w *watcher) Run(ctx context.Context) error {
	for {
		err := w.dial(ctx)

		if ctx.Err() != nil {
			return ctx.Err()
		}

		logging.Log(
			w.Logger,
			"unable to watch '%s' target: %s",
			w.Target.Name,
			err,
		)

		if err := w.Failures.Sleep(ctx, err); err != nil {
			return err
		}
	}
}

// dial attempts to dial the target.
func (w *watcher) dial(ctx context.Context) error {
	dial := w.Dial
	if dial == nil {
		dial = DefaultDialer
	}

	conn, err := dial(ctx, w.Target)
	if err != nil {
		return err
	}
	defer conn.Close()

	return w.watch(ctx, conn)
}

// watch checks if conn supports the config API and notifies the observer
// accordingly.
//
// It blocks until ctx is canceled or the connection is severed, at which point
// the observer is notified of the disconnection.
func (w *watcher) watch(
	ctx context.Context,
	conn *grpc.ClientConn,
) error {
	// Make a context that closes the gRPC stream when this function exists, as
	// streams have no Close() method.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := pb.NewConfigClient(conn).Watch(ctx, &pb.WatchRequest{})
	if err != nil {
		return err
	}

	client := &Client{
		Client: api.NewClient(conn),
		Target: w.Target,
	}

	// The server sends exactly one response. The result of this call will be an
	// empty response if the server implements the config API, or an error if it
	// doesn't.
	if _, err := stream.Recv(); err != nil {
		return err
	}

	w.Observer.ClientConnected(client)
	defer w.Observer.ClientDisconnected(client)

	// We've successfully queried the server, so if there is an error now its a
	// regular disconnection and we should retry immediately, therefore we reset
	// the failure counter.
	w.Failures.Reset()

	_, err = stream.Recv()
	return err
}
