package discovery

import (
	"context"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/dodeca/logging"
	"github.com/dogmatiq/linger/backoff"
)

// Inspector queries connected clients and notifies an application observer
// of the applications it has.
type Inspector struct {
	// Observer is notified when an application is discovered via a config API
	// client. It must not be nil.
	Observer ApplicationObserver

	// Ignore is a predicate function that returns true if the given application
	// should be ignored.
	Ignore func(configkit.Application) bool

	// BackoffStrategy controls how long to wait between inspection retries.
	BackoffStrategy backoff.Strategy

	// Logger is the target for log messages about inspection failures.
	Logger logging.Logger
}

// Run queries a client in order to publish application available/unavailable
// notifications to the observer.
//
// It retries until the applications are successfully discovered. After which it
// blocks until ctx is canceled, or returns immediately if no applications were
// found.
func (i *Inspector) Run(ctx context.Context, c *Client) error {
	ctr := &backoff.Counter{
		Strategy: i.BackoffStrategy,
	}

	for {
		err := i.list(ctx, c)

		if err == nil {
			return nil
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		logging.Log(
			i.Logger,
			"unable to inspect applications on '%s' target: %s",
			c.Target.Name,
			err,
		)

		if err := ctr.Sleep(ctx, err); err != nil {
			return err
		}
	}
}

// list queries the client for applications and notifies the observer
// accordingly.
func (i *Inspector) list(ctx context.Context, c *Client) error {
	configs, err := c.ListApplications(ctx)
	if err != nil {
		return err
	}

	empty := true

	for _, cfg := range configs {
		if i.Ignore != nil && i.Ignore(cfg) {
			continue
		}

		empty = false
		app := &Application{
			Application: cfg,
			Client:      c,
		}

		i.Observer.ApplicationAvailable(app)
		defer i.Observer.ApplicationUnavailable(app)
	}

	if empty {
		return nil
	}

	<-ctx.Done()
	return ctx.Err()
}
