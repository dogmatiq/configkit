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

		XDescribe("func MessageNames()", func() {
			It("returns the expected message names", func() {
				Expect(cfg.MessageNames()).To(Equal(
					EntityMessageNames{
						Roles: message.NameRoles{
							cfixtures.MessageATypeName: message.CommandRole,
							cfixtures.MessageBTypeName: message.EventRole,
							cfixtures.MessageCTypeName: message.CommandRole,
							cfixtures.MessageDTypeName: message.EventRole,
							cfixtures.MessageETypeName: message.EventRole,
							cfixtures.MessageFTypeName: message.EventRole,
							cfixtures.MessageTTypeName: message.TimeoutRole,
						},
						Produced: message.NameRoles{
							cfixtures.MessageATypeName: message.CommandRole,
							cfixtures.MessageBTypeName: message.EventRole,
							cfixtures.MessageCTypeName: message.CommandRole,
							cfixtures.MessageDTypeName: message.EventRole,
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
							cfixtures.MessageFTypeName: message.EventRole,
							cfixtures.MessageTTypeName: message.TimeoutRole,
						},
					},
				))
			})
		})

		XDescribe("func MessageTypes()", func() {
			It("returns the expected message types", func() {
				Expect(cfg.MessageTypes()).To(Equal(
					EntityMessageTypes{
						Roles: message.TypeRoles{
							cfixtures.MessageAType: message.CommandRole,
							cfixtures.MessageBType: message.EventRole,
							cfixtures.MessageCType: message.CommandRole,
							cfixtures.MessageDType: message.EventRole,
							cfixtures.MessageEType: message.EventRole,
							cfixtures.MessageFType: message.EventRole,
							cfixtures.MessageTType: message.TimeoutRole,
						},
						Produced: message.TypeRoles{
							cfixtures.MessageAType: message.CommandRole,
							cfixtures.MessageBType: message.EventRole,
							cfixtures.MessageCType: message.CommandRole,
							cfixtures.MessageDType: message.EventRole,
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
							cfixtures.MessageFType: message.EventRole,
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

		XDescribe("func Handlers()", func() {
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

		XDescribe("func RichHandlers()", func() {
			It("returns a set containing all handlers in the application", func() {
				Expect(cfg.Handlers()).To(Equal(
					NewRichHandlerSet(
						FromAggregate(aggregate),
						FromProcess(process),
						FromIntegration(integration),
						FromProjection(projection),
					),
				))
			})
		})

		XDescribe("func ForeignMessageNames()", func() {
			It("returns the set of messages that belong to another application", func() {
				Expect(cfg.ForeignMessageNames()).To(Equal(
					message.NameRoles{
						cfixtures.MessageXTypeName: message.CommandRole,
						cfixtures.MessageYTypeName: message.EventRole,
					},
				))
			})
		})

		XDescribe("func ForeignMessageTypes()", func() {
			It("returns the set of messages that belong to another application", func() {
				Expect(cfg.ForeignMessageTypes()).To(Equal(
					message.TypeRoles{
						cfixtures.MessageXType: message.CommandRole,
						cfixtures.MessageYType: message.EventRole,
					},
				))
			})
		})

		Describe("func Application()", func() {
			It("returns the underlying application", func() {
				Expect(cfg.Application()).To(BeIdenticalTo(app))
			})
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
		// 	When("the app does not configure an identity", func() {
		// 		BeforeEach(func() {
		// 			app.ConfigureFunc = nil
		// 		})

		// 		It("returns a descriptive error", func() {
		// 			_, err := FromApplication(app)

		// 			Expect(err).To(Equal(
		// 				Error(
		// 					"*fixtures.Application.Configure() did not call ApplicationConfigurer.Identity()",
		// 				),
		// 			))
		// 		})
		// 	})

		// 	When("the app configures multiple identities", func() {
		// 		BeforeEach(func() {
		// 			app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
		// 				c.Identity("<name>", "<key>")
		// 				c.Identity("<other>", "<key>")
		// 			}
		// 		})

		// 		It("returns a descriptive error", func() {
		// 			_, err := FromApplication(app)

		// 			Expect(err).To(Equal(
		// 				Error(
		// 					`*fixtures.Application.Configure() has already called ApplicationConfigurer.Identity("<name>", "<key>")`,
		// 				),
		// 			))
		// 		})
		// 	})

		// 	When("the app configures an invalid application name", func() {
		// 		BeforeEach(func() {
		// 			app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
		// 				c.Identity("\t \n", "<app-key>")
		// 			}
		// 		})

		// 		It("returns a descriptive error", func() {
		// 			_, err := FromApplication(app)

		// 			Expect(err).To(Equal(
		// 				Error(
		// 					`*fixtures.Application.Configure() called ApplicationConfigurer.Identity() with an invalid name "\t \n", names must be non-empty, printable UTF-8 strings with no whitespace`,
		// 				),
		// 			))
		// 		})
		// 	})

		// 	When("the app configures an invalid application key", func() {
		// 		BeforeEach(func() {
		// 			app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
		// 				c.Identity("<app>", "\t \n")
		// 			}
		// 		})

		// 		It("returns a descriptive error", func() {
		// 			_, err := FromApplication(app)

		// 			Expect(err).To(Equal(
		// 				Error(
		// 					`*fixtures.Application.Configure() called ApplicationConfigurer.Identity() with an invalid key "\t \n", keys must be non-empty, printable UTF-8 strings with no whitespace`,
		// 				),
		// 			))
		// 		})
		// 	})

		// 	When("the app contains an invalid handler configurations", func() {
		// 		It("returns an error when an aggregate is misconfigured", func() {
		// 			aggregate.ConfigureFunc = nil

		// 			_, err := FromApplication(app)

		// 			Expect(err).Should(HaveOccurred())
		// 		})

		// 		It("returns an error when a process is misconfigured", func() {
		// 			process.ConfigureFunc = nil

		// 			_, err := FromApplication(app)

		// 			Expect(err).Should(HaveOccurred())
		// 		})

		// 		It("returns an error when an integration is misconfigured", func() {
		// 			integration.ConfigureFunc = nil

		// 			_, err := FromApplication(app)

		// 			Expect(err).Should(HaveOccurred())
		// 		})

		// 		It("returns an error when a projection is misconfigured", func() {
		// 			projection.ConfigureFunc = nil

		// 			_, err := FromApplication(app)

		// 			Expect(err).Should(HaveOccurred())
		// 		})
		// 	})

		// 	When("the app contains conflicting handler identities", func() {
		// 		It("returns an error when an aggregate name is in conflict", func() {
		// 			aggregate.ConfigureFunc = func(c dogma.AggregateConfigurer) {
		// 				c.Identity("<process>", "<aggregate-key>") // conflict!
		// 				c.ConsumesCommandType(fixtures.MessageA{})
		// 				c.ProducesEventType(fixtures.MessageE{})
		// 			}

		// 			app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
		// 				c.Identity("<app>", "<app-key>")
		// 				c.RegisterProcess(process)
		// 				c.RegisterAggregate(aggregate) // register the conflicting aggregate last
		// 			}

		// 			_, err := FromApplication(app)

		// 			Expect(err).To(Equal(
		// 				Error(
		// 					`*fixtures.AggregateMessageHandler can not use the handler name "<process>", because it is already used by *fixtures.ProcessMessageHandler`,
		// 				),
		// 			))
		// 		})

		// 		It("returns an error when an aggregate key is in conflict", func() {
		// 			aggregate.ConfigureFunc = func(c dogma.AggregateConfigurer) {
		// 				c.Identity("<aggregate>", "<process-key>") // conflict!
		// 				c.ConsumesCommandType(fixtures.MessageA{})
		// 				c.ProducesEventType(fixtures.MessageE{})
		// 			}

		// 			app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
		// 				c.Identity("<app>", "<app-key>")
		// 				c.RegisterProcess(process)
		// 				c.RegisterAggregate(aggregate) // register the conflicting aggregate last
		// 			}

		// 			_, err := FromApplication(app)

		// 			Expect(err).To(Equal(
		// 				Error(
		// 					`*fixtures.AggregateMessageHandler can not use the handler key "<process-key>", because it is already used by *fixtures.ProcessMessageHandler`,
		// 				),
		// 			))
		// 		})

		// 		It("returns an error when a process name is in conflict", func() {
		// 			process.ConfigureFunc = func(c dogma.ProcessConfigurer) {
		// 				c.Identity("<aggregate>", "<process-key>") // conflict!
		// 				c.ConsumesEventType(fixtures.MessageB{})
		// 				c.ProducesCommandType(fixtures.MessageC{})
		// 			}

		// 			_, err := FromApplication(app)

		// 			Expect(err).To(Equal(
		// 				Error(
		// 					`*fixtures.ProcessMessageHandler can not use the handler name "<aggregate>", because it is already used by *fixtures.AggregateMessageHandler`,
		// 				),
		// 			))
		// 		})

		// 		It("returns an error when a process key is in conflict", func() {
		// 			process.ConfigureFunc = func(c dogma.ProcessConfigurer) {
		// 				c.Identity("<process>", "<aggregate-key>") // conflict!
		// 				c.ConsumesEventType(fixtures.MessageB{})
		// 				c.ProducesCommandType(fixtures.MessageC{})
		// 			}

		// 			_, err := FromApplication(app)

		// 			Expect(err).To(Equal(
		// 				Error(
		// 					`*fixtures.ProcessMessageHandler can not use the handler key "<aggregate-key>", because it is already used by *fixtures.AggregateMessageHandler`,
		// 				),
		// 			))
		// 		})

		// 		It("returns an error when an integration name is in conflict", func() {
		// 			integration.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
		// 				c.Identity("<process>", "<integration-key>") // conflict!
		// 				c.ConsumesCommandType(fixtures.MessageC{})
		// 				c.ProducesEventType(fixtures.MessageF{})
		// 			}

		// 			_, err := FromApplication(app)

		// 			Expect(err).To(Equal(
		// 				Error(
		// 					`*fixtures.IntegrationMessageHandler can not use the handler name "<process>", because it is already used by *fixtures.ProcessMessageHandler`,
		// 				),
		// 			))
		// 		})

		// 		It("returns an error when an integration key is in conflict", func() {
		// 			integration.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
		// 				c.Identity("<integration>", "<process-key>") // conflict!
		// 				c.ConsumesCommandType(fixtures.MessageC{})
		// 				c.ProducesEventType(fixtures.MessageF{})
		// 			}

		// 			_, err := FromApplication(app)

		// 			Expect(err).To(Equal(
		// 				Error(
		// 					`*fixtures.IntegrationMessageHandler can not use the handler key "<process-key>", because it is already used by *fixtures.ProcessMessageHandler`,
		// 				),
		// 			))
		// 		})

		// 		It("returns an error when a projection name is in conflict", func() {
		// 			projection.ConfigureFunc = func(c dogma.ProjectionConfigurer) {
		// 				c.Identity("<integration>", "<projection-key>") // conflict!
		// 				c.ConsumesEventType(fixtures.MessageD{})
		// 			}

		// 			_, err := FromApplication(app)

		// 			Expect(err).To(Equal(
		// 				Error(
		// 					`*fixtures.ProjectionMessageHandler can not use the handler name "<integration>", because it is already used by *fixtures.IntegrationMessageHandler`,
		// 				),
		// 			))
		// 		})

		// 		It("returns an error when a projection key is in conflict", func() {
		// 			projection.ConfigureFunc = func(c dogma.ProjectionConfigurer) {
		// 				c.Identity("<projection>", "<integration-key>") // conflict!
		// 				c.ConsumesEventType(fixtures.MessageD{})
		// 			}

		// 			_, err := FromApplication(app)

		// 			Expect(err).To(Equal(
		// 				Error(
		// 					`*fixtures.ProjectionMessageHandler can not use the handler key "<integration-key>", because it is already used by *fixtures.IntegrationMessageHandler`,
		// 				),
		// 			))
		// 		})
		// 	})

		// 	It("returns an error when the app contains multiple consumers of the same command", func() {
		// 		integration.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
		// 			c.Identity("<integration>", "<integration-key>")
		// 			c.ConsumesCommandType(fixtures.MessageA{}) // conflict with <aggregate>
		// 			c.ProducesEventType(fixtures.MessageF{})
		// 		}

		// 		_, err := FromApplication(app)

		// 		Expect(err).To(Equal(
		// 			Error(
		// 				`the "<integration>" handler can not consume fixtures.MessageA commands because they are already consumed by "<aggregate>"`,
		// 			),
		// 		))
		// 	})

		// 	It("returns an error when the app contains multiple producers of the same event", func() {
		// 		integration.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
		// 			c.Identity("<integration>", "<integration-key>")
		// 			c.ConsumesCommandType(fixtures.MessageC{})
		// 			c.ProducesEventType(fixtures.MessageE{}) // conflict with <aggregate>
		// 		}

		// 		_, err := FromApplication(app)

		// 		Expect(err).To(Equal(
		// 			Error(
		// 				`the "<integration>" handler can not produce fixtures.MessageE events because they are already produced by "<aggregate>"`,
		// 			),
		// 		))
		// 	})

		// 	It("does not return an error when the app contains multiple processes that schedule the same timeout", func() {
		// 		process1 := &fixtures.ProcessMessageHandler{
		// 			ConfigureFunc: func(c dogma.ProcessConfigurer) {
		// 				c.Identity("<process-1>", "<process-1-key>")
		// 				c.ConsumesEventType(fixtures.MessageB{})
		// 				c.ProducesCommandType(fixtures.MessageC{})
		// 				c.SchedulesTimeoutType(fixtures.MessageT{})
		// 			},
		// 		}

		// 		process2 := &fixtures.ProcessMessageHandler{
		// 			ConfigureFunc: func(c dogma.ProcessConfigurer) {
		// 				c.Identity("<process-2>", "<process-2-key>")
		// 				c.ConsumesEventType(fixtures.MessageB{})
		// 				c.ProducesCommandType(fixtures.MessageC{})
		// 				c.SchedulesTimeoutType(fixtures.MessageT{})
		// 			},
		// 		}

		// 		app := &fixtures.Application{
		// 			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
		// 				c.Identity("<app>", "<app-key>")
		// 				c.RegisterProcess(process1)
		// 				c.RegisterProcess(process2)
		// 			},
		// 		}

		// 		_, err := FromApplication(app)

		// 		Expect(err).ShouldNot(HaveOccurred())
		// 	})

		// 	When("multiple handlers use a single message type in differing roles", func() {
		// 		It("returns an error when a conflict occurs with a consumed message", func() {
		// 			process.ConfigureFunc = func(c dogma.ProcessConfigurer) {
		// 				c.Identity("<process>", "<process-key>")
		// 				c.ConsumesEventType(fixtures.MessageA{}) // conflict with <aggregate>
		// 				c.ProducesCommandType(fixtures.MessageC{})
		// 			}

		// 			_, err := FromApplication(app)

		// 			Expect(err).To(Equal(
		// 				Error(
		// 					`the "<process>" handler configures fixtures.MessageA as an event but "<aggregate>" configures it as a command`,
		// 				),
		// 			))
		// 		})

		// 		It("returns an error when a conflict occurs with a produced message", func() {
		// 			process.ConfigureFunc = func(c dogma.ProcessConfigurer) {
		// 				c.Identity("<process>", "<process-key>")
		// 				c.ConsumesEventType(fixtures.MessageB{})
		// 				c.ProducesCommandType(fixtures.MessageE{}) // conflict with <aggregate>
		// 			}

		// 			_, err := FromApplication(app)

		// 			Expect(err).To(Equal(
		// 				Error(
		// 					`the "<process>" handler configures fixtures.MessageE as a command but "<aggregate>" configures it as an event`,
		// 				),
		// 			))
		// 		})
		// 	})
	)
})
