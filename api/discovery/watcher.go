package discovery

import (
	"context"

	"github.com/dogmatiq/configkit/api"
	"github.com/dogmatiq/configkit/api/internal/pb"
	"github.com/dogmatiq/dodeca/logging"
	"github.com/dogmatiq/linger/backoff"
	"google.golang.org/grpc"
)

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

// watch checks if conn supports the config API and notifies obs accordingly.
//
// It blocks until ctx is canceled or the connection is severed, at which point
// the obs is notified of the disconnection.
func (w *watcher) watch(
	ctx context.Context,
	conn *grpc.ClientConn,
) error {
	stream, err := pb.NewConfigClient(conn).Watch(ctx, &pb.WatchRequest{})
	if err != nil {
		return err
	}

	client := api.NewClient(conn)
	w.Observer.ClientConnected(client)
	defer w.Observer.ClientDisconnected(client)

	// We've successfully queried the server, so if there is an error now its a
	// regular disconnection and we should retry immediately, therefore we reset
	// the failure counter.
	w.Failures.Reset()

	// Wait for the server to disconnect.
	_, err = stream.Recv()
	return err
}
