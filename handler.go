package configkit

import (
	"fmt"
	"reflect"

	"github.com/dogmatiq/configkit/internal/validation"
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Handler is a specialization of the Entity interface for message handlers.
type Handler interface {
	Entity

	// HandlerType returns the type of handler.
	HandlerType() HandlerType

	// IsDisabled returns true if the handler is disabled.
	IsDisabled() bool
}

// RichHandler is a specialization of the Handler interface that exposes
// information about the Go types used to implement the Dogma application.
type RichHandler interface {
	RichEntity

	// HandlerType returns the type of handler.
	HandlerType() HandlerType

	// IsDisabled returns true if the handler is disabled.
	IsDisabled() bool
}

// IsHandlerEqual compares two handlers for equality.
//
// It returns true if both handlers:
//
//  1. have the same identity
//  2. produce and consume the same message types
//  3. are implemented using the same Go types
//
// Point 3. refers to the type used to implement the dogma.Aggregate,
// dogma.Process, dogma.Integration or dogma.Projection interface (not the type
// used to implement the configkit.Handler interface).
//
// This definition of equality relies on the fact that no single Go type can
// implement more than one Dogma handler interface because they all contain
// a Configure() method with different signatures.
func IsHandlerEqual(a, b Handler) bool {
	return a.Identity() == b.Identity() &&
		a.TypeName() == b.TypeName() &&
		a.HandlerType() == b.HandlerType() &&
		a.IsDisabled() == b.IsDisabled() &&
		a.MessageNames().IsEqual(b.MessageNames())
}

func configureRoutes[T dogma.Route](
	types *EntityMessages[message.Type],
	routes []T,
	handlerIdent Identity,
	handlerType reflect.Type,
) {
	if *types == nil {
		*types = EntityMessages[message.Type]{}
	}

	for _, route := range routes {
		switch route := any(route).(type) {
		case dogma.HandlesCommandRoute:
			configureConsumerRoute(*types, route.Type, "HandlesCommand", handlerIdent, handlerType)
		case dogma.RecordsEventRoute:
			configureProducerRoute(*types, route.Type, "RecordsEvent", handlerIdent, handlerType)
		case dogma.HandlesEventRoute:
			configureConsumerRoute(*types, route.Type, "HandlesEvent", handlerIdent, handlerType)
		case dogma.ExecutesCommandRoute:
			configureProducerRoute(*types, route.Type, "ExecutesCommand", handlerIdent, handlerType)
		case dogma.SchedulesTimeoutRoute:
			configureConsumerRoute(*types, route.Type, "SchedulesTimeout", handlerIdent, handlerType)
			configureProducerRoute(*types, route.Type, "SchedulesTimeout", handlerIdent, handlerType)
		default:
			panic(fmt.Sprintf("unsupported route type: %T", route))
		}
	}
}

func configureConsumerRoute(
	types EntityMessages[message.Type],
	messageType reflect.Type,
	routeFunc string,
	handlerIdent Identity,
	handlerType reflect.Type,
) {
	types.Update(
		message.TypeFromReflect(messageType),
		func(t message.Type, em *EntityMessage) {
			if em.IsConsumed {
				validation.Panicf(
					"%s is configured with multiple %s() routes for %s, should these refer to different message types?",
					handlerDisplayName(handlerIdent, handlerType),
					routeFunc,
					t,
				)
			}

			em.Kind = t.Kind()
			em.IsConsumed = true
		},
	)
}

func configureProducerRoute(
	types EntityMessages[message.Type],
	messageType reflect.Type,
	routeFunc string,
	handlerIdent Identity,
	handlerType reflect.Type,
) {
	types.Update(
		message.TypeFromReflect(messageType),
		func(t message.Type, em *EntityMessage) {
			if em.IsProduced {
				validation.Panicf(
					"%s is configured with multiple %s() routes for %s, should these refer to different message types?",
					handlerDisplayName(handlerIdent, handlerType),
					routeFunc,
					t,
				)
			}

			em.Kind = t.Kind()
			em.IsProduced = true
		},
	)
}

func handlerDisplayName(
	handlerIdent Identity,
	handlerType reflect.Type,
) string {
	s := handlerType.String()

	if !handlerIdent.IsZero() {
		s += " (" + handlerIdent.Name + ")"
	}

	return s
}

// mustHaveConsumerRoute panics if the handler is not configured to handle any
// messages of the given kind.
func mustHaveConsumerRoute(
	types *EntityMessages[message.Type],
	kind message.Kind,
	handlerIdent Identity,
	handlerType reflect.Type,
) {
	for _, k := range types.Consumed() {
		if k == kind {
			return
		}
	}

	validation.Panicf(
		`%s is not configured to handle any %ss, at least one Handles%s() route must be added within Configure()`,
		handlerDisplayName(handlerIdent, handlerType),
		kind,
		cases.Title(language.English).String(kind.String()),
	)
}

// mustHaveProducerRoute panics if the handler is not configured to produce any
// messages of the given kind.
func mustHaveProducerRoute(
	types *EntityMessages[message.Type],
	kind message.Kind,
	handlerIdent Identity,
	handlerType reflect.Type,
) {
	for _, k := range types.Produced() {
		if k == kind {
			return
		}
	}

	verb := message.MapKind(kind, "execute", "record", "schedule")
	routeFunc := message.MapKind(kind, "ExecutesCommand", "RecordsEvent", "SchedulesTimeout")

	validation.Panicf(
		`%s is not configured to %s any %ss, at least one %s() route must be added within Configure()`,
		handlerDisplayName(handlerIdent, handlerType),
		verb,
		kind,
		routeFunc,
	)
}
