package simpledns

import (
	"context"
	"net"
	"time"

	"github.com/dogmatiq/configkit/api/discovery"
	"github.com/dogmatiq/linger"
)

// Resolver is an interface for the subset of net.Resolver used by the
// discoverer.
type Resolver interface {
	LookupHost(ctx context.Context, host string) ([]string, error)
}

const (
	// DefaultInterval is the default interval at which DNS queries are performed.
	DefaultInterval = 10 * time.Second

	// DefaultServerPort is the default server TCP port.
	DefaultServerPort = "https"
)

// Discoverer periodically performs a DNS query to discover API servers and
// notifies a target observer.
type Discoverer struct {
	QueryHost  string
	ServerPort string
	Observer   discovery.TargetObserver
	Interval   time.Duration
	Resolver   Resolver

	targets map[string]*discovery.Target
}

// Run performs discovery until ctx is canceled or an error occurs.
func (d *Discoverer) Run(ctx context.Context) error {
	for {
		addrs, err := d.query(ctx)
		if err != nil {
			return err
		}

		d.update(addrs)

		if err := linger.Sleep(ctx, d.Interval, DefaultInterval); err != nil {
			return err
		}
	}
}

// update sends notifications to the observer about the targets that have become
// available/unavailable based on new query results.
func (d *Discoverer) update(addrs []string) {
	port := d.ServerPort
	if port == "" {
		port = DefaultServerPort
	}

	prev := d.targets
	d.targets = make(map[string]*discovery.Target, len(addrs))

	for _, a := range addrs {
		n := net.JoinHostPort(a, port)

		t, ok := prev[n]

		if !ok {
			t = &discovery.Target{
				Name: n,
				MetaData: discovery.MetaData{
					QueryHostKey: d.QueryHost,
				},
			}

			d.Observer.TargetAvailable(t)
		}

		d.targets[n] = t
	}

	for n, t := range prev {
		if _, ok := d.targets[n]; !ok {
			d.Observer.TargetUnavailable(t)
		}
	}
}

// query performs a DNS query to find API targets.
func (d *Discoverer) query(ctx context.Context) ([]string, error) {
	r := d.Resolver
	if r == nil {
		r = net.DefaultResolver
	}

	addrs, err := r.LookupHost(ctx, d.QueryHost)
	if err != nil {
		if x, ok := err.(*net.DNSError); ok {
			if x.IsTemporary || x.IsNotFound {
				return nil, nil
			}
		}

		return nil, err
	}

	return addrs, nil
}
