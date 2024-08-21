package configkit_test

import (
	"context"
	"errors"
	"reflect"

	. "github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/fixtures" // can't dot-import due to conflicts
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func FromProcess()", func() {
	var handler *ProcessMessageHandlerStub

	BeforeEach(func() {
		handler = &ProcessMessageHandlerStub{
			ConfigureFunc: func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", processKey)
				c.Routes(
					dogma.HandlesEvent[EventStub[TypeA]](),
					dogma.HandlesEvent[EventStub[TypeB]](),
					dogma.ExecutesCommand[CommandStub[TypeA]](),
					dogma.SchedulesTimeout[TimeoutStub[TypeA]](),
				)
			},
		}
	})

	When("the configuration is valid", func() {
		var cfg RichProcess

		JustBeforeEach(func() {
			cfg = FromProcess(handler)
		})

		Describe("func Identity()", func() {
			It("returns the handler identity", func() {
				Expect(cfg.Identity()).To(Equal(
					MustNewIdentity("<name>", processKey),
				))
			})
		})

		Describe("func MessageNames()", func() {
			It("returns the expected message names", func() {
				Expect(cfg.MessageNames()).To(Equal(
					EntityMessageNames{
						Produced: message.NameRoles{
							message.NameFor[CommandStub[TypeA]](): message.CommandRole,
							message.NameFor[TimeoutStub[TypeA]](): message.TimeoutRole,
						},
						Consumed: message.NameRoles{
							message.NameFor[EventStub[TypeA]]():   message.EventRole,
							message.NameFor[EventStub[TypeB]]():   message.EventRole,
							message.NameFor[TimeoutStub[TypeA]](): message.TimeoutRole,
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
							message.TypeFor[CommandStub[TypeA]](): message.CommandRole,
							message.TypeFor[TimeoutStub[TypeA]](): message.TimeoutRole,
						},
						Consumed: message.TypeRoles{
							message.TypeFor[EventStub[TypeA]]():   message.EventRole,
							message.TypeFor[EventStub[TypeB]]():   message.EventRole,
							message.TypeFor[TimeoutStub[TypeA]](): message.TimeoutRole,
						},
					},
				))
			})
		})

		Describe("func TypeName()", func() {
			It("returns the fully-qualified type name of the handler", func() {
				Expect(cfg.TypeName()).To(Equal("*github.com/dogmatiq/enginekit/enginetest/stubs.ProcessMessageHandlerStub"))
			})
		})

		Describe("func ReflectType()", func() {
			It("returns the type of the handler", func() {
				Expect(cfg.ReflectType()).To(Equal(reflect.TypeOf(handler)))
			})
		})

		Describe("func AcceptVisitor()", func() {
			It("calls the appropriate method on the visitor", func() {
				v := &fixtures.Visitor{
					VisitProcessFunc: func(_ context.Context, c Process) error {
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
				v := &fixtures.RichVisitor{
					VisitRichProcessFunc: func(_ context.Context, c RichProcess) error {
						Expect(c).To(BeIdenticalTo(cfg))
						return errors.New("<error>")
					},
				}

				err := cfg.AcceptRichVisitor(context.Background(), v)
				Expect(err).To(MatchError("<error>"))
			})
		})

		Describe("func HandlerType()", func() {
			It("returns ProcessHandlerType", func() {
				Expect(cfg.HandlerType()).To(Equal(ProcessHandlerType))
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
				handler.ConfigureFunc = func(c dogma.ProcessConfigurer) {
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
			fn func(dogma.ProcessConfigurer),
		) {
			handler.ConfigureFunc = fn

			var err error
			func() {
				defer Recover(&err)
				FromProcess(handler)
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
			`*stubs.ProcessMessageHandlerStub is configured without an identity, Identity() must be called exactly once within Configure()`,
			func(c dogma.ProcessConfigurer) {
				c.Routes(
					dogma.HandlesEvent[EventStub[TypeA]](),
					dogma.ExecutesCommand[CommandStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler configures multiple identities",
			`*stubs.ProcessMessageHandlerStub is configured with multiple identities (<name>/bea52cf4-e403-4b18-819d-88ade7836308 and <other>/bea52cf4-e403-4b18-819d-88ade7836308), Identity() must be called exactly once within Configure()`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", processKey)
				c.Identity("<other>", processKey)
				c.Routes(
					dogma.HandlesEvent[EventStub[TypeA]](),
					dogma.ExecutesCommand[CommandStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler configures an invalid name",
			`*stubs.ProcessMessageHandlerStub is configured with an invalid identity, invalid name "\t \n", names must be non-empty, printable UTF-8 strings with no whitespace`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("\t \n", processKey)
				c.Routes(
					dogma.HandlesEvent[EventStub[TypeA]](),
					dogma.ExecutesCommand[CommandStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler configures an invalid key",
			`*stubs.ProcessMessageHandlerStub is configured with an invalid identity, invalid key "\t \n", keys must be RFC 4122 UUIDs`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", "\t \n")
				c.Routes(
					dogma.HandlesEvent[EventStub[TypeA]](),
					dogma.ExecutesCommand[CommandStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler does not configure any event routes",
			`*stubs.ProcessMessageHandlerStub (<name>) is not configured to handle any events, at least one HandlesEvent() route must be added within Configure()`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", processKey)
				c.Routes(
					dogma.ExecutesCommand[CommandStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler configures multiple routes for the same event",
			`*stubs.ProcessMessageHandlerStub (<name>) is configured with multiple HandlesEvent() routes for stubs.EventStub[TypeA], should these refer to different message types?`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", processKey)
				c.Routes(
					dogma.HandlesEvent[EventStub[TypeA]](),
					dogma.HandlesEvent[EventStub[TypeA]](),
					dogma.ExecutesCommand[CommandStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler does not configure command routes",
			`*stubs.ProcessMessageHandlerStub (<name>) is not configured to execute any commands, at least one ExecutesCommand() route must be added within Configure()`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", processKey)
				c.Routes(
					dogma.HandlesEvent[EventStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler configures multiple routes for the same command",
			`*stubs.ProcessMessageHandlerStub (<name>) is configured with multiple ExecutesCommand() routes for stubs.CommandStub[TypeA], should these refer to different message types?`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", processKey)
				c.Routes(
					dogma.HandlesEvent[EventStub[TypeA]](),
					dogma.ExecutesCommand[CommandStub[TypeA]](),
					dogma.ExecutesCommand[CommandStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler configures multiple routes for the same timeout",
			`*stubs.ProcessMessageHandlerStub (<name>) is configured with multiple SchedulesTimeout() routes for stubs.TimeoutStub[TypeA], should these refer to different message types?`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", processKey)
				c.Routes(
					dogma.HandlesEvent[EventStub[TypeA]](),
					dogma.ExecutesCommand[CommandStub[TypeA]](),
					dogma.SchedulesTimeout[TimeoutStub[TypeA]](),
					dogma.SchedulesTimeout[TimeoutStub[TypeA]](),
				)
			},
		),
		Entry(
			"when the handler configures the same message type with different roles",
			`*stubs.ProcessMessageHandlerStub (<name>) is configured to use stubs.EventStub[TypeA] as both an event and a timeout`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", processKey)
				c.Routes(
					dogma.HandlesEvent[EventStub[TypeA]](),
					dogma.SchedulesTimeout[EventStub[TypeA]](),
				)
			},
		),
		Entry(
			"when an error occurs before the identity is configured it omits the handler name",
			`*stubs.ProcessMessageHandlerStub is configured with multiple HandlesEvent() routes for stubs.EventStub[TypeA], should these refer to different message types?`,
			func(c dogma.ProcessConfigurer) {
				c.Routes(
					dogma.HandlesEvent[EventStub[TypeA]](),
					dogma.HandlesEvent[EventStub[TypeA]](),
				)
			},
		),
	)
})
