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

	// Ignore is a list of application identity keys to ignore.
	Ignore []string

	apps []*Application
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

	defer func() {
		for _, a := range i.apps {
			i.Observer.ApplicationUnavailable(a)
		}
	}()

	for _, a := range apps {
		if i.isIgnored(a) {
			continue
		}

		x := &Application{
			Application: a,
			Client:      c,
		}

		i.apps = append(i.apps, x)
		i.Observer.ApplicationAvailable(x)
	}

	<-ctx.Done()
	return ctx.Err()
}

// isIgnored returns true if app should be ignored.
func (i *Inspector) isIgnored(app configkit.Application) bool {
	for _, k := range i.Ignore {
		if k == app.Identity().Key {
			return true
		}
	}

	return false
}
