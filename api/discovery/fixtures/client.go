package fixtures

import (
	"sync"

	"github.com/dogmatiq/configkit/api/discovery"
)

// ClientObserver is a mock of the discovery.ClientObserver interface.
type ClientObserver struct {
	m                      sync.Mutex
	ClientConnectedFunc    func(*discovery.Client)
	ClientDisconnectedFunc func(*discovery.Client)
}

// ClientConnected calls o.ClientConnectedFunc(c) if it is non-nil.
func (o *ClientObserver) ClientConnected(c *discovery.Client) {
	if o.ClientConnectedFunc != nil {
		o.m.Lock()
		defer o.m.Unlock()
		o.ClientConnectedFunc(c)
	}
}

// ClientDisconnected calls o.ClientDisconnectedFunc(c) if it is non-nil.
func (o *ClientObserver) ClientDisconnected(c *discovery.Client) {
	if o.ClientDisconnectedFunc != nil {
		o.m.Lock()
		defer o.m.Unlock()
		o.ClientDisconnectedFunc(c)
	}
}
