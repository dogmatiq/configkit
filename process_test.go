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

var _ = Describe("func FromProcess()", func() {
	var handler *fixtures.ProcessMessageHandler

	BeforeEach(func() {
		handler = &fixtures.ProcessMessageHandler{
			ConfigureFunc: func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesEventType(fixtures.MessageA{})
				c.ConsumesEventType(fixtures.MessageB{})
				c.ProducesCommandType(fixtures.MessageC{})
				c.SchedulesTimeoutType(fixtures.MessageT{})
			},
		}
	})

	When("the configuration is valid", func() {
		var cfg RichProcess

		BeforeEach(func() {
			cfg = FromProcess(handler)
		})

		Describe("func Identity()", func() {
			It("returns the handler identity", func() {
				Expect(cfg.Identity()).To(Equal(
					MustNewIdentity("<name>", "<key>"),
				))
			})
		})

		Describe("func MessageNames()", func() {
			It("returns the expected message names", func() {
				Expect(cfg.MessageNames()).To(Equal(
					EntityMessageNames{
						Roles: message.NameRoles{
							cfixtures.MessageATypeName: message.EventRole,
							cfixtures.MessageBTypeName: message.EventRole,
							cfixtures.MessageCTypeName: message.CommandRole,
							cfixtures.MessageTTypeName: message.TimeoutRole,
						},
						Produced: message.NameRoles{
							cfixtures.MessageCTypeName: message.CommandRole,
							cfixtures.MessageTTypeName: message.TimeoutRole,
						},
						Consumed: message.NameRoles{
							cfixtures.MessageATypeName: message.EventRole,
							cfixtures.MessageBTypeName: message.EventRole,
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
						Roles: message.TypeRoles{
							cfixtures.MessageAType: message.EventRole,
							cfixtures.MessageBType: message.EventRole,
							cfixtures.MessageCType: message.CommandRole,
							cfixtures.MessageTType: message.TimeoutRole,
						},
						Produced: message.TypeRoles{
							cfixtures.MessageCType: message.CommandRole,
							cfixtures.MessageTType: message.TimeoutRole,
						},
						Consumed: message.TypeRoles{
							cfixtures.MessageAType: message.EventRole,
							cfixtures.MessageBType: message.EventRole,
							cfixtures.MessageTType: message.TimeoutRole,
						},
					},
				))
			})
		})

		Describe("func TypeName()", func() {
			It("returns the fully-qualified type name of the handler", func() {
				Expect(cfg.TypeName()).To(Equal("*github.com/dogmatiq/dogma/fixtures.ProcessMessageHandler"))
			})
		})

		Describe("func ReflectType()", func() {
			It("returns the type of the handler", func() {
				Expect(cfg.ReflectType()).To(Equal(reflect.TypeOf(handler)))
			})
		})

		Describe("func AcceptVisitor()", func() {
			It("calls the appropriate method on the visitor", func() {
				v := &cfixtures.Visitor{
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
				v := &cfixtures.RichVisitor{
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
			`*fixtures.ProcessMessageHandler is configured without an identity, Identity() must be called exactly once within Configure()`,
			func(c dogma.ProcessConfigurer) {
				c.ConsumesEventType(fixtures.MessageA{})
				c.ProducesCommandType(fixtures.MessageC{})
			},
		),
		Entry(
			"when the handler configures multiple identities",
			`*fixtures.ProcessMessageHandler is configured with multiple identities (<name>/<key> and <other>/<key>), Identity() must be called exactly once within Configure()`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", "<key>")
				c.Identity("<other>", "<key>")
				c.ConsumesEventType(fixtures.MessageA{})
				c.ProducesCommandType(fixtures.MessageC{})
			},
		),
		Entry(
			"when the handler configures an invalid name",
			`*fixtures.ProcessMessageHandler is configured with an invalid identity, invalid name "\t \n", names must be non-empty, printable UTF-8 strings with no whitespace`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("\t \n", "<key>")
				c.ConsumesEventType(fixtures.MessageA{})
				c.ProducesCommandType(fixtures.MessageC{})
			},
		),
		Entry(
			"when the handler configures an invalid key",
			`*fixtures.ProcessMessageHandler is configured with an invalid identity, invalid key "\t \n", keys must be non-empty, printable UTF-8 strings with no whitespace`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", "\t \n")
				c.ConsumesEventType(fixtures.MessageA{})
				c.ProducesCommandType(fixtures.MessageC{})
			},
		),
		Entry(
			"when the handler does not configure any consumed event types",
			`*fixtures.ProcessMessageHandler is not configured to consume any events, ConsumesEventType() must be called at least once within Configure()`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", "<key>")
				c.ProducesCommandType(fixtures.MessageC{})
			},
		),
		Entry(
			"when the handler configures the same consumed event type multiple times",
			`*fixtures.ProcessMessageHandler is configured to consume the fixtures.MessageA event more than once, should this refer to different message types?`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesEventType(fixtures.MessageA{})
				c.ConsumesEventType(fixtures.MessageA{})
				c.ProducesCommandType(fixtures.MessageC{})
			},
		),
		Entry(
			"when the handler does not configure any produced commands",
			`*fixtures.ProcessMessageHandler is not configured to produce any commands, ProducesCommandType() must be called at least once within Configure()`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesEventType(fixtures.MessageA{})
			},
		),
		Entry(
			"when the handler configures the same produced command type multiple times",
			`*fixtures.ProcessMessageHandler is configured to produce the fixtures.MessageC command more than once, should this refer to different message types?`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesEventType(fixtures.MessageA{})
				c.ProducesCommandType(fixtures.MessageC{})
				c.ProducesCommandType(fixtures.MessageC{})
			},
		),
		Entry(
			"when the handler configures the same scheudled timeout type multiple times",
			`*fixtures.ProcessMessageHandler is configured to schedule the fixtures.MessageT timeout more than once, should this refer to different message types?`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesEventType(fixtures.MessageA{})
				c.ProducesCommandType(fixtures.MessageC{})
				c.SchedulesTimeoutType(fixtures.MessageT{})
				c.SchedulesTimeoutType(fixtures.MessageT{})
			},
		),
		Entry(
			"when the handler configures the same message type with different roles",
			`*fixtures.ProcessMessageHandler is configured to use fixtures.MessageA as both an event and a timeout`,
			func(c dogma.ProcessConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesEventType(fixtures.MessageA{})
				c.SchedulesTimeoutType(fixtures.MessageA{})
			},
		),
	)
})
