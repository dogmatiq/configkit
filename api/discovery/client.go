package discovery

import (
	"github.com/dogmatiq/configkit/api"
	"google.golang.org/grpc"
)

// Client is an API client that is aware of the target it connects to.
type Client struct {
	api.Client

	// Target is the discovered gRPC target that the client connects to.
	Target *Target

	// Connection is the gRPC connection to the target.
	Connection *grpc.ClientConn
}

// ClientObserver is notified when connections to config API servers are
// established and severed.
type ClientObserver interface {
	// ClientConnected is called when a connection to a config API server is
	// established.
	ClientConnected(*Client)

	// ClientDisconnected is called when a connection to a config API server is
	// severed.
	ClientDisconnected(*Client)
}

// ClientObserverSet is a client observer that publishes to other observers.
//
// It implements both the ClientObserver and ClientPublisher interfaces.
type ClientObserverSet struct {
	observerSet
}

// NewClientObserverSet registers the given observers with a new observer set
// and returns it.
func NewClientObserverSet(observers ...ClientObserver) *ClientObserverSet {
	s := &ClientObserverSet{}

	for _, o := range observers {
		s.RegisterClientObserver(o)
	}

	return s
}

// RegisterClientObserver registers o to be notified when connections to config
// API servers are established and servered.
func (s *ClientObserverSet) RegisterClientObserver(o ClientObserver) {
	s.register(o, func(e interface{}) {
		o.ClientConnected(e.(*Client))
	})
}

// UnregisterClientObserver stops o from being notified when connections to
// config API servers are established and servered.
func (s *ClientObserverSet) UnregisterClientObserver(o ClientObserver) {
	s.unregister(o, func(e interface{}) {
		o.ClientDisconnected(e.(*Client))
	})
}

// ClientConnected notifies the registered observers that c has connected.
func (s *ClientObserverSet) ClientConnected(c *Client) {
	s.add(c, func(o interface{}) {
		o.(ClientObserver).ClientConnected(c)
	})
}

// ClientDisconnected notifies the registered observers that c has disconnected.
func (s *ClientObserverSet) ClientDisconnected(c *Client) {
	s.remove(c, func(o interface{}) {
		o.(ClientObserver).ClientDisconnected(c)
	})
}
