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

func fromProjection(h dogma.ProjectionMessageHandler) (*richProjection, *projectionConfigurer) {
	cfg := &richProjection{
		handlerEntity: handlerEntity{
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
			handler: &cfg.handlerEntity,
		},
		projection: cfg,
	}

	h.Configure(c)

	return cfg, c
}

// richProjection is the default implementation of [RichProjection].
type richProjection struct {
	handlerEntity

	impl           dogma.ProjectionMessageHandler
	deliveryPolicy dogma.ProjectionDeliveryPolicy
}

func (h *richProjection) AcceptVisitor(ctx context.Context, v Visitor) error {
	return v.VisitProjection(ctx, h)
}

func (h *richProjection) AcceptRichVisitor(ctx context.Context, v RichVisitor) error {
	return v.VisitRichProjection(ctx, h)
}

func (h *richProjection) HandlerType() HandlerType {
	return ProjectionHandlerType
}

func (h *richProjection) Handler() dogma.ProjectionMessageHandler {
	return h.impl
}

func (h *richProjection) DeliveryPolicy() dogma.ProjectionDeliveryPolicy {
	return h.deliveryPolicy
}
