package configkit

import (
	"context"
	"reflect"

	"github.com/dogmatiq/configkit/internal/typename/goreflect"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/message"
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

// FromAggregate returns the configuration for an aggregate message handler.
//
// It panics if the handler is configured incorrectly. Use Recover() to convert
// configuration related panic values to errors.
func FromAggregate(h dogma.AggregateMessageHandler) RichAggregate {
	cfg := fromAggregateUnvalidated(h)
	cfg.mustValidate()
	return cfg
}

func fromAggregateUnvalidated(h dogma.AggregateMessageHandler) *richAggregate {
	cfg := &richAggregate{handler: h}
	h.Configure(&aggregateConfigurer{cfg})
	return cfg
}

// richAggregate the default implementation of [RichAggregate].
type richAggregate struct {
	ident      Identity
	types      EntityMessages[message.Type]
	isDisabled bool
	handler    dogma.AggregateMessageHandler
}

func (h *richAggregate) Identity() Identity {
	return h.ident
}

func (h *richAggregate) MessageNames() EntityMessages[message.Name] {
	return asMessageNames(h.types)
}

func (h *richAggregate) MessageTypes() EntityMessages[message.Type] {
	return h.types
}

func (h *richAggregate) TypeName() string {
	return goreflect.NameOf(h.ReflectType())
}

func (h *richAggregate) ReflectType() reflect.Type {
	return reflect.TypeOf(h.handler)
}

func (h *richAggregate) IsDisabled() bool {
	return h.isDisabled
}

func (h *richAggregate) AcceptVisitor(ctx context.Context, v Visitor) error {
	return v.VisitAggregate(ctx, h)
}

func (h *richAggregate) AcceptRichVisitor(ctx context.Context, v RichVisitor) error {
	return v.VisitRichAggregate(ctx, h)
}

func (h *richAggregate) HandlerType() HandlerType {
	return AggregateHandlerType
}

func (h *richAggregate) Handler() dogma.AggregateMessageHandler {
	return h.handler
}

func (h *richAggregate) isConfigured() bool {
	return !h.ident.IsZero() || len(h.types) != 0
}

func (h *richAggregate) mustValidate() {
	mustHaveValidIdentity(h.Identity(), h.ReflectType())
	mustHaveConsumerRoute(&h.types, message.CommandKind, h.Identity(), h.ReflectType())
	mustHaveProducerRoute(&h.types, message.EventKind, h.Identity(), h.ReflectType())
}

// aggregateConfigurer is the default implementation of
// [dogma.AggregateConfigurer].
type aggregateConfigurer struct {
	config *richAggregate
}

func (c *aggregateConfigurer) Identity(name, key string) {
	configureIdentity(&c.config.ident, name, key, c.config.ReflectType())
}

func (c *aggregateConfigurer) Routes(routes ...dogma.AggregateRoute) {
	configureRoutes(&c.config.types, routes, c.config.ident, c.config.ReflectType())
}

func (c *aggregateConfigurer) Disable(...dogma.DisableOption) {
	c.config.isDisabled = true
}
