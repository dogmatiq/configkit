package configkit

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
func (c *entityConfigurer) Identity(n, k string) {
	c.configured = true
	configureIdentity(
		c.entity.rt,
		&c.entity.ident,
		n,
		k,
	)
}

func (c *entityConfigurer) isConfigured() bool {
	return c.configured
}

// mustValidate panics if the configuration is invalid.
func (c *entityConfigurer) mustValidate() {
	mustValidateIdentity(c.entity.rt, c.entity.ident)
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
