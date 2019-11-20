package configkit

import "github.com/dogmatiq/dogma"

type process interface {
	// Handler returns the handler implementation.
	Handler() dogma.ProcessMessageHandler
}

// Process is an interface that represents the configuration of a Dogma process
// message handler.
type Process interface {
	Handler
	process
}

// PortableProcess is a specialization of the Process interface that does not
// require the Go types used to implement the Dogma process message handler.
type PortableProcess interface {
	PortableHandler
	process
}

// RichProcess is a specialization of the Process interface that has access to
// the Go types used to implement the Dogma process message handler.
type RichProcess interface {
	RichHandler
	process
}
