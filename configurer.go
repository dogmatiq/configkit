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

func configureIdentity(
	entityIdent *Identity,
	name, key string,
	entityType reflect.Type,
) {
	if !entityIdent.IsZero() {
		validation.Panicf(
			"%s is configured with multiple identities (%s and %s/%s), Identity() must be called exactly once within Configure()",
			entityType,
			*entityIdent,
			name,
			key,
		)
	}

	var err error
	*entityIdent, err = NewIdentity(name, key)

	if err != nil {
		validation.Panicf(
			"%s is configured with an invalid identity, %s",
			entityType,
			err,
		)
	}
}

func mustHaveValidIdentity(
	entityIdent Identity,
	entityType reflect.Type,
) {
	if entityIdent.IsZero() {
		validation.Panicf(
			"%s is configured without an identity, Identity() must be called exactly once within Configure()",
			entityType,
		)
	}
}

func configureRoute(
	types *EntityMessageTypes,
	route dogma.Route,
	handlerIdent Identity,
	handlerType reflect.Type,
) {
	switch route := route.(type) {
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
		titleCase(role.String()),
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

func titleCase(s string) string {
	return cases.Title(language.English).String(s)
}
