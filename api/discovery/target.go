package discovery

import (
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
	// TargetAvailable is called when a target is becomes available.
	TargetAvailable(*Target)

	// TargetUnavailable is called when a target becomes unavailable.
	TargetUnavailable(*Target)
}

// TargetPublisher is an interface that allows target observers to be registered
// to receive notifications.
type TargetPublisher interface {
	// RegisterTargetObserver registers o to be notified when targets become
	// available and unavailable.
	RegisterTargetObserver(o TargetObserver)

	// UnregisterTargetObserver stops o from being notified when targets become
	// available and unavailable.
	UnregisterTargetObserver(o TargetObserver)
}
