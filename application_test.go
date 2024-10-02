package configkit_test

import (
	"context"
	"errors"
	"reflect"

	. "github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	"github.com/jmalloc/gomegax"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func FromApplication()", func() {

	var (
		aggregate   *AggregateMessageHandlerStub
		process     *ProcessMessageHandlerStub
		integration *IntegrationMessageHandlerStub
		projection  *ProjectionMessageHandlerStub
		app         *ApplicationStub
	)

	BeforeEach(func() {
		aggregate = &AggregateMessageHandlerStub{
			ConfigureFunc: func(c dogma.AggregateConfigurer) {
				c.Identity("<aggregate>", aggregateKey)
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			},
		}

		process = &ProcessMessageHandlerStub{
			ConfigureFunc: func(c dogma.ProcessConfigurer) {
				c.Identity("<process>", processKey)
				c.Routes(
					dogma.HandlesEvent[EventStub[TypeA]](), // shared with <projection>
					dogma.HandlesEvent[EventStub[TypeB]](),
					dogma.ExecutesCommand[CommandStub[TypeB]](),
					dogma.SchedulesTimeout[TimeoutStub[TypeA]](),
				)
			},
		}

		integration = &IntegrationMessageHandlerStub{
			ConfigureFunc: func(c dogma.IntegrationConfigurer) {
				c.Identity("<integration>", integrationKey)
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeB]](),
					dogma.RecordsEvent[EventStub[TypeC]](),
				)
			},
		}

		projection = &ProjectionMessageHandlerStub{
			ConfigureFunc: func(c dogma.ProjectionConfigurer) {
				c.Identity("<projection>", projectionKey)
				c.Routes(
					dogma.HandlesEvent[EventStub[TypeA]](), // shared with <process>
					dogma.HandlesEvent[EventStub[TypeD]](),
				)
			},
		}

		disabled := &ProjectionMessageHandlerStub{
			ConfigureFunc: func(c dogma.ProjectionConfigurer) {
				// Verify that disabled handlers with no identity / route
				// configuration are excluded from the application
				// configuration.
				c.Disable()
			},
		}

		app = &ApplicationStub{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("<app>", appKey)
				c.RegisterAggregate(aggregate)
				c.RegisterProcess(process)
				c.RegisterIntegration(integration)
				c.RegisterProjection(projection)
				c.RegisterProjection(disabled)
			},
		}
	})

	When("the configuration is valid", func() {
		var cfg RichApplication

		BeforeEach(func() {
			cfg = FromApplication(app)
		})

		Describe("func Identity()", func() {
			It("returns the application identity", func() {
				Expect(cfg.Identity()).To(Equal(
					MustNewIdentity("<app>", appKey),
				))
			})
		})

		Describe("func MessageNames()", func() {
			It("returns the expected message names", func() {
				Expect(cfg.MessageNames()).To(gomegax.EqualX(
					EntityMessages[message.Name]{
						message.NameOf(CommandA1): {
							Kind:       message.CommandKind,
							IsConsumed: true,
						},
						message.NameOf(CommandB1): {
							Kind:       message.CommandKind,
							IsProduced: true,
							IsConsumed: true,
						},
						message.NameOf(EventA1): {
							Kind:       message.EventKind,
							IsProduced: true,
							IsConsumed: true,
						},
						message.NameOf(EventB1): {
							Kind:       message.EventKind,
							IsConsumed: true,
						},
						message.NameOf(EventC1): {
							Kind:       message.EventKind,
							IsProduced: true,
						},
						message.NameOf(EventD1): {
							Kind:       message.EventKind,
							IsConsumed: true,
						},
						message.NameOf(TimeoutA1): {
							Kind:       message.TimeoutKind,
							IsProduced: true,
							IsConsumed: true,
						},
					},
				))
			})
		})

		Describe("func MessageTypes()", func() {
			It("returns the expected message types", func() {
				Expect(cfg.MessageTypes()).To(Equal(
					EntityMessages[message.Type]{
						message.TypeOf(CommandA1): {
							Kind:       message.CommandKind,
							IsConsumed: true,
						},
						message.TypeOf(CommandB1): {
							Kind:       message.CommandKind,
							IsProduced: true,
							IsConsumed: true,
						},
						message.TypeOf(EventA1): {
							Kind:       message.EventKind,
							IsProduced: true,
							IsConsumed: true,
						},
						message.TypeOf(EventB1): {
							Kind:       message.EventKind,
							IsConsumed: true,
						},
						message.TypeOf(EventC1): {
							Kind:       message.EventKind,
							IsProduced: true,
						},
						message.TypeOf(EventD1): {
							Kind:       message.EventKind,
							IsConsumed: true,
						},
						message.TypeOf(TimeoutA1): {
							Kind:       message.TimeoutKind,
							IsProduced: true,
							IsConsumed: true,
						},
					},
				))
			})
		})

		Describe("func TypeName()", func() {
			It("returns the fully-qualified type name of the application", func() {
				Expect(cfg.TypeName()).To(Equal("*github.com/dogmatiq/enginekit/enginetest/stubs.ApplicationStub"))
			})
		})

		Describe("func ReflectType()", func() {
			It("returns the type of the application", func() {
				Expect(cfg.ReflectType()).To(Equal(reflect.TypeOf(app)))
			})
		})

		Describe("func AcceptVisitor()", func() {
			It("calls the appropriate method on the visitor", func() {
				v := &visitorStub{
					VisitApplicationFunc: func(_ context.Context, c Application) error {
						Expect(c).To(BeIdenticalTo(cfg))
						return errors.New("<error>")
					},
				}

				err := cfg.AcceptVisitor(context.Background(), v)
				Expect(err).To(MatchError("<error>"))
			})
		})

		Describe("func AcceptRichVisitor()", func() {
			It("calls the appropriate method on the visitor", func() {
				v := &richVisitorStub{
					VisitRichApplicationFunc: func(_ context.Context, c RichApplication) error {
						Expect(c).To(BeIdenticalTo(cfg))
						return errors.New("<error>")
					},
				}

				err := cfg.AcceptRichVisitor(context.Background(), v)
				Expect(err).To(MatchError("<error>"))
			})
		})

		Describe("func Handlers()", func() {
			It("returns a set containing all handlers in the application", func() {
				Expect(cfg.Handlers()).To(Equal(
					NewHandlerSet(
						FromAggregate(aggregate),
						FromProcess(process),
						FromIntegration(integration),
						FromProjection(projection),
					),
				))
			})
		})

		Describe("func RichHandlers()", func() {
			It("returns a set containing all handlers in the application", func() {
				Expect(cfg.RichHandlers()).To(Equal(
					NewRichHandlerSet(
						FromAggregate(aggregate),
						FromProcess(process),
						FromIntegration(integration),
						FromProjection(projection),
					),
				))
			})
		})

		Describe("func Application()", func() {
			It("returns the underlying application", func() {
				Expect(cfg.Application()).To(BeIdenticalTo(app))
			})
		})

		It("does not panic when the app name is shared with handler", func() {
			aggregate.ConfigureFunc = func(c dogma.AggregateConfigurer) {
				c.Identity("<app>", aggregateKey)
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			}

			Expect(func() {
				FromApplication(app)
			}).NotTo(Panic())
		})

		It("does not panic when the app contains multiple processes that schedule the same timeout", func() {
			process1 := &ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("<process-1>", "51621cac-73e2-48fa-95ad-8c3d06ab2ac3")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeB]](),
						dogma.ExecutesCommand[CommandStub[TypeB]](),
						dogma.SchedulesTimeout[TimeoutStub[TypeA]](),
					)
				},
			}

			process2 := &ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("<process-2>", "97abc0e1-39c8-434a-8ff2-1f0e2d37486e")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeB]](),
						dogma.ExecutesCommand[CommandStub[TypeB]](),
						dogma.SchedulesTimeout[TimeoutStub[TypeA]](),
					)
				},
			}

			app := &ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("<app>", appKey)
					c.RegisterProcess(process1)
					c.RegisterProcess(process2)
				},
			}

			Expect(func() {
				FromApplication(app)
			}).NotTo(Panic())
		})
	})

	DescribeTable(
		"when the configuration is invalid",
		func(
			msg string,
			fn func(),
		) {
			fn()

			var err error
			func() {
				defer Recover(&err)
				FromApplication(app)
			}()

			Expect(err).Should(HaveOccurred())
			if msg != "" {
				Expect(err).To(MatchError(msg))
			}
		},
		Entry(
			"when the app does not configure anything",
			"", // any error
			func() {
				app.ConfigureFunc = nil
			},
		),
		Entry(
			"when the app does not configure an identity",
			`*stubs.ApplicationStub is configured without an identity, Identity() must be called exactly once within Configure()`,
			func() {
				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.RegisterAggregate(aggregate)
				}
			},
		),
		Entry(
			"when the app configures multiple identities",
			`*stubs.ApplicationStub is configured with multiple identities (<name>/59a82a24-a181-41e8-9b93-17a6ce86956e and <other>/59a82a24-a181-41e8-9b93-17a6ce86956e), Identity() must be called exactly once within Configure()`,
			func() {
				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.Identity("<name>", appKey)
					c.Identity("<other>", appKey)
					c.RegisterAggregate(aggregate)
				}
			},
		),
		Entry(
			"when the app configures an invalid name",
			`*stubs.ApplicationStub is configured with an invalid identity, invalid name "\t \n", names must be non-empty, printable UTF-8 strings with no whitespace`,
			func() {
				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.Identity("\t \n", appKey)
					c.RegisterAggregate(aggregate)
				}
			},
		),
		Entry(
			"when the app configures an invalid key",
			`*stubs.ApplicationStub is configured with an invalid identity, invalid key "\t \n", keys must be RFC 4122 UUIDs`,
			func() {
				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.Identity("<name>", "\t \n")
					c.RegisterAggregate(aggregate)
				}
			},
		),
		Entry(
			"when the app configures an identity that conflicts with a handler",
			`*stubs.ApplicationStub can not use the application key "14769f7f-87fe-48dd-916e-5bcab6ba6aca", because it is already used by *stubs.AggregateMessageHandlerStub`,
			func() {
				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.RegisterAggregate(aggregate)
					c.Identity("<app>", aggregateKey) // conflict
				}
			},
		),
		Entry(
			"when a handler is registered with a key that conflicts with the app",
			`*stubs.AggregateMessageHandlerStub can not use the handler key "59a82a24-a181-41e8-9b93-17a6ce86956e", because it is already used by *stubs.ApplicationStub`,
			func() {
				aggregate.ConfigureFunc = func(c dogma.AggregateConfigurer) {
					c.Identity("<aggregate>", appKey) // conflict!
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
						dogma.RecordsEvent[EventStub[TypeA]](),
					)
				}
			},
		),
		Entry(
			"when the app contains an invalid aggregate configuration",
			"", // any error
			func() {
				aggregate.ConfigureFunc = nil
			},
		),
		Entry(
			"when the app contains an invalid process configuration",
			"", // any error
			func() {
				process.ConfigureFunc = nil
			},
		),
		Entry(
			"when the app contains an invalid integration configuration",
			"", // any error
			func() {
				integration.ConfigureFunc = nil
			},
		),
		Entry(
			"when the app contains an invalid projection configuration",
			"", // any error
			func() {
				projection.ConfigureFunc = nil
			},
		),
		Entry(
			"the app contains handlers with conflicting names",
			`*stubs.AggregateMessageHandlerStub can not use the handler name "<process>", because it is already used by *stubs.ProcessMessageHandlerStub`,
			func() {
				aggregate.ConfigureFunc = func(c dogma.AggregateConfigurer) {
					c.Identity("<process>", aggregateKey) // conflict!
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
						dogma.RecordsEvent[EventStub[TypeA]](),
					)
				}

				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.Identity("<app>", appKey)
					c.RegisterProcess(process)
					c.RegisterAggregate(aggregate) // register the conflicting aggregate last
				}
			},
		),
		Entry(
			"the app contains handlers with conflicting keys",
			`*stubs.AggregateMessageHandlerStub can not use the handler key "bea52cf4-e403-4b18-819d-88ade7836308", because it is already used by *stubs.ProcessMessageHandlerStub`,
			func() {
				aggregate.ConfigureFunc = func(c dogma.AggregateConfigurer) {
					c.Identity("<aggregate>", processKey) // conflict!
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
						dogma.RecordsEvent[EventStub[TypeA]](),
					)
				}

				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.Identity("<app>", appKey)
					c.RegisterProcess(process)
					c.RegisterAggregate(aggregate) // register the conflicting aggregate last
				}
			},
		),
		Entry(
			"when the app contains multiple handlers of the same command",
			`*stubs.IntegrationMessageHandlerStub (<integration>) can not handle stubs.CommandStub[TypeA] commands because they are already configured to be handled by *stubs.AggregateMessageHandlerStub (<aggregate>)`,
			func() {
				integration.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
					c.Identity("<integration>", integrationKey)
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](), // conflict with <aggregate>
						dogma.RecordsEvent[EventStub[TypeC]](),
					)
				}
			},
		),
		Entry(
			"when the app contains multiple handlers that record the same event",
			`*stubs.IntegrationMessageHandlerStub (<integration>) can not record stubs.EventStub[TypeA] events because they are already configured to be recorded by *stubs.AggregateMessageHandlerStub (<aggregate>)`,
			func() {
				integration.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
					c.Identity("<integration>", integrationKey)
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeB]](),
						dogma.RecordsEvent[EventStub[TypeA]](), // conflict with <aggregate>
					)
				}
			},
		),
	)
})

