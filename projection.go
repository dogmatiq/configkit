package configkit

import "github.com/dogmatiq/dogma"

type projection interface {
	// Handler returns the handler implementation.
	Handler() dogma.ProjectionMessageHandler
}

// Projection is an interface that represents the configuration of a Dogma
// projection message handler.
type Projection interface {
	Handler
	projection
}

// PortableProjection is a specialization of the Projection interface that does
// not require the Go types used to implement the Dogma projection message
// handler.
type PortableProjection interface {
	PortableHandler
	projection
}

// RichProjection is a specialization of the Projection interface that has
// access to the Go types used to implement the Dogma projection message
// handler.
type RichProjection interface {
	RichHandler
	projection
}
