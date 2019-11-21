package configkit

import "github.com/dogmatiq/dogma"

// Projection is an interface that represents the configuration of a Dogma
// projection message handler.
type Projection interface {
	Handler

	// Handler returns the handler implementation.
	Handler() dogma.ProjectionMessageHandler
}

// RichProjection is a specialization of the Projection interface that has
// access to the Go types used to implement the Dogma projection message
// handler.
type RichProjection interface {
	RichHandler

	// Handler returns the handler implementation.
	Handler() dogma.ProjectionMessageHandler
}
