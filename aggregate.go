package configkit

import "github.com/dogmatiq/dogma"

type aggregate interface {
	// Handler returns the handler implementation.
	Handler() dogma.AggregateMessageHandler
}

// Aggregate is an interface that represents the configuration of a Dogma
// aggregate message handler.
type Aggregate interface {
	Handler
	aggregate
}

// PortableAggregate is a specialization of the Aggregate interface that does
// not require the Go types used to implement the Dogma aggregate message
// handler.
type PortableAggregate interface {
	PortableHandler
	aggregate
}

// RichAggregate is a specialization of the Aggregate interface that has access
// to the Go types used to implement the Dogma aggregate message handler.
type RichAggregate interface {
	RichHandler
	aggregate
}
