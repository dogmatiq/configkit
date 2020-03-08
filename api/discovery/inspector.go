package discovery

import (
	"context"

	"github.com/dogmatiq/configkit"
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
}

// Inspect queries a client in order to publish application
// available/unavailable notifications to the observer.
//
// It blocks until ctx is canceled.
func (i *Inspector) Inspect(ctx context.Context, c *Client) error {
	apps, err := c.ListApplications(ctx)
	if err != nil {
		return err
	}

	empty := true

	for _, a := range apps {
		if i.Ignore != nil && i.Ignore(a) {
			continue
		}

		x := &Application{
			Application: a,
			Client:      c,
		}

		empty = false
		i.Observer.ApplicationAvailable(x)
		defer i.Observer.ApplicationUnavailable(x)
	}

	if empty {
		return nil
	}

	<-ctx.Done()
	return ctx.Err()
}
