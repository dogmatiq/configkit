package configkit

import (
	"context"
	"reflect"

	"github.com/dogmatiq/configkit/message"
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
	cfg := &process{
		entity: entity{
			rt: reflect.TypeOf(h),
		},
		impl: h,
	}

	c := &handlerConfigurer{
		target: &cfg.entity,
	}

	h.Configure(c)

	c.validate()
	c.mustConsume(message.EventRole)
	c.mustProduce(message.CommandRole)

	return cfg
}

// process is an implementation of RichProcess.
type process struct {
	entity

	impl dogma.ProcessMessageHandler
}

func (h *process) AcceptVisitor(ctx context.Context, v Visitor) error {
	return v.VisitProcess(ctx, h)
}

func (h *process) AcceptRichVisitor(ctx context.Context, v RichVisitor) error {
	return v.VisitRichProcess(ctx, h)
}

func (h *process) HandlerType() HandlerType {
	return ProcessHandlerType
}

func (h *process) Handler() dogma.ProcessMessageHandler {
	return h.impl
}
