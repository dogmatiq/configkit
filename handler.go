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
//  2. produce and consume the same messages, with the same roles
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
	types *EntityMessageTypes,
	routes []T,
	handlerIdent Identity,
	handlerType reflect.Type,
) {
	for _, route := range routes {
		switch route := any(route).(type) {
		case dogma.HandlesCommandRoute:
			configureConsumerRoute(types, route.Type, message.CommandRole, "HandlesCommand", handlerIdent, handlerType)
		case dogma.RecordsEventRoute:
			configureProducerRoute(types, route.Type, message.EventRole, "RecordsEvent", handlerIdent, handlerType)
		case dogma.HandlesEventRoute:
			configureConsumerRoute(types, route.Type, message.EventRole, "HandlesEvent", handlerIdent, handlerType)
		case dogma.ExecutesCommandRoute:
			configureProducerRoute(types, route.Type, message.CommandRole, "ExecutesCommand", handlerIdent, handlerType)
		case dogma.SchedulesTimeoutRoute:
			configureConsumerRoute(types, route.Type, message.TimeoutRole, "SchedulesTimeout", handlerIdent, handlerType)
			configureProducerRoute(types, route.Type, message.TimeoutRole, "SchedulesTimeout", handlerIdent, handlerType)
		default:
			panic(fmt.Sprintf("unsupported route type: %T", route))
		}
	}
}

func configureConsumerRoute(
	types *EntityMessageTypes,
	messageType reflect.Type,
	role message.Role,
	routeFunc string,
	handlerIdent Identity,
	handlerType reflect.Type,
) {
	t := message.TypeFromReflect(messageType)

	guardAgainstConflictingRoles(
		types,
		t,
		role,
		handlerIdent,
		handlerType,
	)

	if types.Consumed.Has(t) {
		validation.Panicf(
			"%s is configured with multiple %s() routes for %s, should these refer to different message types?",
			handlerDisplayName(handlerIdent, handlerType),
			routeFunc,
			t,
		)
	}

	if types.Consumed == nil {
		types.Consumed = message.TypeRoles{}
	}
	types.Consumed.Add(t, role)
}

func configureProducerRoute(
	types *EntityMessageTypes,
	messageType reflect.Type,
	role message.Role,
	routeFunc string,
	handlerIdent Identity,
	handlerType reflect.Type,
) {
	t := message.TypeFromReflect(messageType)

	guardAgainstConflictingRoles(
		types,
		t,
		role,
		handlerIdent,
		handlerType,
	)

	if types.Produced.Has(t) {
		validation.Panicf(
			"%s is configured with multiple %s() routes for %s, should these refer to different message types?",
			handlerDisplayName(handlerIdent, handlerType),
			routeFunc,
			t,
		)
	}

	if types.Produced == nil {
		types.Produced = message.TypeRoles{}
	}
	types.Produced.Add(t, role)
}

func guardAgainstConflictingRoles(
	types *EntityMessageTypes,
	messageType message.Type,
	proposed message.Role,
	handlerIdent Identity,
	handlerType reflect.Type,
) {
	existing, ok := types.RoleOf(messageType)

	if !ok || existing == proposed {
		return
	}

	validation.Panicf(
		"%s is configured to use %s as both a %s and a %s",
		handlerDisplayName(handlerIdent, handlerType),
		messageType,
		existing,
		proposed,
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
// messages of the given role.
func mustHaveConsumerRoute(
	types EntityMessageTypes,
	role message.Role,
	handlerIdent Identity,
	handlerType reflect.Type,
) {
	for _, r := range types.Consumed {
		if r == role {
			return
		}
	}

	validation.Panicf(
		`%s is not configured to handle any %ss, at least one Handles%s() route must be added within Configure()`,
		handlerDisplayName(handlerIdent, handlerType),
		role,
		cases.Title(language.English).String(role.String()),
	)
}

// mustHaveProducerRoute panics if the handler is not configured to produce any
// messages of the given role.
func mustHaveProducerRoute(
	types EntityMessageTypes,
	role message.Role,
	handlerIdent Identity,
	handlerType reflect.Type,
) {
	for _, r := range types.Produced {
		if r == role {
			return
		}
	}

	verb := ""
	routeFunc := ""

	switch role {
	case message.CommandRole:
		verb = "execute"
		routeFunc = "ExecutesCommand"
	case message.EventRole:
		verb = "record"
		routeFunc = "RecordsEvent"
	case message.TimeoutRole:
		verb = "schedule"
		routeFunc = "SchedulesTimeout"
	}

	validation.Panicf(
		`%s is not configured to %s any %ss, at least one %s() route must be added within Configure()`,
		handlerDisplayName(handlerIdent, handlerType),
		verb,
		role,
		routeFunc,
	)
}
