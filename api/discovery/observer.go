package discovery

import (
	"sync"

	"github.com/dogmatiq/configkit/api"
)

// Observer is notified when config API servers are discovered.
type Observer interface {
	// Connected is called when a connection to a config API server is
	// established.
	Connected(*api.Client)

	// Disconnected is called when a connection to a config API server is
	// severed.
	Disconnected(*api.Client)
}

// ObserverSet is an Observer that dispatches to other observers.
type ObserverSet struct {
	m         sync.RWMutex
	observers map[Observer]struct{}
	clients   map[*api.Client]struct{}
}

// Add registers o to be notified when connections to config API servers are
// established and servered.
func (s *ObserverSet) Add(o Observer) {
	s.m.Lock()
	defer s.m.Unlock()

	if s.observers == nil {
		s.observers = map[Observer]struct{}{}
	} else if _, ok := s.observers[o]; ok {
		return
	}

	s.observers[o] = struct{}{}
	s.notifyOne(Observer.Connected, o)
}

// Remove stops o from being notified when connections to config API servers are
// established and servered.
func (s *ObserverSet) Remove(o Observer) {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.observers[o]; !ok {
		return
	}

	delete(s.observers, o)
	s.notifyOne(Observer.Disconnected, o)
}

// Connected is called when a connection to a config API server is established.
func (s *ObserverSet) Connected(c *api.Client) {
	s.m.Lock()
	defer s.m.Unlock()

	if s.clients == nil {
		s.clients = map[*api.Client]struct{}{}
	} else if _, ok := s.clients[c]; ok {
		return
	}

	s.clients[c] = struct{}{}
	s.notifyAll(Observer.Connected, c)
}

// Disconnected is called when a connection to a config API server is servered.
func (s *ObserverSet) Disconnected(c *api.Client) {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.clients[c]; !ok {
		return
	}

	delete(s.clients, c)
	s.notifyAll(Observer.Disconnected, c)
}

// notifyAll notifies all observers about a change to one client.
func (s *ObserverSet) notifyAll(
	fn func(Observer, *api.Client),
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
func (s *ObserverSet) notifyOne(
	fn func(Observer, *api.Client),
	o Observer,
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
