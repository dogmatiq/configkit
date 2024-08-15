package configkit

import (
	"reflect"

	"github.com/dogmatiq/configkit/internal/validation"
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

func mustValidateIdentity(
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
