package configkit

import (
	"github.com/dogmatiq/dogma"
)

// Integration is an interface that represents the configuration of a Dogma
// integration message handler.
type Integration interface {
	Handler
}

// RichIntegration is a specialization of Integration that exposes information
// about the Go types used to implement the underlying Dogma handler.
type RichIntegration interface {
	RichHandler

	// Handler returns the underlying message handler.
	Handler() dogma.IntegrationMessageHandler
}
