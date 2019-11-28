package configkit

import (
	"context"
	"reflect"

	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
)

// Projection is an interface that represents the configuration of a Dogma
// projection message handler.
type Projection interface {
	Handler
}

// RichProjection is a specialization of Projection that exposes information
// about the Go types used to implement the underlying Dogma handler.
type RichProjection interface {
	RichHandler

	// Handler returns the underlying message handler.
	Handler() dogma.ProjectionMessageHandler
}

// FromProjection returns the configuration for a projection message handler.
//
// It panics if the handler is configured incorrectly. Use Recover() to convert
// configuration related panic values to errors.
func FromProjection(h dogma.ProjectionMessageHandler) RichProjection {
	cfg := &projection{
		entity: entity{
			rt: reflect.TypeOf(h),
		},
		impl: h,
	}

	c := &handlerConfigurer{
		entityConfigurer: entityConfigurer{
			entity: &cfg.entity,
		},
	}

	h.Configure(c)

	c.validate()
	c.mustConsume(message.EventRole)

	return cfg
}

// projection is an implementation of RichProjection.
type projection struct {
	entity

	impl dogma.ProjectionMessageHandler
}

func (h *projection) AcceptVisitor(ctx context.Context, v Visitor) error {
	return v.VisitProjection(ctx, h)
}

func (h *projection) AcceptRichVisitor(ctx context.Context, v RichVisitor) error {
	return v.VisitRichProjection(ctx, h)
}

func (h *projection) HandlerType() HandlerType {
	return ProjectionHandlerType
}

func (h *projection) Handler() dogma.ProjectionMessageHandler {
	return h.impl
}
