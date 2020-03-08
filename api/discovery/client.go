package discovery

import (
	"context"
	"sync"

	"github.com/dogmatiq/configkit/api"
	"google.golang.org/grpc"
)

// Client is an API client that is aware of the target it connects to.
type Client struct {
	*api.Client

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

// ClientObserverSet is a ClientObserver that publishes to other observers.
type ClientObserverSet struct {
	m         sync.RWMutex
	observers map[ClientObserver]struct{}
	clients   map[*Client]struct{}
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
func (s *ClientObserverSet) ClientConnected(c *Client) {
	s.m.Lock()
	defer s.m.Unlock()

	if s.clients == nil {
		s.clients = map[*Client]struct{}{}
	} else if _, ok := s.clients[c]; ok {
		return
	}

	s.clients[c] = struct{}{}
	s.notifyAll(ClientObserver.ClientConnected, c)
}

// ClientDisconnected notifies the registered observers that c has disconnected.
func (s *ClientObserverSet) ClientDisconnected(c *Client) {
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
	fn func(ClientObserver, *Client),
	c *Client,
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
	fn func(ClientObserver, *Client),
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

// ClientTask is a function executed by a ClientExecutor.
type ClientTask func(context.Context, *Client)

// ClientExecutor is a ClientObserver that executes a function in a new
// goroutine whenever a client connects.
type ClientExecutor struct {
	executor

	// Task is the function to execute when a client connects.
	// The context is canceled when the target becomes unavailable.
	Task ClientTask

	// Parent is the parent context under which the function is called.
	// If it is nil, context.Background() is used.
	Parent context.Context
}

// ClientConnected starts a new goroutine for the given client.
func (e *ClientExecutor) ClientConnected(c *Client) {
	e.start(e.Parent, c, func(ctx context.Context) {
		e.Task(ctx, c)
	})
}

// ClientDisconnected cancels the context associated with any existing goroutine
// for the given client and waits for the goroutine to exit.
func (e *ClientExecutor) ClientDisconnected(c *Client) {
	e.stop(c)
}
