package configkit

import (
	"context"
	"reflect"

	"github.com/dogmatiq/configkit/message"
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
	cfg := &integration{
		handler: handler{
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
			handler: &cfg.handler,
		},
	}

	h.Configure(c)

	c.validate()
	c.mustConsume(message.CommandRole)

	return cfg
}

// integration is an implementation of RichIntegration.
type integration struct {
	handler

	impl dogma.IntegrationMessageHandler
}

func (h *integration) AcceptVisitor(ctx context.Context, v Visitor) error {
	return v.VisitIntegration(ctx, h)
}

func (h *integration) AcceptRichVisitor(ctx context.Context, v RichVisitor) error {
	return v.VisitRichIntegration(ctx, h)
}

func (h *integration) HandlerType() HandlerType {
	return IntegrationHandlerType
}

func (h *integration) Handler() dogma.IntegrationMessageHandler {
	return h.impl
}
