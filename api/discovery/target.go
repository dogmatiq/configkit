package discovery

import (
	"context"

	"google.golang.org/grpc"
)

// MetaData is a container for meta-data about a target.
type MetaData map[interface{}]interface{}

// Target represents some dialable gRPC target, typically a single gRPC server.
type Target struct {
	// Name is the target name used to dial the endpoint. The syntax is defined
	// in https://github.com/grpc/grpc/blob/master/doc/naming.md.
	Name string

	// Options is a set of grpc.DialOptions used when dialing this target.
	// The options must not include grpc.WithBlock().
	Options []grpc.DialOption

	// MetaData contains driver-specific meta-data about the target.
	MetaData MetaData
}

// TargetObserver is notified when config API targets are discovered.
type TargetObserver interface {
	// TargetAvailable is called when a target becomes available.
	TargetAvailable(*Target)

	// TargetUnavailable is called when a target becomes unavailable.
	TargetUnavailable(*Target)
}

// TargetObserverSet is a target observer that publishes to other observers.
//
// It implements both the TargetObserver and TargetPublisher interfaces.
type TargetObserverSet struct {
	observerSet
}

// NewTargetObserverSet registers the given observers with a new observer set
// and returns it.
func NewTargetObserverSet(observers ...TargetObserver) *TargetObserverSet {
	s := &TargetObserverSet{}

	for _, o := range observers {
		s.RegisterTargetObserver(o)
	}

	return s
}

// RegisterTargetObserver registers o to be notified when targets become
// available and unavailable.
func (s *TargetObserverSet) RegisterTargetObserver(o TargetObserver) {
	s.register(o, func(e interface{}) {
		o.TargetAvailable(e.(*Target))
	})
}

// UnregisterTargetObserver stops o from being notified when targets become
// available and unavailable.
func (s *TargetObserverSet) UnregisterTargetObserver(o TargetObserver) {
	s.unregister(o, func(e interface{}) {
		o.TargetUnavailable(e.(*Target))
	})
}

// TargetAvailable notifies the registered observers that t is available.
func (s *TargetObserverSet) TargetAvailable(t *Target) {
	s.add(t, func(o interface{}) {
		o.(TargetObserver).TargetAvailable(t)
	})
}

// TargetUnavailable notifies the registered observers that t is unavailable.
func (s *TargetObserverSet) TargetUnavailable(t *Target) {
	s.remove(t, func(o interface{}) {
		o.(TargetObserver).TargetUnavailable(t)
	})
}

// TargetTask is a function executed by a TargetExecutor.
type TargetTask func(context.Context, *Target)

// TargetExecutor is a TargetObserver that executes a function in a new
// goroutine whenever a target becomes available.
type TargetExecutor struct {
	executor

	// Task is the function to execute when a target becomes available.
	// The context is canceled when the target becomes unavailable.
	Task TargetTask

	// Parent is the parent context under which the function is called.
	// If it is nil, context.Background() is used.
	Parent context.Context
}

// TargetAvailable starts a new goroutine for the given target.
func (e *TargetExecutor) TargetAvailable(t *Target) {
	e.start(e.Parent, t, func(ctx context.Context) {
		e.Task(ctx, t)
	})
}

// TargetUnavailable cancels the context associated with any existing goroutine
// for the given target and waits for the goroutine to exit.
func (e *TargetExecutor) TargetUnavailable(t *Target) {
	e.stop(t)
}
