package discovery

import (
	"context"
	"sync"

	"github.com/dogmatiq/configkit"
)

// Application is an application configuration that is aware of the client it
// was obtained from.
type Application struct {
	configkit.Application

	// Client is the client that queried the application configuration.
	Client *Client
}

// ApplicationObserver is notified when Dogma applications are discovered.
type ApplicationObserver interface {
	// ApplicationAvailable is called when an application becomes available.
	ApplicationAvailable(*Application)

	// ApplicationUnavailable is called when an application becomes unavailable.
	ApplicationUnavailable(*Application)
}

// ApplicationObserverSet is an ApplicationObserver that publishes to other
// observers.
type ApplicationObserverSet struct {
	m         sync.RWMutex
	observers map[ApplicationObserver]struct{}
	apps      map[*Application]struct{}
}

// NewApplicationObserverSet registers the given observers with a new observer
// set and returns it.
func NewApplicationObserverSet(observers ...ApplicationObserver) *ApplicationObserverSet {
	s := &ApplicationObserverSet{}

	for _, o := range observers {
		s.RegisterApplicationObserver(o)
	}

	return s
}

// RegisterApplicationObserver registers o to be notified when applications
// become available and unavailable.
func (s *ApplicationObserverSet) RegisterApplicationObserver(o ApplicationObserver) {
	s.m.Lock()
	defer s.m.Unlock()

	if s.observers == nil {
		s.observers = map[ApplicationObserver]struct{}{}
	} else if _, ok := s.observers[o]; ok {
		return
	}

	s.observers[o] = struct{}{}
	s.notifyOne(ApplicationObserver.ApplicationAvailable, o)
}

// UnregisterApplicationObserver stops o from being notified when applications
// become available and unavailable.
func (s *ApplicationObserverSet) UnregisterApplicationObserver(o ApplicationObserver) {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.observers[o]; !ok {
		return
	}

	delete(s.observers, o)
	s.notifyOne(ApplicationObserver.ApplicationUnavailable, o)
}

// ApplicationAvailable notifies the registered observers that t is available.
func (s *ApplicationObserverSet) ApplicationAvailable(a *Application) {
	s.m.Lock()
	defer s.m.Unlock()

	if s.apps == nil {
		s.apps = map[*Application]struct{}{}
	} else if _, ok := s.apps[a]; ok {
		return
	}

	s.apps[a] = struct{}{}
	s.notifyAll(ApplicationObserver.ApplicationAvailable, a)
}

// ApplicationUnavailable notifies the registered observers that t is unavailable.
func (s *ApplicationObserverSet) ApplicationUnavailable(a *Application) {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.apps[a]; !ok {
		return
	}

	delete(s.apps, a)
	s.notifyAll(ApplicationObserver.ApplicationUnavailable, a)
}

// notifyAll notifies all observers about a change to one application.
func (s *ApplicationObserverSet) notifyAll(
	fn func(ApplicationObserver, *Application),
	a *Application,
) {
	var g sync.WaitGroup

	g.Add(len(s.observers))

	for o := range s.observers {
		o := o // capture loop variable

		go func() {
			defer g.Done()
			fn(o, a)
		}()
	}

	g.Wait()
}

// notifyOne notifies one observer about a change to all applications.
func (s *ApplicationObserverSet) notifyOne(
	fn func(ApplicationObserver, *Application),
	o ApplicationObserver,
) {
	var g sync.WaitGroup

	g.Add(len(s.apps))

	for t := range s.apps {
		t := t // capture loop variable

		go func() {
			defer g.Done()
			fn(o, t)
		}()
	}

	g.Wait()
}

// ApplicationTask is a function executed by an ApplicationExecutor.
type ApplicationTask func(context.Context, *Application)

// ApplicationExecutor is an ApplicationObserver that executes a function in a
// new goroutine whenever an application becomes available.
type ApplicationExecutor struct {
	executor

	// Task is the function to execute when an application becomes available.
	// The context is canceled when the application becomes unavailable.
	Task ApplicationTask

	// Parent is the parent context under which the function is called.
	// If it is nil, context.Background() is used.
	Parent context.Context
}

// ApplicationAvailable starts a new goroutine for the given application.
func (e *ApplicationExecutor) ApplicationAvailable(a *Application) {
	e.start(e.Parent, a, func(ctx context.Context) {
		e.Task(ctx, a)
	})
}

// ApplicationUnavailable cancels the context associated with any existing
// goroutine for the given application and waits for the goroutine to exit.
func (e *ApplicationExecutor) ApplicationUnavailable(a *Application) {
	e.stop(a)
}
