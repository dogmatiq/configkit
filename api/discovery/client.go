package discovery

import (
	"sync"

	"github.com/dogmatiq/configkit/api"
)

// ClientPublisher is an interface that allows client observers to be registered to
// receive notifications.
type ClientPublisher interface {
	// RegisterClientObserver registers o to be notified when connections to
	// config API servers are established and servered.
	RegisterClientObserver(o ClientObserver)

	// UnregisterClientObserver stops o from being notified when connections to
	// config API servers are established and servered.
	UnregisterClientObserver(o ClientObserver)
}

// ClientObserver is notified when connections to config API servers are
// established and severed.
type ClientObserver interface {
	// ClientConnected is called when a connection to a config API server is
	// established.
	ClientConnected(*api.Client)

	// ClientDisconnected is called when a connection to a config API server is
	// severed.
	ClientDisconnected(*api.Client)
}

// ClientObserverSet is a client observer that publishes to other observers.
//
// It implements both the ClientObserver and ClientPublisher interfaces.
type ClientObserverSet struct {
	m         sync.RWMutex
	observers map[ClientObserver]struct{}
	clients   map[*api.Client]struct{}
}

// RegisterClientObserver registers o to be notified when connections to config
// API servers are established and servered.
func (s *ClientObserverSet) RegisterClientObserver(o ClientObserver) {
	s.m.Lock()
	defer s.m.Unlock()

	if s.observers == nil {
		s.observers = map[ClientObserver]struct{}{}
	} else if _, ok := s.observers[o]; ok {
		return
	}

	s.observers[o] = struct{}{}
	s.notifyOne(ClientObserver.ClientConnected, o)
}

// UnregisterClientObserver stops o from being notified when connections to
// config API servers are established and servered.
func (s *ClientObserverSet) UnregisterClientObserver(o ClientObserver) {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.observers[o]; !ok {
		return
	}

	delete(s.observers, o)
	s.notifyOne(ClientObserver.ClientDisconnected, o)
}

// ClientConnected notifies the registered observers that c has connected.
func (s *ClientObserverSet) ClientConnected(c *api.Client) {
	s.m.Lock()
	defer s.m.Unlock()

	if s.clients == nil {
		s.clients = map[*api.Client]struct{}{}
	} else if _, ok := s.clients[c]; ok {
		return
	}

	s.clients[c] = struct{}{}
	s.notifyAll(ClientObserver.ClientConnected, c)
}

// ClientDisconnected notifies the registered observers that c has disconnected.
func (s *ClientObserverSet) ClientDisconnected(c *api.Client) {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.clients[c]; !ok {
		return
	}

	delete(s.clients, c)
	s.notifyAll(ClientObserver.ClientDisconnected, c)
}

// notifyAll notifies all observers about a change to one client.
func (s *ClientObserverSet) notifyAll(
	fn func(ClientObserver, *api.Client),
	c *api.Client,
) {
	var g sync.WaitGroup

	g.Add(len(s.observers))

	for o := range s.observers {
		o := o // capture loop variable

		go func() {
			defer g.Done()
			fn(o, c)
		}()
	}

	g.Wait()
}

// notifyOne notifies one observer about a change to all clients.
func (s *ClientObserverSet) notifyOne(
	fn func(ClientObserver, *api.Client),
	o ClientObserver,
) {
	var g sync.WaitGroup

	g.Add(len(s.clients))

	for c := range s.clients {
		c := c // capture loop variable

		go func() {
			defer g.Done()
			fn(o, c)
		}()
	}

	g.Wait()
}
