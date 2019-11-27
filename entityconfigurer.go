package configkit

// handlerConfigurer is a partial implementation of the configurer interfaces for
// all of the Dogma entity types.
//
// - dogma.ApplicationConfigurer
// - dogma.AggregateConfigurer
// - dogma.ProcessConfigurer
// - dogma.IntegrationConfigurer
// - dogma.ProjectionConfigurer
type entityConfigurer struct {
	// target is the entity to populate with the configuration values.
	target *entity
}

// Identity sets the handler's identity.
func (c *entityConfigurer) Identity(name string, key string) {
	if !c.target.ident.IsZero() {
		Panicf(
			"%s is configured with multiple identities (%s and %s/%s), Identity() must be called exactly once within Configure()",
			c.target.rt,
			c.target.ident,
			name,
			key,
		)
	}

	var err error
	c.target.ident, err = NewIdentity(name, key)

	if err != nil {
		Panicf(
			"%s is configured with an invalid identity, %s",
			c.target.rt,
			err,
		)
	}
}

// validate panics if the configuration is invalid.
func (c *entityConfigurer) validate() {
	if c.target.ident.IsZero() {
		Panicf(
			"%s is configured without an identity, Identity() must be called exactly once within Configure()",
			c.target.rt,
		)
	}
}
