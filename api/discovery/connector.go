package discovery

import (
	"context"
	"sync"

	"github.com/dogmatiq/dodeca/logging"
	"github.com/dogmatiq/linger/backoff"
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
