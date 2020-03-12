package static

import (
	"context"

	"github.com/dogmatiq/configkit/api/discovery"
)

// Discoverer notifies a target observer of a fixed list of targets.
type Discoverer struct {
	Targets  []*discovery.Target
	Observer discovery.TargetObserver
}

// Run adds the notifies the observer of the targets availability until ctx is
// canceled.
func (d *Discoverer) Run(ctx context.Context) error {
	for _, t := range d.Targets {
		d.Observer.TargetAvailable(t)
		defer d.Observer.TargetUnavailable(t)
	}

	<-ctx.Done()
	return ctx.Err()
}
