package fixtures

import (
	"sync"

	"github.com/dogmatiq/configkit/api/discovery"
)

var _ discovery.ApplicationObserver = (*ApplicationObserver)(nil)

// ApplicationObserver is a mock of the discovery.ApplicationObserver interface.
type ApplicationObserver struct {
	m                          sync.Mutex
	ApplicationAvailableFunc   func(*discovery.Application)
	ApplicationUnavailableFunc func(*discovery.Application)
}

// ApplicationAvailable calls o.ApplicationAvailableFunc(a) if it is non-nil.
func (o *ApplicationObserver) ApplicationAvailable(a *discovery.Application) {
	if o.ApplicationAvailableFunc != nil {
		o.m.Lock()
		defer o.m.Unlock()
		o.ApplicationAvailableFunc(a)
	}
}

// ApplicationUnavailable calls o.ApplicationUnavailableFunc(a) if it is non-nil.
func (o *ApplicationObserver) ApplicationUnavailable(a *discovery.Application) {
	if o.ApplicationUnavailableFunc != nil {
		o.m.Lock()
		defer o.m.Unlock()
		o.ApplicationUnavailableFunc(a)
	}
}
