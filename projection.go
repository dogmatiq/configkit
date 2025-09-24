package configkit

import (
	"context"
	"reflect"

	"github.com/dogmatiq/configkit/internal/typename/goreflect"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/enginekit/message"
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
	cfg := fromProjectionUnvalidated(h)
	cfg.mustValidate()
	return cfg
}

func fromProjectionUnvalidated(h dogma.ProjectionMessageHandler) *richProjection {
	cfg := &richProjection{handler: h}
	h.Configure(&projectionConfigurer{cfg})
	return cfg
}

// richProjection is an implementation of RichProjection.
type richProjection struct {
	ident      Identity
	types      EntityMessages[message.Type]
	isDisabled bool
	handler    dogma.ProjectionMessageHandler
}

func (h *richProjection) Identity() Identity {
	return h.ident
}

func (h *richProjection) MessageNames() EntityMessages[message.Name] {
	return asMessageNames(h.types)
}

func (h *richProjection) MessageTypes() EntityMessages[message.Type] {
	return h.types
}

func (h *richProjection) TypeName() string {
	return goreflect.NameOf(h.ReflectType())
}

func (h *richProjection) ReflectType() reflect.Type {
	return reflect.TypeOf(h.handler)
}

func (h *richProjection) IsDisabled() bool {
	return h.isDisabled
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
	return h.handler
}

func (h *richProjection) isConfigured() bool {
	return !h.ident.IsZero() || len(h.types) != 0
}

func (h *richProjection) mustValidate() {
	mustHaveValidIdentity(h.Identity(), h.ReflectType())
	mustHaveConsumerRoute(&h.types, message.EventKind, h.Identity(), h.ReflectType())
}

// projectionConfigurer is the default implementation of
// [dogma.ProjectionConfigurer].
type projectionConfigurer struct {
	config *richProjection
}

func (c *projectionConfigurer) Identity(name, key string) {
	configureIdentity(&c.config.ident, name, key, c.config.ReflectType())
}

func (c *projectionConfigurer) Routes(routes ...dogma.ProjectionRoute) {
	configureRoutes(&c.config.types, routes, c.config.ident, c.config.ReflectType())
}

func (c *projectionConfigurer) Disable(...dogma.DisableOption) {
	c.config.isDisabled = true
}
