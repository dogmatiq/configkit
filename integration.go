package configkit

import (
	"context"
	"reflect"

	"github.com/dogmatiq/dogma"
)

// Integration is an interface that represents the configuration of a Dogma
// integration message handler.
type Integration interface {
	Handler
}

// RichIntegration is a specialization of Integration that exposes information
// about the Go types used to implement the underlying Dogma handler.
type RichIntegration interface {
	RichHandler

	// Handler returns the underlying message handler.
	Handler() dogma.IntegrationMessageHandler
}

// FromIntegration returns the configuration for an integration message handler.
//
// It panics if the handler is configured incorrectly. Use Recover() to convert
// configuration related panic values to errors.
func FromIntegration(h dogma.IntegrationMessageHandler) RichIntegration {
	cfg, c := fromIntegration(h)
	c.mustValidate()
	return cfg
}

func fromIntegration(h dogma.IntegrationMessageHandler) (*richIntegration, *integrationConfigurer) {
	cfg := &richIntegration{
		handlerEntity: handlerEntity{
			entity: entity{
				rt: reflect.TypeOf(h),
			},
		},
		impl: h,
	}

	c := &integrationConfigurer{
		handlerConfigurer: handlerConfigurer{
			entityConfigurer: entityConfigurer{
				entity: &cfg.entity,
			},
			handler: &cfg.handlerEntity,
		},
	}

	h.Configure(c)

	return cfg, c
}

// richIntegration is the default implementation of [RichIntegration].
type richIntegration struct {
	handlerEntity

	impl dogma.IntegrationMessageHandler
}

func (h *richIntegration) AcceptVisitor(ctx context.Context, v Visitor) error {
	return v.VisitIntegration(ctx, h)
}

func (h *richIntegration) AcceptRichVisitor(ctx context.Context, v RichVisitor) error {
	return v.VisitRichIntegration(ctx, h)
}

func (h *richIntegration) HandlerType() HandlerType {
	return IntegrationHandlerType
}

func (h *richIntegration) Handler() dogma.IntegrationMessageHandler {
	return h.impl
}
