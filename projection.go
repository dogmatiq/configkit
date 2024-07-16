package configkit

import (
	"context"
	"reflect"

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

	// DeliveryPolicy returns the projection's delivery policy.
	DeliveryPolicy() dogma.ProjectionDeliveryPolicy
}

// FromProjection returns the configuration for a projection message handler.
//
// It panics if the handler is configured incorrectly. Use Recover() to convert
// configuration related panic values to errors.
func FromProjection(h dogma.ProjectionMessageHandler) RichProjection {
	cfg, c := fromProjection(h)
	c.mustValidate()
	return cfg
}

func fromProjection(h dogma.ProjectionMessageHandler) (*projection, *projectionConfigurer) {
	cfg := &projection{
		handler: handler{
			entity: entity{
				rt: reflect.TypeOf(h),
			},
		},
		impl:           h,
		deliveryPolicy: dogma.UnicastProjectionDeliveryPolicy{},
	}

	c := &projectionConfigurer{
		handlerConfigurer: handlerConfigurer{
			entityConfigurer: entityConfigurer{
				entity: &cfg.entity,
			},
			handler: &cfg.handler,
		},
		projection: cfg,
	}

	h.Configure(c)

	return cfg, c
}

// projection is an implementation of RichProjection.
type projection struct {
	handler

	impl           dogma.ProjectionMessageHandler
	deliveryPolicy dogma.ProjectionDeliveryPolicy
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

func (h *projection) DeliveryPolicy() dogma.ProjectionDeliveryPolicy {
	return h.deliveryPolicy
}
