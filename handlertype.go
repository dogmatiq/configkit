package configkit

import (
	"bytes"
	"fmt"

	"github.com/dogmatiq/configkit/message"
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
		return fmt.Errorf("invalid handler type: %s", string(t))
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

// IsConsumerOf returns true if handlers of type t can consume messages with the
// given role.
func (t HandlerType) IsConsumerOf(r message.Role) bool {
	return r.Is(t.Consumes()...)
}

// IsProducerOf returns true if handlers of type t can produce messages with the
// given role.
func (t HandlerType) IsProducerOf(r message.Role) bool {
	return r.Is(t.Produces()...)
}

// Consumes returns the roles of messages that can be consumed by handlers of
// this type.
func (t HandlerType) Consumes() []message.Role {
	t.MustValidate()

	switch t {
	case AggregateHandlerType:
		return []message.Role{message.CommandRole}
	case ProcessHandlerType:
		return []message.Role{message.EventRole, message.TimeoutRole}
	case IntegrationHandlerType:
		return []message.Role{message.CommandRole}
	default: // ProjectionHandlerType
		return []message.Role{message.EventRole}
	}
}

// Produces returns the roles of messages that can be produced by handlers of
// this type.
func (t HandlerType) Produces() []message.Role {
	t.MustValidate()

	switch t {
	case AggregateHandlerType:
		return []message.Role{message.EventRole}
	case ProcessHandlerType:
		return []message.Role{message.CommandRole, message.TimeoutRole}
	case IntegrationHandlerType:
		return []message.Role{message.EventRole}
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

// ConsumersOf returns the handler types that can consume messages with the
// given role.
func ConsumersOf(r message.Role) []HandlerType {
	r.MustValidate()

	switch r {
	case message.CommandRole:
		return []HandlerType{AggregateHandlerType, IntegrationHandlerType}
	case message.EventRole:
		return []HandlerType{ProcessHandlerType, ProjectionHandlerType}
	default: // message.TimeoutRole
		return []HandlerType{ProcessHandlerType}
	}
}

// ProducersOf returns the handler types that can produces messages with the
// given role.
func ProducersOf(r message.Role) []HandlerType {
	r.MustValidate()

	switch r {
	case message.CommandRole:
		return []HandlerType{ProcessHandlerType}
	case message.EventRole:
		return []HandlerType{AggregateHandlerType, IntegrationHandlerType}
	default: // message.TimeoutRole
		return []HandlerType{ProcessHandlerType}
	}
}
