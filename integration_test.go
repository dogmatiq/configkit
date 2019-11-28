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

var _ = Describe("func FromIntegration()", func() {
	var handler *fixtures.IntegrationMessageHandler

	BeforeEach(func() {
		handler = &fixtures.IntegrationMessageHandler{
			ConfigureFunc: func(c dogma.IntegrationConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ConsumesCommandType(fixtures.MessageB{})
				c.ProducesEventType(fixtures.MessageE{})
			},
		}
	})

	When("the configuration is valid", func() {
		var cfg RichIntegration

		BeforeEach(func() {
			cfg = FromIntegration(handler)
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
							cfixtures.MessageATypeName: message.CommandRole,
							cfixtures.MessageBTypeName: message.CommandRole,
							cfixtures.MessageETypeName: message.EventRole,
						},
						Produced: message.NameRoles{
							cfixtures.MessageETypeName: message.EventRole,
						},
						Consumed: message.NameRoles{
							cfixtures.MessageATypeName: message.CommandRole,
							cfixtures.MessageBTypeName: message.CommandRole,
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
							cfixtures.MessageAType: message.CommandRole,
							cfixtures.MessageBType: message.CommandRole,
							cfixtures.MessageEType: message.EventRole,
						},
						Produced: message.TypeRoles{
							cfixtures.MessageEType: message.EventRole,
						},
						Consumed: message.TypeRoles{
							cfixtures.MessageAType: message.CommandRole,
							cfixtures.MessageBType: message.CommandRole,
						},
					},
				))
			})
		})

		Describe("func TypeName()", func() {
			It("returns the fully-qualified type name of the handler", func() {
				Expect(cfg.TypeName()).To(Equal("*github.com/dogmatiq/dogma/fixtures.IntegrationMessageHandler"))
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
				v := &cfixtures.RichVisitor{
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

		When("the handler does not configure any produced events", func() {
			BeforeEach(func() {
				handler.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
					c.Identity("<name>", "<key>")
					c.ConsumesCommandType(fixtures.MessageA{})
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
			`*fixtures.IntegrationMessageHandler is configured without an identity, Identity() must be called exactly once within Configure()`,
			func(c dogma.IntegrationConfigurer) {
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ProducesEventType(fixtures.MessageE{})
			},
		),
		Entry(
			"when the handler configures multiple identities",
			`*fixtures.IntegrationMessageHandler is configured with multiple identities (<name>/<key> and <other>/<key>), Identity() must be called exactly once within Configure()`,
			func(c dogma.IntegrationConfigurer) {
				c.Identity("<name>", "<key>")
				c.Identity("<other>", "<key>")
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ProducesEventType(fixtures.MessageE{})
			},
		),
		Entry(
			"when the handler configures an invalid name",
			`*fixtures.IntegrationMessageHandler is configured with an invalid identity, invalid name "\t \n", names must be non-empty, printable UTF-8 strings with no whitespace`,
			func(c dogma.IntegrationConfigurer) {
				c.Identity("\t \n", "<key>")
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ProducesEventType(fixtures.MessageE{})
			},
		),
		Entry(
			"when the handler configures an invalid key",
			`*fixtures.IntegrationMessageHandler is configured with an invalid identity, invalid key "\t \n", keys must be non-empty, printable UTF-8 strings with no whitespace`,
			func(c dogma.IntegrationConfigurer) {
				c.Identity("<name>", "\t \n")
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ProducesEventType(fixtures.MessageE{})
			},
		),
		Entry(
			"when the handler does not configure any consumed command types",
			`*fixtures.IntegrationMessageHandler (<name>) is not configured to consume any commands, ConsumesCommandType() must be called at least once within Configure()`,
			func(c dogma.IntegrationConfigurer) {
				c.Identity("<name>", "<key>")
				c.ProducesEventType(fixtures.MessageE{})
			},
		),
		Entry(
			"when the handler configures the same consumed command type multiple times",
			`*fixtures.IntegrationMessageHandler (<name>) is configured to consume the fixtures.MessageA command more than once, should this refer to different message types?`,
			func(c dogma.IntegrationConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ProducesEventType(fixtures.MessageE{})
			},
		),
		Entry(
			"when the handler configures the same produced event type multiple times",
			`*fixtures.IntegrationMessageHandler (<name>) is configured to produce the fixtures.MessageE event more than once, should this refer to different message types?`,
			func(c dogma.IntegrationConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ProducesEventType(fixtures.MessageE{})
				c.ProducesEventType(fixtures.MessageE{})
			},
		),
		Entry(
			"when the handler configures the same message type with different roles",
			`*fixtures.IntegrationMessageHandler (<name>) is configured to use fixtures.MessageA as both a command and an event`,
			func(c dogma.IntegrationConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ProducesEventType(fixtures.MessageA{})
			},
		),
	)
})
