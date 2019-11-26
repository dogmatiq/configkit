package configkit

import (
	"github.com/dogmatiq/dogma"
)

// Aggregate is an interface that represents the configuration of a Dogma
// aggregate message handler.
type Aggregate interface {
	Handler
}

// RichAggregate is a specialization of Aggregate that exposes information
// about the Go types used to implement the underlying Dogma handler.
type RichAggregate interface {
	RichHandler

	// Handler returns the underlying message handler.
	Handler() dogma.AggregateMessageHandler
}
