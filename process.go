package configkit

import (
	"context"
	"reflect"

	"github.com/dogmatiq/configkit/internal/typename/goreflect"
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
	cfg := fromProcessUnvalidated(h)
	cfg.mustValidate()
	return cfg
}

func fromProcessUnvalidated(h dogma.ProcessMessageHandler) *richProcess {
	cfg := &richProcess{handler: h}
	h.Configure(&processConfigurer{cfg})
	return cfg
}

// richProcess is the default implementation of [RichProcess].
type richProcess struct {
	ident      Identity
	types      EntityMessageTypes
	isDisabled bool
	handler    dogma.ProcessMessageHandler
}

func (h *richProcess) Identity() Identity {
	return h.ident
}

func (h *richProcess) MessageNames() EntityMessageNames {
	return h.types.asNames()
}

func (h *richProcess) MessageTypes() EntityMessageTypes {
	return h.types
}

func (h *richProcess) TypeName() string {
	return goreflect.NameOf(h.ReflectType())
}

func (h *richProcess) ReflectType() reflect.Type {
	return reflect.TypeOf(h.handler)
}

func (h *richProcess) IsDisabled() bool {
	return h.isDisabled
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
	return h.handler
}

func (h *richProcess) isConfigured() bool {
	return !h.ident.IsZero() ||
		h.types.Consumed.Len() != 0 ||
		h.types.Produced.Len() != 0
}

func (h *richProcess) mustValidate() {
	mustHaveValidIdentity(h.Identity(), h.ReflectType())
	mustHaveConsumerRoute(h.types, message.EventRole, h.Identity(), h.ReflectType())
	mustHaveProducerRoute(h.types, message.CommandRole, h.Identity(), h.ReflectType())
}

// processConfigurer is the default implementation of [dogma.ProcessConfigurer].
type processConfigurer struct {
	config *richProcess
}

func (c *processConfigurer) Identity(name, key string) {
	configureIdentity(&c.config.ident, name, key, c.config.ReflectType())
}

func (c *processConfigurer) Routes(routes ...dogma.ProcessRoute) {
	configureRoutes(&c.config.types, routes, c.config.ident, c.config.ReflectType())
}

func (c *processConfigurer) Disable(...dogma.DisableOption) {
	c.config.isDisabled = true
}
