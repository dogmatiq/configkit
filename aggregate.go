package configkit

import (
	"context"
	"reflect"

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

// FromAggregate returns the configuration for an aggregate message handler.
//
// It panics if the handler is configured incorrectly. Use Recover() to convert
// configuration related panic values to errors.
func FromAggregate(h dogma.AggregateMessageHandler) RichAggregate {
	cfg, c := fromAggregate(h)
	c.mustValidate()
	return cfg
}

func fromAggregate(h dogma.AggregateMessageHandler) (*aggregate, *aggregateConfigurer) {
	cfg := &aggregate{
		handler: handler{
			entity: entity{
				rt: reflect.TypeOf(h),
			},
		},
		impl: h,
	}

	c := &aggregateConfigurer{
		handlerConfigurer: handlerConfigurer{
			entityConfigurer: entityConfigurer{
				entity: &cfg.entity,
			},
			handler: &cfg.handler,
		},
	}

	h.Configure(c)

	return cfg, c
}

// aggregate is an implementation of RichAggregate.
type aggregate struct {
	handler

	impl dogma.AggregateMessageHandler
}

func (h *aggregate) AcceptVisitor(ctx context.Context, v Visitor) error {
	return v.VisitAggregate(ctx, h)
}

func (h *aggregate) AcceptRichVisitor(ctx context.Context, v RichVisitor) error {
	return v.VisitRichAggregate(ctx, h)
}

func (h *aggregate) HandlerType() HandlerType {
	return AggregateHandlerType
}

func (h *aggregate) Handler() dogma.AggregateMessageHandler {
	return h.impl
}
