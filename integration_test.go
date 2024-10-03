package configkit_test

import (
	"context"
	"errors"
	"reflect"

	. "github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func FromIntegration()", func() {
	var handler *IntegrationMessageHandlerStub

	BeforeEach(func() {
		handler = &IntegrationMessageHandlerStub{
			ConfigureFunc: func(c dogma.IntegrationConfigurer) {
				c.Identity("<name>", integrationKey)
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.HandlesCommand[CommandStub[TypeB]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			},
		}
	})

	When("the configuration is valid", func() {
		var cfg RichIntegration

		JustBeforeEach(func() {
			cfg = FromIntegration(handler)
		})

		Describe("func Identity()", func() {
			It("returns the handler identity", func() {
				Expect(cfg.Identity()).To(Equal(
					MustNewIdentity("<name>", integrationKey),
				))
			})
		})

		Describe("func MessageNames()", func() {
			It("returns the expected message names", func() {
				Expect(cfg.MessageNames()).To(Equal(
					EntityMessages[message.Name]{
						message.NameOf(EventA1): {
							Kind:       message.EventKind,
							IsProduced: true,
						},
						message.NameOf(CommandA1): {
							Kind:       message.CommandKind,
							IsConsumed: true,
						},
						message.NameOf(CommandB1): {
							Kind:       message.CommandKind,
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
						message.TypeOf(EventA1): {
							Kind:       message.EventKind,
							IsProduced: true,
						},
						message.TypeOf(CommandA1): {
							Kind:       message.CommandKind,
							IsConsumed: true,
						},
						message.TypeOf(CommandB1): {
							Kind:       message.CommandKind,
							IsConsumed: true,
						},
					},
				))
			})
		})

		Describe("func TypeName()", func() {
			It("returns the fully-qualified type name of the handler", func() {
				Expect(cfg.TypeName()).To(Equal("*github.com/dogmatiq/enginekit/enginetest/stubs.IntegrationMessageHandlerStub"))
			})
		})

		Describe("func ReflectType()", func() {
			It("returns the type of the handler", func() {
				Expect(cfg.ReflectType()).To(Equal(reflect.TypeOf(handler)))
			})
		})

		Describe("func AcceptVisitor()", func() {
			It("calls the appropriate method on the visitor", func() {
				v := &visitorStub{
					VisitIntegrationFunc: func(_ context.Context, c Integration) error {
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
					VisitRichIntegrationFunc: func(_ context.Context, c RichIntegration) error {
						Expect(c).To(BeIdenticalTo(cfg))
						return errors.New("<error>")
					},
				}

				err := cfg.AcceptRichVisitor(context.Background(), v)
				Expect(err).To(MatchError("<error>"))
			})
		})

		Describe("func HandlerType()", func() {
			It("returns IntegrationHandlerType", func() {
				Expect(cfg.HandlerType()).To(Equal(IntegrationHandlerType))
			})
		})

		Describe("func Handler()", func() {
			It("returns the underlying handler", func() {
				Expect(cfg.Handler()).To(BeIdenticalTo(handler))
			})
		})

		When("the handler is disabled", func() {
			BeforeEach(func() {
				configure := handler.ConfigureFunc
				handler.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
					configure(c)
					c.Disable()
				}
			})

			Describe("func IsDisabled()", func() {
				It("returns true", func() {
					Expect(cfg.IsDisabled()).To(BeTrue())
				})
			})
		})

		When("the handler does not configure any event routes", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
					c.Identity("<name>", integrationKey)
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
					)
				}
			})

			It("does not panic", func() {
				FromIntegration(handler)
			})
		})
	})

	DescribeTable(
		"when the configuration is invalid",
		func(
			msg string,
			fn func(dogma.IntegrationConfigurer),
		) {
			handler.ConfigureFunc = fn

			var err error
			func() {
				defer Recover(&err)
				FromIntegration(handler)
			}()

			Expect(err).Should(HaveOccurred())
			if msg != "" {
				Expect(err).To(MatchError(msg))
			}
		},
		Entry(
			"when the handler does not configure anything",
			"", // any error
			nil,
		),
		Entry(
			"when the handler does not configure an identity",
			`*stubs.IntegrationMessageHandlerStub is configured without an identity, Identity() must be called exactly once within Configure()`,
			func(c dogma.IntegrationConfigurer) {
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler configures multiple identities",
			`*stubs.IntegrationMessageHandlerStub is configured with multiple identities (<name>/e28f056e-e5a0-4ee7-aaf1-1d1fe02fb6e3 and <other>/e28f056e-e5a0-4ee7-aaf1-1d1fe02fb6e3), Identity() must be called exactly once within Configure()`,
			func(c dogma.IntegrationConfigurer) {
				c.Identity("<name>", integrationKey)
				c.Identity("<other>", integrationKey)
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler configures an invalid name",
			`*stubs.IntegrationMessageHandlerStub is configured with an invalid identity, invalid name "\t \n", names must be non-empty, printable UTF-8 strings with no whitespace`,
			func(c dogma.IntegrationConfigurer) {
				c.Identity("\t \n", integrationKey)
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler configures an invalid key",
			`*stubs.IntegrationMessageHandlerStub is configured with an invalid identity, invalid key "\t \n", keys must be RFC 4122 UUIDs`,
			func(c dogma.IntegrationConfigurer) {
				c.Identity("<name>", "\t \n")
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler does not configure any command routes",
			`*stubs.IntegrationMessageHandlerStub (<name>) is not configured to handle any commands, at least one HandlesCommand() route must be added within Configure()`,
			func(c dogma.IntegrationConfigurer) {
				c.Identity("<name>", integrationKey)
				c.Routes(
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler configures multiple routes for the same command",
			`*stubs.IntegrationMessageHandlerStub (<name>) is configured with multiple HandlesCommand() routes for stubs.CommandStub[TypeA], should these refer to different message types?`,
			func(c dogma.IntegrationConfigurer) {
				c.Identity("<name>", integrationKey)
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler configures multiple routes for the same event",
			`*stubs.IntegrationMessageHandlerStub (<name>) is configured with multiple RecordsEvent() routes for stubs.EventStub[TypeA], should these refer to different message types?`,
			func(c dogma.IntegrationConfigurer) {
				c.Identity("<name>", integrationKey)
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			},
		),
		Entry(
			"when an error occurs before the identity is configured it omits the handler name",
			`*stubs.IntegrationMessageHandlerStub is configured with multiple HandlesCommand() routes for stubs.CommandStub[TypeA], should these refer to different message types?`,
			func(c dogma.IntegrationConfigurer) {
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.HandlesCommand[CommandStub[TypeA]](),
				)
			},
		),
	)
})
