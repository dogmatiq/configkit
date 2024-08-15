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
	rt reflect.Type,
	id *Identity,
	n, k string,
) {
	if !id.IsZero() {
		validation.Panicf(
			"%s is configured with multiple identities (%s and %s/%s), Identity() must be called exactly once within Configure()",
			rt,
			*id,
			n,
			k,
		)
	}

	var err error
	*id, err = NewIdentity(n, k)

	if err != nil {
		validation.Panicf(
			"%s is configured with an invalid identity, %s",
			rt,
			err,
		)
	}
}

func mustHaveValidIdentity(
	rt reflect.Type,
	id Identity,
) {
	if id.IsZero() {
		validation.Panicf(
			"%s is configured without an identity, Identity() must be called exactly once within Configure()",
			rt,
		)
	}
}

func configureRoute(
	handlerType reflect.Type,
	handlerIdent Identity,
	types *EntityMessageTypes,
	route dogma.Route,
) {
	switch route := route.(type) {
	case dogma.HandlesCommandRoute:
		configureConsumerRoute(handlerType, handlerIdent, types, "HandlesCommand", route.Type, message.CommandRole)
	case dogma.RecordsEventRoute:
		configureProducerRoute(handlerType, handlerIdent, types, "RecordsEvent", route.Type, message.EventRole)
	case dogma.HandlesEventRoute:
		configureConsumerRoute(handlerType, handlerIdent, types, "HandlesEvent", route.Type, message.EventRole)
	case dogma.ExecutesCommandRoute:
		configureProducerRoute(handlerType, handlerIdent, types, "ExecutesCommand", route.Type, message.CommandRole)
	case dogma.SchedulesTimeoutRoute:
		configureConsumerRoute(handlerType, handlerIdent, types, "SchedulesTimeout", route.Type, message.TimeoutRole)
		configureProducerRoute(handlerType, handlerIdent, types, "SchedulesTimeout", route.Type, message.TimeoutRole)
	default:
		panic(fmt.Sprintf("unsupported route type: %T", route))
	}
}

func configureConsumerRoute(
	handlerType reflect.Type,
	handlerIdent Identity,
	types *EntityMessageTypes,
	routeFunc string,
	t reflect.Type,
	r message.Role,
) {
	mt := message.TypeFromReflect(t)

	guardAgainstConflictingRoles(
		handlerType,
		handlerIdent,
		types,
		mt,
		r,
	)

	if types.Consumed.Has(mt) {
		validation.Panicf(
			"%s is configured with multiple %s() routes for %s, should these refer to different message types?",
			handlerDisplayName(handlerType, handlerIdent),
			routeFunc,
			mt,
		)
	}

	if types.Consumed == nil {
		types.Consumed = message.TypeRoles{}
	}
	types.Consumed.Add(mt, r)
}

func configureProducerRoute(
	handlerType reflect.Type,
	handlerIdent Identity,
	types *EntityMessageTypes,
	routeFunc string,
	t reflect.Type,
	r message.Role,
) {
	mt := message.TypeFromReflect(t)

	guardAgainstConflictingRoles(
		handlerType,
		handlerIdent,
		types,
		mt,
		r,
	)

	if types.Produced.Has(mt) {
		validation.Panicf(
			"%s is configured with multiple %s() routes for %s, should these refer to different message types?",
			handlerDisplayName(handlerType, handlerIdent),
			routeFunc,
			mt,
		)
	}

	if types.Produced == nil {
		types.Produced = message.TypeRoles{}
	}
	types.Produced.Add(mt, r)
}

func guardAgainstConflictingRoles(
	handlerType reflect.Type,
	handlerIdent Identity,
	types *EntityMessageTypes,
	mt message.Type,
	r message.Role,
) {
	x, ok := types.RoleOf(mt)

	if !ok || x == r {
		return
	}

	validation.Panicf(
		"%s is configured to use %s as both a %s and a %s",
		handlerDisplayName(handlerType, handlerIdent),
		mt,
		x,
		r,
	)
}

func handlerDisplayName(
	rt reflect.Type,
	id Identity,
) string {
	s := rt.String()

	if !id.IsZero() {
		s += " (" + id.Name + ")"
	}

	return s
}

// mustHaveConsumerRoute panics if the handler is not configured to handle any
// messages of the given role.
func mustHaveConsumerRoute(
	handlerType reflect.Type,
	handlerIdent Identity,
	types EntityMessageTypes,
	r message.Role,
) {
	for _, x := range types.Consumed {
		if x == r {
			return
		}
	}

	validation.Panicf(
		`%s is not configured to handle any %ss, at least one Handles%s() route must be added within Configure()`,
		handlerDisplayName(handlerType, handlerIdent),
		r,
		titleCase(r.String()),
	)
}

// mustHaveProducerRoute panics if the handler is not configured to produce any
// messages of the given role.
func mustHaveProducerRoute(
	handlerType reflect.Type,
	handlerIdent Identity,
	types EntityMessageTypes,
	r message.Role,
) {
	for _, x := range types.Produced {
		if x == r {
			return
		}
	}

	verb := ""
	routeFunc := ""

	switch r {
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
		handlerDisplayName(handlerType, handlerIdent),
		verb,
		r,
		routeFunc,
	)
}

func titleCase(s string) string {
	return cases.Title(language.English).String(s)
}
