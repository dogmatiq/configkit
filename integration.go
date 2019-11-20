package configkit

import "github.com/dogmatiq/dogma"

type integration interface {
	// Handler returns the handler implementation.
	Handler() dogma.IntegrationMessageHandler
}

// Integration is an interface that represents the configuration of a Dogma
// integration message handler.
type Integration interface {
	Handler
	integration
}

// PortableIntegration is a specialization of the Integration interface that
// does not require the Go types used to implement the Dogma integration message
// handler.
type PortableIntegration interface {
	PortableHandler
	integration
}

// RichIntegration is a specialization of the Integration interface that has
// access to the Go types used to implement the Dogma integration message
// handler.
type RichIntegration interface {
	RichHandler
	integration
}
