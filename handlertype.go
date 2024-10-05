package configkit

import (
	"bytes"
	"fmt"
	"slices"

	"github.com/dogmatiq/configkit/internal/validation"
	"github.com/dogmatiq/enginekit/message"
)

// HandlerType is an enumeration of the types of handlers.
type HandlerType string

const (
	// AggregateHandlerType is the handler type for dogma.AggregateMessageHandler.
	AggregateHandlerType HandlerType = "aggregate"

	// ProcessHandlerType is the handler type for dogma.ProcessMessageHandler.
	ProcessHandlerType HandlerType = "process"

	// IntegrationHandlerType is the handler type for dogma.IntegrationMessageHandler.
	IntegrationHandlerType HandlerType = "integration"

	// ProjectionHandlerType is the handler type for dogma.ProjectionMessageHandler.
	ProjectionHandlerType HandlerType = "projection"
)

// HandlerTypes is a slice of the valid handler types.
var HandlerTypes = []HandlerType{
	AggregateHandlerType,
	ProcessHandlerType,
	IntegrationHandlerType,
	ProjectionHandlerType,
}

const (
	aggregateHandlerTypeShortString   = "AGG"
	processHandlerTypeShortString     = "PRC"
	integrationHandlerTypeShortString = "INT"
	projectionHandlerTypeShortString  = "PRJ"
)

var (
	aggregateHandlerTypeBinary   = []byte{'A'}
	processHandlerTypeBinary     = []byte{'P'}
	integrationHandlerTypeBinary = []byte{'I'}
	projectionHandlerTypeBinary  = []byte{'R'}
)

// Validate returns an error if t is not a valid handler type.
func (t HandlerType) Validate() error {
	switch t {
	case AggregateHandlerType,
		ProcessHandlerType,
		IntegrationHandlerType,
		ProjectionHandlerType:
		return nil
	default:
		return validation.Errorf("invalid handler type: %s", string(t))
	}
}

// MustValidate panics if t is not a valid handler type.
func (t HandlerType) MustValidate() {
	if err := t.Validate(); err != nil {
		panic(err)
	}
}

// Is returns true if t is one of the given types.
func (t HandlerType) Is(types ...HandlerType) bool {
	t.MustValidate()

	for _, x := range types {
		x.MustValidate()

		if t == x {
			return true
		}
	}

	return false
}

// MustBe panics if t is not one of the given types.
func (t HandlerType) MustBe(types ...HandlerType) {
	if !t.Is(types...) {
		panic("unexpected handler type: " + t.String())
	}
}

// MustNotBe panics if t is one of the given types.
func (t HandlerType) MustNotBe(types ...HandlerType) {
	if t.Is(types...) {
		panic("unexpected handler type: " + t.String())
	}
}

// IsConsumerOf returns true if handlers of type t can consume messages of the
// given kind.
func (t HandlerType) IsConsumerOf(k message.Kind) bool {
	return slices.Contains(t.Consumes(), k)
}

// IsProducerOf returns true if handlers of type t can produce messages of the
// given kind.
func (t HandlerType) IsProducerOf(k message.Kind) bool {
	return slices.Contains(t.Produces(), k)
}

// Consumes returns the kind of messages that can be consumed by handlers of
// type t.
func (t HandlerType) Consumes() []message.Kind {
	t.MustValidate()

	switch t {
	case AggregateHandlerType:
		return []message.Kind{message.CommandKind}
	case ProcessHandlerType:
		return []message.Kind{message.EventKind, message.TimeoutKind}
	case IntegrationHandlerType:
		return []message.Kind{message.CommandKind}
	default: // ProjectionHandlerType
		return []message.Kind{message.EventKind}
	}
}

// Produces returns the kinds of messages that can be produced by handlers of
// type t.
func (t HandlerType) Produces() []message.Kind {
	t.MustValidate()

	switch t {
	case AggregateHandlerType:
		return []message.Kind{message.EventKind}
	case ProcessHandlerType:
		return []message.Kind{message.CommandKind, message.TimeoutKind}
	case IntegrationHandlerType:
		return []message.Kind{message.EventKind}
	default: // ProjectionHandlerType
		return nil
	}
}

// ShortString returns a short (3-character) representation of the handler type.
func (t HandlerType) ShortString() string {
	t.MustValidate()

	switch t {
	case AggregateHandlerType:
		return aggregateHandlerTypeShortString
	case ProcessHandlerType:
		return processHandlerTypeShortString
	case IntegrationHandlerType:
		return integrationHandlerTypeShortString
	default: // ProjectionHandlerType
		return projectionHandlerTypeShortString
	}
}

// String returns a string representation of the handler type.
func (t HandlerType) String() string {
	if err := t.Validate(); err != nil {
		return "<" + err.Error() + ">"
	}

	return string(t)
}

// MarshalText returns a UTF-8 representation of the handler type.
func (t HandlerType) MarshalText() ([]byte, error) {
	return []byte(t), t.Validate()
}

// UnmarshalText unmarshals a type from its UTF-8 representation.
func (t *HandlerType) UnmarshalText(text []byte) error {
	x := HandlerType(text)

	if x.Validate() != nil {
		return fmt.Errorf("invalid text representation of handler type: %s", text)
	}

	*t = x
	return nil
}

// MarshalBinary returns a binary representation of the handler type.
func (t HandlerType) MarshalBinary() ([]byte, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}

	switch t {
	case AggregateHandlerType:
		return aggregateHandlerTypeBinary, nil
	case ProcessHandlerType:
		return processHandlerTypeBinary, nil
	case IntegrationHandlerType:
		return integrationHandlerTypeBinary, nil
	default: // ProjectionHandlerType
		return projectionHandlerTypeBinary, nil
	}
}

// UnmarshalBinary unmarshals a type from its binary representation.
func (t *HandlerType) UnmarshalBinary(data []byte) error {
	if bytes.Equal(data, aggregateHandlerTypeBinary) {
		*t = AggregateHandlerType
	} else if bytes.Equal(data, processHandlerTypeBinary) {
		*t = ProcessHandlerType
	} else if bytes.Equal(data, integrationHandlerTypeBinary) {
		*t = IntegrationHandlerType
	} else if bytes.Equal(data, projectionHandlerTypeBinary) {
		*t = ProjectionHandlerType
	} else {
		return fmt.Errorf("invalid binary representation of handler type: %s", data)
	}

	return nil
}

// ConsumersOf returns the handler types that can consume messages of kind k.
func ConsumersOf(k message.Kind) []HandlerType {
	return message.MapByKind(
		k,
		[]HandlerType{AggregateHandlerType, IntegrationHandlerType},
		[]HandlerType{ProcessHandlerType, ProjectionHandlerType},
		[]HandlerType{ProcessHandlerType},
	)
}

// ProducersOf returns the handler types that can produces messages of kind k.
func ProducersOf(k message.Kind) []HandlerType {
	return message.MapByKind(
		k,
		[]HandlerType{ProcessHandlerType},
		[]HandlerType{AggregateHandlerType, IntegrationHandlerType},
		[]HandlerType{ProcessHandlerType},
	)
}
