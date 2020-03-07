package fixtures

import (
	"sync"

	"github.com/dogmatiq/configkit/api/discovery"
)

// TargetObserver is a mock of the discovery.TargetObserver interface.
type TargetObserver struct {
	m                     sync.Mutex
	TargetAvailableFunc   func(*discovery.Target)
	TargetUnavailableFunc func(*discovery.Target)
}

// TargetAvailable calls o.TargetAvailableFunc(t) if it is non-nil.
func (o *TargetObserver) TargetAvailable(t *discovery.Target) {
	if o.TargetAvailableFunc != nil {
		o.m.Lock()
		defer o.m.Unlock()
		o.TargetAvailableFunc(t)
	}
}

// TargetUnavailable calls o.TargetUnavailableFunc(t) if it is non-nil.
func (o *TargetObserver) TargetUnavailable(t *discovery.Target) {
	if o.TargetUnavailableFunc != nil {
		o.m.Lock()
		defer o.m.Unlock()
		o.TargetUnavailableFunc(t)
	}
}
