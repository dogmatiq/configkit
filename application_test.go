package configkit_test

import (
	"context"
	"errors"
	"reflect"

	. "github.com/dogmatiq/configkit"
	cfixtures "github.com/dogmatiq/configkit/fixtures" // can't dot-import due to conflicts
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures" // can't dot-import due to conflicts
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func FromApplication()", func() {
	var (
		aggregate   *fixtures.AggregateMessageHandler
		process     *fixtures.ProcessMessageHandler
		integration *fixtures.IntegrationMessageHandler
		projection  *fixtures.ProjectionMessageHandler
		app         *fixtures.Application
	)

	BeforeEach(func() {
		aggregate = &fixtures.AggregateMessageHandler{
			ConfigureFunc: func(c dogma.AggregateConfigurer) {
				c.Identity("<aggregate>", "<aggregate-key>")
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ProducesEventType(fixtures.MessageE{})
			},
		}

		process = &fixtures.ProcessMessageHandler{
			ConfigureFunc: func(c dogma.ProcessConfigurer) {
				c.Identity("<process>", "<process-key>")
				c.ConsumesEventType(fixtures.MessageB{})
				c.ConsumesEventType(fixtures.MessageE{}) // shared with <projection>
				c.ProducesCommandType(fixtures.MessageC{})
				c.SchedulesTimeoutType(fixtures.MessageT{})
			},
		}

		integration = &fixtures.IntegrationMessageHandler{
			ConfigureFunc: func(c dogma.IntegrationConfigurer) {
				c.Identity("<integration>", "<integration-key>")
				c.ConsumesCommandType(fixtures.MessageC{})
				c.ProducesEventType(fixtures.MessageF{})
			},
		}

		projection = &fixtures.ProjectionMessageHandler{
			ConfigureFunc: func(c dogma.ProjectionConfigurer) {
				c.Identity("<projection>", "<projection-key>")
				c.ConsumesEventType(fixtures.MessageD{})
				c.ConsumesEventType(fixtures.MessageE{}) // shared with <process>
			},
		}

		app = &fixtures.Application{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("<app>", "<app-key>")
				c.RegisterAggregate(aggregate)
				c.RegisterProcess(process)
				c.RegisterIntegration(integration)
				c.RegisterProjection(projection)
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
					MustNewIdentity("<app>", "<app-key>"),
				))
			})
		})

		Describe("func MessageNames()", func() {
			It("returns the expected message names", func() {
				Expect(cfg.MessageNames()).To(Equal(
					EntityMessageNames{
						Produced: message.NameRoles{
							cfixtures.MessageCTypeName: message.CommandRole,
							cfixtures.MessageETypeName: message.EventRole,
							cfixtures.MessageFTypeName: message.EventRole,
							cfixtures.MessageTTypeName: message.TimeoutRole,
						},
						Consumed: message.NameRoles{
							cfixtures.MessageATypeName: message.CommandRole,
							cfixtures.MessageBTypeName: message.EventRole,
							cfixtures.MessageCTypeName: message.CommandRole,
							cfixtures.MessageDTypeName: message.EventRole,
							cfixtures.MessageETypeName: message.EventRole,
							cfixtures.MessageTTypeName: message.TimeoutRole,
						},
					},
				))
			})
		})

		Describe("func MessageTypes()", func() {
			It("returns the expected message types", func() {
				Expect(cfg.MessageTypes()).To(Equal(
					EntityMessageTypes{
						Produced: message.TypeRoles{
							cfixtures.MessageCType: message.CommandRole,
							cfixtures.MessageEType: message.EventRole,
							cfixtures.MessageFType: message.EventRole,
							cfixtures.MessageTType: message.TimeoutRole,
						},
						Consumed: message.TypeRoles{
							cfixtures.MessageAType: message.CommandRole,
							cfixtures.MessageBType: message.EventRole,
							cfixtures.MessageCType: message.CommandRole,
							cfixtures.MessageDType: message.EventRole,
							cfixtures.MessageEType: message.EventRole,
							cfixtures.MessageTType: message.TimeoutRole,
						},
					},
				))
			})
		})

		Describe("func TypeName()", func() {
			It("returns the fully-qualified type name of the application", func() {
				Expect(cfg.TypeName()).To(Equal("*github.com/dogmatiq/dogma/fixtures.Application"))
			})
		})

		Describe("func ReflectType()", func() {
			It("returns the type of the application", func() {
				Expect(cfg.ReflectType()).To(Equal(reflect.TypeOf(app)))
			})
		})

		Describe("func AcceptVisitor()", func() {
			It("calls the appropriate method on the visitor", func() {
				v := &cfixtures.Visitor{
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
				v := &cfixtures.RichVisitor{
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
				c.Identity("<app>", "<aggregate-key>")
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ProducesEventType(fixtures.MessageE{})
			}

			Expect(func() {
				FromApplication(app)
			}).NotTo(Panic())
		})

		It("does not panic when the app contains multiple processes that schedule the same timeout", func() {
			process1 := &fixtures.ProcessMessageHandler{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("<process-1>", "<process-1-key>")
					c.ConsumesEventType(fixtures.MessageB{})
					c.ProducesCommandType(fixtures.MessageC{})
					c.SchedulesTimeoutType(fixtures.MessageT{})
				},
			}

			process2 := &fixtures.ProcessMessageHandler{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("<process-2>", "<process-2-key>")
					c.ConsumesEventType(fixtures.MessageB{})
					c.ProducesCommandType(fixtures.MessageC{})
					c.SchedulesTimeoutType(fixtures.MessageT{})
				},
			}

			app := &fixtures.Application{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("<app>", "<app-key>")
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
			`*fixtures.Application is configured without an identity, Identity() must be called exactly once within Configure()`,
			func() {
				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.RegisterAggregate(aggregate)
				}
			},
		),
		Entry(
			"when the app configures multiple identities",
			`*fixtures.Application is configured with multiple identities (<name>/<key> and <other>/<key>), Identity() must be called exactly once within Configure()`,
			func() {
				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.Identity("<name>", "<key>")
					c.Identity("<other>", "<key>")
					c.RegisterAggregate(aggregate)
				}
			},
		),
		Entry(
			"when the app configures an invalid name",
			`*fixtures.Application is configured with an invalid identity, invalid name "\t \n", names must be non-empty, printable UTF-8 strings with no whitespace`,
			func() {
				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.Identity("\t \n", "<key>")
					c.RegisterAggregate(aggregate)
				}
			},
		),
		Entry(
			"when the app configures an invalid key",
			`*fixtures.Application is configured with an invalid identity, invalid key "\t \n", keys must be non-empty, printable UTF-8 strings with no whitespace`,
			func() {
				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.Identity("<name>", "\t \n")
					c.RegisterAggregate(aggregate)
				}
			},
		),
		Entry(
			"when the app configures an identity that conflicts with a handler",
			`*fixtures.Application can not use the application key "<aggregate-key>", because it is already used by *fixtures.AggregateMessageHandler`,
			func() {
				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.RegisterAggregate(aggregate)
					c.Identity("<app>", "<aggregate-key>") // conflict
				}
			},
		),
		Entry(
			"when a handler is registered with a key that conflicts with the app",
			`*fixtures.AggregateMessageHandler can not use the handler key "<app-key>", because it is already used by *fixtures.Application`,
			func() {
				aggregate.ConfigureFunc = func(c dogma.AggregateConfigurer) {
					c.Identity("<aggregate>", "<app-key>") // conflict!
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageE{})
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
			`*fixtures.AggregateMessageHandler can not use the handler name "<process>", because it is already used by *fixtures.ProcessMessageHandler`,
			func() {
				aggregate.ConfigureFunc = func(c dogma.AggregateConfigurer) {
					c.Identity("<process>", "<aggregate-key>") // conflict!
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageE{})
				}

				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.Identity("<app>", "<app-key>")
					c.RegisterProcess(process)
					c.RegisterAggregate(aggregate) // register the conflicting aggregate last
				}
			},
		),
		Entry(
			"the app contains handlers with conflicting keys",
			`*fixtures.AggregateMessageHandler can not use the handler key "<process-key>", because it is already used by *fixtures.ProcessMessageHandler`,
			func() {
				aggregate.ConfigureFunc = func(c dogma.AggregateConfigurer) {
					c.Identity("<aggregate>", "<process-key>") // conflict!
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageE{})
				}

				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.Identity("<app>", "<app-key>")
					c.RegisterProcess(process)
					c.RegisterAggregate(aggregate) // register the conflicting aggregate last
				}
			},
		),
		Entry(
			"when the app contains multiple consumers of the same command",
			`*fixtures.IntegrationMessageHandler (<integration>) can not consume fixtures.MessageA commands because they are already consumed by *fixtures.AggregateMessageHandler (<aggregate>)`,
			func() {
				integration.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
					c.Identity("<integration>", "<integration-key>")
					c.ConsumesCommandType(fixtures.MessageA{}) // conflict with <aggregate>
					c.ProducesEventType(fixtures.MessageF{})
				}
			},
		),
		Entry(
			"when the app contains multiple producers of the same event",
			`*fixtures.IntegrationMessageHandler (<integration>) can not produce fixtures.MessageE events because they are already produced by *fixtures.AggregateMessageHandler (<aggregate>)`,
			func() {
				integration.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
					c.Identity("<integration>", "<integration-key>")
					c.ConsumesCommandType(fixtures.MessageC{})
					c.ProducesEventType(fixtures.MessageE{}) // conflict with <aggregate>
				}
			},
		),
		Entry(
			"when multiple handlers use a single message type in differing roles",
			`*fixtures.ProjectionMessageHandler (<projection>) configures fixtures.MessageA as an event but *fixtures.AggregateMessageHandler (<aggregate>) configures it as a command`,
			func() {
				projection.ConfigureFunc = func(c dogma.ProjectionConfigurer) {
					c.Identity("<projection>", "<projection-key>")
					c.ConsumesEventType(fixtures.MessageA{}) // conflict with <aggregate>
				}

				app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
					c.Identity("<app>", "<app-key>")
					c.RegisterAggregate(aggregate)
					c.RegisterProjection(projection)
				}
			},
		),
	)
})

