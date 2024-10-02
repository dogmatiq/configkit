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

var _ = Describe("func FromAggregate()", func() {
	var handler *AggregateMessageHandlerStub

	BeforeEach(func() {
		handler = &AggregateMessageHandlerStub{
			ConfigureFunc: func(c dogma.AggregateConfigurer) {
				c.Identity("<name>", aggregateKey)
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.HandlesCommand[CommandStub[TypeB]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			},
		}
	})

	When("the configuration is valid", func() {
		var cfg RichAggregate

		JustBeforeEach(func() {
			cfg = FromAggregate(handler)
		})

		Describe("func Identity()", func() {
			It("returns the handler identity", func() {
				Expect(cfg.Identity()).To(Equal(
					MustNewIdentity("<name>", aggregateKey),
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
				Expect(cfg.TypeName()).To(Equal("*github.com/dogmatiq/enginekit/enginetest/stubs.AggregateMessageHandlerStub"))
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
					VisitAggregateFunc: func(_ context.Context, c Aggregate) error {
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
					VisitRichAggregateFunc: func(_ context.Context, c RichAggregate) error {
						Expect(c).To(BeIdenticalTo(cfg))
						return errors.New("<error>")
					},
				}

				err := cfg.AcceptRichVisitor(context.Background(), v)
				Expect(err).To(MatchError("<error>"))
			})
		})

		Describe("func HandlerType()", func() {
			It("returns AggregateHandlerType", func() {
				Expect(cfg.HandlerType()).To(Equal(AggregateHandlerType))
			})
		})

		Describe("func IsDisabled()", func() {
			It("returns false", func() {
				Expect(cfg.IsDisabled()).To(BeFalse())
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
				handler.ConfigureFunc = func(c dogma.AggregateConfigurer) {
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
	})

	DescribeTable(
		"when the configuration is invalid",
		func(
			msg string,
			fn func(dogma.AggregateConfigurer),
		) {
			handler.ConfigureFunc = fn

			var err error
			func() {
				defer Recover(&err)
				FromAggregate(handler)
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
			`*stubs.AggregateMessageHandlerStub is configured without an identity, Identity() must be called exactly once within Configure()`,
			func(c dogma.AggregateConfigurer) {
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler configures multiple identities",
			`*stubs.AggregateMessageHandlerStub is configured with multiple identities (<name>/14769f7f-87fe-48dd-916e-5bcab6ba6aca and <other>/14769f7f-87fe-48dd-916e-5bcab6ba6aca), Identity() must be called exactly once within Configure()`,
			func(c dogma.AggregateConfigurer) {
				c.Identity("<name>", aggregateKey)
				c.Identity("<other>", aggregateKey)
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler configures an invalid name",
			`*stubs.AggregateMessageHandlerStub is configured with an invalid identity, invalid name "\t \n", names must be non-empty, printable UTF-8 strings with no whitespace`,
			func(c dogma.AggregateConfigurer) {
				c.Identity("\t \n", appKey)
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler configures an invalid key",
			`*stubs.AggregateMessageHandlerStub is configured with an invalid identity, invalid key "\t \n", keys must be RFC 4122 UUIDs`,
			func(c dogma.AggregateConfigurer) {
				c.Identity("<name>", "\t \n")
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler does not configure any command routes",
			`*stubs.AggregateMessageHandlerStub (<name>) is not configured to handle any commands, at least one HandlesCommand() route must be added within Configure()`,
			func(c dogma.AggregateConfigurer) {
				c.Identity("<name>", aggregateKey)
				c.Routes(
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler configures multiple routes for the same command",
			`*stubs.AggregateMessageHandlerStub (<name>) is configured with multiple HandlesCommand() routes for stubs.CommandStub[TypeA], should these refer to different message types?`,
			func(c dogma.AggregateConfigurer) {
				c.Identity("<name>", aggregateKey)
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler does not configure any event routes",
			`*stubs.AggregateMessageHandlerStub (<name>) is not configured to record any events, at least one RecordsEvent() route must be added within Configure()`,
			func(c dogma.AggregateConfigurer) {
				c.Identity("<name>", aggregateKey)
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler configures multiple routes for the same event",
			`*stubs.AggregateMessageHandlerStub (<name>) is configured with multiple RecordsEvent() routes for stubs.EventStub[TypeA], should these refer to different message types?`,
			func(c dogma.AggregateConfigurer) {
				c.Identity("<name>", aggregateKey)
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
					dogma.RecordsEvent[EventStub[TypeA]](),
				)
			},
		),
		Entry(
			"when an error occurs before the identity is configured it omits the handler name",
			`*stubs.AggregateMessageHandlerStub is configured with multiple HandlesCommand() routes for stubs.CommandStub[TypeA], should these refer to different message types?`,
			func(c dogma.AggregateConfigurer) {
				c.Routes(
					dogma.HandlesCommand[CommandStub[TypeA]](),
					dogma.HandlesCommand[CommandStub[TypeA]](),
				)
			},
		),
	)
})
