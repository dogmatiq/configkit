package configkit

import (
	"context"
	"reflect"

	"github.com/dogmatiq/dogma"
)

// Process is an interface that represents the configuration of a Dogma process
// message handler.
type Process interface {
	Handler
}

// RichProcess is a specialization of Process that exposes information about the
// Go types used to implement the underlying Dogma handler.
type RichProcess interface {
	RichHandler

	// Handler returns the underlying message handler.
	Handler() dogma.ProcessMessageHandler
}

// FromProcess returns the configuration for a process message handler.
//
// It panics if the handler is configured incorrectly. Use Recover() to convert
// configuration related panic values to errors.
func FromProcess(h dogma.ProcessMessageHandler) RichProcess {
	cfg, c := fromProcess(h)
	c.mustValidate()
	return cfg
}

func fromProcess(h dogma.ProcessMessageHandler) (*richProcess, *processConfigurer) {
	cfg := &richProcess{
		handlerEntity: handlerEntity{
			entity: entity{
				rt: reflect.TypeOf(h),
			},
		},
		impl: h,
	}

	c := &processConfigurer{
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

// richProcess is the default implementation of [RichProcess].
type richProcess struct {
	handlerEntity

	impl dogma.ProcessMessageHandler
}

func (h *richProcess) AcceptVisitor(ctx context.Context, v Visitor) error {
	return v.VisitProcess(ctx, h)
}

func (h *richProcess) AcceptRichVisitor(ctx context.Context, v RichVisitor) error {
	return v.VisitRichProcess(ctx, h)
}

func (h *richProcess) HandlerType() HandlerType {
	return ProcessHandlerType
}

func (h *richProcess) Handler() dogma.ProcessMessageHandler {
	return h.impl
}
