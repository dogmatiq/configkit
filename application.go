package configkit

import (
	"context"
	"reflect"

	"github.com/dogmatiq/configkit/internal/typename/goreflect"
	"github.com/dogmatiq/configkit/internal/validation"
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
)

// Application is an interface that represents the configuration of a Dogma
// application.
type Application interface {
	Entity

	// Handlers returns the handlers within this application.
	Handlers() HandlerSet
}

// RichApplication is a specialization of Application that exposes information
// about the Go types used to implement the Dogma application.
type RichApplication interface {
	RichEntity

	// Handlers returns the handlers within this application.
	Handlers() HandlerSet

	// RichHandlers returns the handlers within this application.
	RichHandlers() RichHandlerSet

	// Application returns the underlying application.
	Application() dogma.Application
}

// FromApplication returns the configuration for an application.
//
// It panics if the application is configured incorrectly. Use Recover() to
// convert configuration related panic values to errors.
func FromApplication(a dogma.Application) RichApplication {
	cfg := &richApplication{app: a}
	a.Configure(&applicationConfigurer{cfg})

	mustHaveValidIdentity(
		cfg.Identity(),
		cfg.ReflectType(),
	)

	return cfg
}

// IsApplicationEqual compares two applications for equality.
//
// It returns true if both applications:
//
//  1. have the same identity
//  2. produce and consume the same messages, with the same roles
//  3. are implemented using the same Go types
//  4. contain equivalent handlers
//
// Point 3. refers to the type used to implement the dogma.Application interface
// (not the type used to implement the configkit.Application interface).
func IsApplicationEqual(a, b Application) bool {
	return a.Identity() == b.Identity() &&
		a.TypeName() == b.TypeName() &&
		a.MessageNames().IsEqual(b.MessageNames()) &&
		a.Handlers().IsEqual(b.Handlers())
}

// richApplication is the default implementation of [RichApplication].
type richApplication struct {
	ident    Identity
	types    EntityMessages[message.Type]
	handlers RichHandlerSet
	app      dogma.Application
}

func (a *richApplication) Identity() Identity {
	return a.ident
}

func (a *richApplication) MessageNames() EntityMessages[message.Name] {
	return asMessageNames(a.types)
}

func (a *richApplication) MessageTypes() EntityMessages[message.Type] {
	return a.types
}

func (a *richApplication) TypeName() string {
	return goreflect.NameOf(a.ReflectType())
}

func (a *richApplication) ReflectType() reflect.Type {
	return reflect.TypeOf(a.app)
}

func (a *richApplication) AcceptVisitor(ctx context.Context, v Visitor) error {
	return v.VisitApplication(ctx, a)
}

func (a *richApplication) AcceptRichVisitor(ctx context.Context, v RichVisitor) error {
	return v.VisitRichApplication(ctx, a)
}

func (a *richApplication) Handlers() HandlerSet {
	return a.handlers.asHandlerSet()
}

func (a *richApplication) RichHandlers() RichHandlerSet {
	return a.handlers
}

func (a *richApplication) Application() dogma.Application {
	return a.app
}

// applicationConfigurer is the default implementation of
// [dogma.ApplicationConfigurer].
type applicationConfigurer struct {
	config *richApplication
}

func (c *applicationConfigurer) Identity(name, key string) {
	if h, ok := c.config.handlers.ByKey(key); ok {
		validation.Panicf(
			`%s can not use the application key "%s", because it is already used by %s`,
			c.config.ReflectType(),
			key,
			h.ReflectType(),
		)
	}

	configureIdentity(&c.config.ident, name, key, c.config.ReflectType())
}

func (c *applicationConfigurer) RegisterAggregate(h dogma.AggregateMessageHandler, _ ...dogma.RegisterAggregateOption) {
	c.registerIfConfigured(fromAggregateUnvalidated(h))
}

func (c *applicationConfigurer) RegisterProcess(h dogma.ProcessMessageHandler, _ ...dogma.RegisterProcessOption) {
	c.registerIfConfigured(fromProcessUnvalidated(h))
}

func (c *applicationConfigurer) RegisterIntegration(h dogma.IntegrationMessageHandler, _ ...dogma.RegisterIntegrationOption) {
	c.registerIfConfigured(fromIntegrationUnvalidated(h))
}

func (c *applicationConfigurer) RegisterProjection(h dogma.ProjectionMessageHandler, _ ...dogma.RegisterProjectionOption) {
	c.registerIfConfigured(fromProjectionUnvalidated(h))
}

type validatableHandler interface {
	RichHandler

	// isConfigured returns true if the handler has been configured in any way
	// beyond being disabled, even if the configuration is invalid.
	isConfigured() bool

	// mustValidate panics if the handler is not configured correctly.
	mustValidate()
}

func (c *applicationConfigurer) registerIfConfigured(
	h validatableHandler,
) {
	if h.IsDisabled() && !h.isConfigured() {
		return
	}

	h.mustValidate()

	c.guardAgainstConflictingIdentities(h)
	c.guardAgainstConflictingRoutes(h)

	if c.config.handlers == nil {
		c.config.handlers = RichHandlerSet{}
		c.config.types = EntityMessages[message.Type]{}
	}

	c.config.handlers.Add(h)
	c.config.types.merge(h.MessageTypes())
}

// guardAgainstConflictingIdentities panics if h's identity conflicts with the
// application or any other handlers.
func (c *applicationConfigurer) guardAgainstConflictingIdentities(h RichHandler) {
	appIdent := c.config.Identity()
	handlerIdent := h.Identity()

	if handlerIdent.Key == appIdent.Key {
		validation.Panicf(
			`%s can not use the handler key "%s", because it is already used by %s`,
			h.ReflectType(),
			handlerIdent.Key,
			c.config.ReflectType(),
		)
	}

	if x, ok := c.config.handlers.ByName(handlerIdent.Name); ok {
		validation.Panicf(
			`%s can not use the handler name "%s", because it is already used by %s`,
			h.ReflectType(),
			handlerIdent.Name,
			x.ReflectType(),
		)
	}

	if x, ok := c.config.handlers.ByKey(handlerIdent.Key); ok {
		validation.Panicf(
			`%s can not use the handler key "%s", because it is already used by %s`,
			h.ReflectType(),
			handlerIdent.Key,
			x.ReflectType(),
		)
	}
}

// guardAgainstConflictingRoutes panics if an h consumes the same commands or
// produces the same events as some existing handler.
func (c *applicationConfigurer) guardAgainstConflictingRoutes(h RichHandler) {
	for mt, em := range h.MessageTypes() {
		if em.Kind == message.CommandKind && em.IsConsumed {
			for _, x := range c.config.handlers.ConsumersOf(mt) {
				validation.Panicf(
					`%s (%s) can not handle %s commands because they are already configured to be handled by %s (%s)`,
					h.ReflectType(),
					h.Identity().Name,
					mt,
					x.ReflectType(),
					x.Identity().Name,
				)
			}
		}

		if em.Kind == message.EventKind && em.IsProduced {
			for _, x := range c.config.handlers.ProducersOf(mt) {
				validation.Panicf(
					`%s (%s) can not record %s events because they are already configured to be recorded by %s (%s)`,
					h.ReflectType(),
					h.Identity().Name,
					mt,
					x.ReflectType(),
					x.Identity().Name,
				)
			}
		}
	}
}
