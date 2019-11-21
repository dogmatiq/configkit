package configkit

import "github.com/dogmatiq/dogma"

// Aggregate is an interface that represents the configuration of a Dogma
// aggregate message handler.
type Aggregate interface {
	Handler

	// Handler returns the handler implementation.
	Handler() dogma.AggregateMessageHandler
}

// RichAggregate is a specialization of the Aggregate interface that has access
// to the Go types used to implement the Dogma aggregate message handler.
type RichAggregate interface {
	RichHandler

	// Handler returns the handler implementation.
	Handler() dogma.AggregateMessageHandler
}
