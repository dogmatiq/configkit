package configkit

import "github.com/dogmatiq/configkit/internal/validation"

// handlerConfigurer is a partial implementation of the configurer interfaces for
// all of the Dogma entity types.
//
//   - [dogma.ApplicationConfigurer]
//   - [dogma.AggregateConfigurer]
//   - [dogma.ProcessConfigurer]
//   - [dogma.IntegrationConfigurer]
//   - [dogma.ProjectionConfigurer]
type entityConfigurer struct {
	// entity is the target entity to populate with the configuration values.
	entity     *entity
	configured bool
}

// Identity sets the entity's identity.
func (c *entityConfigurer) Identity(n string, k string) {
	c.configured = true

	if !c.entity.ident.IsZero() {
		validation.Panicf(
			"%s is configured with multiple identities (%s and %s/%s), Identity() must be called exactly once within Configure()",
			c.entity.rt,
			c.entity.ident,
			n,
			k,
		)
	}

	var err error
	c.entity.ident, err = NewIdentity(n, k)

	if err != nil {
		validation.Panicf(
			"%s is configured with an invalid identity, %s",
			c.entity.rt,
			err,
		)
	}
}

func (c *entityConfigurer) isConfigured() bool {
	return c.configured
}

// mustValidate panics if the configuration is invalid.
func (c *entityConfigurer) mustValidate() {
	if c.entity.ident.IsZero() {
		validation.Panicf(
			"%s is configured without an identity, Identity() must be called exactly once within Configure()",
			c.entity.rt,
		)
	}
}

// displayName returns a human-readable string used to refer to the entity in
// error messages.
func (c *entityConfigurer) displayName() string {
	s := c.entity.rt.String()

	if !c.entity.ident.IsZero() {
		s += " (" + c.entity.ident.Name + ")"
	}

	return s
}