var _ = Describe("func IsApplicationEqual()", func() {
	It("returns true if the two applications are equivalent", func() {
		app := &fixtures.Application{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("<app>", "<app-key>")
				c.RegisterProjection(&fixtures.ProjectionMessageHandler{
					ConfigureFunc: func(c dogma.ProjectionConfigurer) {
						c.Identity("<projection>", "<projection-key>")
						c.ConsumesEventType(fixtures.MessageE{})
					},
				})
			},
		}

		a := FromApplication(app)
		b := FromApplication(app)

		Expect(IsApplicationEqual(a, b)).To(BeTrue())
	})

	// aliasedApplication is a mock of dogma.Application that has a different Go
	// type name to fixtures.Application, used to test the type-name comparison
	// logic in IsApplicationEqual().
	type aliasedApplication struct {
		fixtures.Application
	}

	DescribeTable(
		"returns false if the applications are not equivalent",
		func(b Application) {
			app := &fixtures.Application{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("<app>", "<app-key>")
					c.RegisterProjection(&fixtures.ProjectionMessageHandler{
						ConfigureFunc: func(c dogma.ProjectionConfigurer) {
							c.Identity("<projection>", "<projection-key>")
							c.ConsumesEventType(fixtures.MessageE{})
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
				Application: fixtures.Application{
					ConfigureFunc: func(c dogma.ApplicationConfigurer) {
						c.Identity("<app>", "<app-key>")
						c.RegisterProjection(&fixtures.ProjectionMessageHandler{
							ConfigureFunc: func(c dogma.ProjectionConfigurer) {
								c.Identity("<projection>", "<projection-key>")
								c.ConsumesEventType(fixtures.MessageE{})
							},
						})
					},
				},
			}),
		),
		Entry(
			"identity name differs",
			FromApplication(&fixtures.Application{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("<app-different>", "<app-key>") // diff
					c.RegisterProjection(&fixtures.ProjectionMessageHandler{
						ConfigureFunc: func(c dogma.ProjectionConfigurer) {
							c.Identity("<projection>", "<projection-key>")
							c.ConsumesEventType(fixtures.MessageE{})
						},
					})
				},
			}),
		),
		Entry(
			"identity key differs",
			FromApplication(&fixtures.Application{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("<app>", "<app-key-different>") // diff
					c.RegisterProjection(&fixtures.ProjectionMessageHandler{
						ConfigureFunc: func(c dogma.ProjectionConfigurer) {
							c.Identity("<projection>", "<projection-key>")
							c.ConsumesEventType(fixtures.MessageE{})
						},
					})
				},
			}),
		),
		Entry(
			"messages differ",
			FromApplication(&fixtures.Application{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("<app>", "<app-key>")
					c.RegisterProjection(&fixtures.ProjectionMessageHandler{
						ConfigureFunc: func(c dogma.ProjectionConfigurer) {
							c.Identity("<projection>", "<projection-key>")
							c.ConsumesEventType(fixtures.MessageE{})
							c.ConsumesEventType(fixtures.MessageX{}) // diff
						},
					})
				},
			}),
		),
		Entry(
			"handlers differ",
			FromApplication(&fixtures.Application{
				ConfigureFunc: func(c dogma.ApplicationConfigurer) {
					c.Identity("<app>", "<app-key>")
					c.RegisterProjection(&fixtures.ProjectionMessageHandler{
						ConfigureFunc: func(c dogma.ProjectionConfigurer) {
							c.Identity("<projection-different>", "<projection-key>") // diff
							c.ConsumesEventType(fixtures.MessageE{})
						},
					})
				},
			}),
		),
	)
})
