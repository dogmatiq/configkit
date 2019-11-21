package configkit

import "github.com/dogmatiq/dogma"

// Integration is an interface that represents the configuration of a Dogma
// integration message handler.
type Integration interface {
	Handler

	// Handler returns the handler implementation.
	Handler() dogma.IntegrationMessageHandler
}

// RichIntegration is a specialization of the Integration interface that has
// access to the Go types used to implement the Dogma integration message
// handler.
type RichIntegration interface {
	RichHandler

	// Handler returns the handler implementation.
	Handler() dogma.IntegrationMessageHandler
}
