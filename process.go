package configkit

import "github.com/dogmatiq/dogma"

// Process is an interface that represents the configuration of a Dogma process
// message handler.
type Process interface {
	Handler

	// Handler returns the handler implementation.
	Handler() dogma.ProcessMessageHandler
}

// RichProcess is a specialization of the Process interface that has access to
// the Go types used to implement the Dogma process message handler.
type RichProcess interface {
	RichHandler

	// Handler returns the handler implementation.
	Handler() dogma.ProcessMessageHandler
}
