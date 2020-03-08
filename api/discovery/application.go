package discovery

import (
	"context"

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
	observerSet
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
	s.register(o, func(e interface{}) {
		o.ApplicationAvailable(e.(*Application))
	})
}

// UnregisterApplicationObserver stops o from being notified when applications
// become available and unavailable.
func (s *ApplicationObserverSet) UnregisterApplicationObserver(o ApplicationObserver) {
	s.unregister(o, func(e interface{}) {
		o.ApplicationUnavailable(e.(*Application))
	})
}

// ApplicationAvailable notifies the registered observers that t is available.
func (s *ApplicationObserverSet) ApplicationAvailable(a *Application) {
	s.add(a, func(o interface{}) {
		o.(ApplicationObserver).ApplicationAvailable(a)
	})
}

// ApplicationUnavailable notifies the registered observers that t is unavailable.
func (s *ApplicationObserverSet) ApplicationUnavailable(a *Application) {
	s.remove(a, func(o interface{}) {
		o.(ApplicationObserver).ApplicationUnavailable(a)
	})
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