var _ = Describe("func IsApplicationEqual()", func() {
	It("returns true if the two applications are equivalent", func() {
		app := &ApplicationStub{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("<app>", appKey)
				c.RegisterProjection(&ProjectionMessageHandlerStub{
					ConfigureFunc: func(c dogma.ProjectionConfigurer) {
						c.Identity("<projection>", projectionKey)
						c.Routes(
							dogma.HandlesEvent[EventStub[TypeA]](),
						)
					},
				})
			},
		}

		a := FromApplication(app)
		b := FromApplication(app)

		Expect(IsApplicationEqual(a, b)).To(BeTrue())
	})

	// aliasedApplication is a mock of [dogma.Application] that has a different
	// Go type name to [ApplicationStub], used to test the type-name comparison
	// logic in IsApplicationEqual().
	type aliasedApplication struct {
		ApplicationStub
	}

	DescribeTable(
		"returns false if the applications are not equivalent",
		func(b Application) {
			app := &ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("<app>", appKey)
					c.RegisterProjection(&ProjectionMessageHandlerStub{
						ConfigureFunc: func(c dogma.ProjectionConfigurer) {
							c.Identity("<projection>", projectionKey)
							c.Routes(
								dogma.HandlesEvent[EventStub[TypeA]](),
							)
						},
					})
				},
			}

			a := FromApplication(app)

			Expect(IsApplicationEqual(a, b)).To(BeFalse())
		},
		Entry(
			"type differs",
			FromApplication(&aliasedApplication{
				ApplicationStub: ApplicationStub{
					ConfigureFunc: func(c dogma.ApplicationConfigurer) {
						c.Identity("<app>", appKey)
						c.RegisterProjection(&ProjectionMessageHandlerStub{
							ConfigureFunc: func(c dogma.ProjectionConfigurer) {
								c.Identity("<projection>", projectionKey)
								c.Routes(
									dogma.HandlesEvent[EventStub[TypeA]](),
								)
							},
						})
					},
				},
			}),
		),
		Entry(
			"identity name differs",
			FromApplication(&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("<app-different>", appKey) // diff
					c.RegisterProjection(&ProjectionMessageHandlerStub{
						ConfigureFunc: func(c dogma.ProjectionConfigurer) {
							c.Identity("<projection>", projectionKey)
							c.Routes(
								dogma.HandlesEvent[EventStub[TypeA]](),
							)
						},
					})
				},
			}),
		),
		Entry(
			"identity key differs",
			FromApplication(&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("<app>", "b7deb466-0fb7-4e89-b4dd-a32cdb1e1823") // diff
					c.RegisterProjection(&ProjectionMessageHandlerStub{
						ConfigureFunc: func(c dogma.ProjectionConfigurer) {
							c.Identity("<projection>", projectionKey)
							c.Routes(
								dogma.HandlesEvent[EventStub[TypeA]](),
							)
						},
					})
				},
			}),
		),
		Entry(
			"messages differ",
			FromApplication(&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("<app>", appKey)
					c.RegisterProjection(&ProjectionMessageHandlerStub{
						ConfigureFunc: func(c dogma.ProjectionConfigurer) {
							c.Identity("<projection>", projectionKey)
							c.Routes(
								dogma.HandlesEvent[EventStub[TypeA]](),
								dogma.HandlesEvent[EventStub[TypeX]](), // diff
							)
						},
					})
				},
			}),
		),
		Entry(
			"handlers differ",
			FromApplication(&ApplicationStub{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("<app>", appKey)
					c.RegisterProjection(&ProjectionMessageHandlerStub{
						ConfigureFunc: func(c dogma.ProjectionConfigurer) {
							c.Identity("<projection-different>", projectionKey) // diff
							c.Routes(
								dogma.HandlesEvent[EventStub[TypeA]](),
							)
						},
					})
				},
			}),
		),
	)
})
