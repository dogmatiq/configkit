package configkit

import (
	"context"
	"reflect"

	"github.com/dogmatiq/configkit/internal/typename/goreflect"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/message"
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
	cfg := fromIntegrationUnvalidated(h)
	cfg.mustValidate()
	return cfg
}

func fromIntegrationUnvalidated(h dogma.IntegrationMessageHandler) *richIntegration {
	cfg := &richIntegration{handler: h}
	h.Configure(&integrationConfigurer{cfg})
	return cfg
}

// richIntegration the default implementation of [RichIntegration].
type richIntegration struct {
	ident      Identity
	types      EntityMessages[message.Type]
	isDisabled bool
	handler    dogma.IntegrationMessageHandler
}

func (h *richIntegration) Identity() Identity {
	return h.ident
}

func (h *richIntegration) MessageNames() EntityMessages[message.Name] {
	return asMessageNames(h.types)
}

func (h *richIntegration) MessageTypes() EntityMessages[message.Type] {
	return h.types
}

func (h *richIntegration) TypeName() string {
	return goreflect.NameOf(h.ReflectType())
}

func (h *richIntegration) ReflectType() reflect.Type {
	return reflect.TypeOf(h.handler)
}

func (h *richIntegration) IsDisabled() bool {
	return h.isDisabled
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
	return h.handler
}

func (h *richIntegration) isConfigured() bool {
	return !h.ident.IsZero() || len(h.types) != 0
}

func (h *richIntegration) mustValidate() {
	mustHaveValidIdentity(h.Identity(), h.ReflectType())
	mustHaveConsumerRoute(&h.types, message.CommandKind, h.Identity(), h.ReflectType())
}

// integrationConfigurer is the default implementation of
// [dogma.IntegrationConfigurer].
type integrationConfigurer struct {
	config *richIntegration
}

func (c *integrationConfigurer) Identity(name, key string) {
	configureIdentity(&c.config.ident, name, key, c.config.ReflectType())
}

func (c *integrationConfigurer) Routes(routes ...dogma.IntegrationRoute) {
	configureRoutes(&c.config.types, routes, c.config.ident, c.config.ReflectType())
}

func (c *integrationConfigurer) Disable(...dogma.DisableOption) {
	c.config.isDisabled = true
}
