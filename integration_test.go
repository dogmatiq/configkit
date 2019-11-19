package config_test

import (
	. "github.com/dogmatiq/configkit"
	. "github.com/dogmatiq/configkit/fixtures"
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/dogma/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ HandlerConfig = &IntegrationConfig{}

var _ = Describe("type IntegrationConfig", func() {
	Describe("func NewIntegrationConfig", func() {
		var handler *IntegrationMessageHandler

		BeforeEach(func() {
			handler = &IntegrationMessageHandler{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("<name>", "<key>")
					c.ConsumesCommandType(MessageA{})
					c.ConsumesCommandType(MessageB{})
					c.ProducesEventType(MessageE{})
				},
			}
		})

		When("the configuration is valid", func() {
			var cfg *IntegrationConfig

			BeforeEach(func() {
				var err error
				cfg, err = NewIntegrationConfig(handler)
				Expect(err).ShouldNot(HaveOccurred())
			})

			Describe("func Identity()", func() {
				It("returns the handler identity", func() {
					Expect(cfg.Identity()).To(Equal(
						MustNewIdentity("<name>", "<key>"),
					))
				})
			})

			Describe("func HandlerType()", func() {
				It("returns IntegrationHandlerType", func() {
					Expect(cfg.HandlerType()).To(Equal(IntegrationHandlerType))
				})
			})

			Describe("func ConsumedMessageTypes()", func() {
				It("returns the expected message types", func() {
					Expect(cfg.ConsumedMessageTypes()).To(Equal(
						MessageRoleMap{
							MessageAType: CommandMessageRole,
							MessageBType: CommandMessageRole,
						},
					))
				})
			})

			Describe("func ProducedMessageTypes()", func() {
				It("returns the expected message types", func() {
					Expect(cfg.ProducedMessageTypes()).To(Equal(
						MessageRoleMap{
							MessageEType: EventMessageRole,
						},
					))
				})
			})

			When("the handler does not configure any produced events", func() {
				BeforeEach(func() {
					handler.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
						c.Identity("<name>", "<key>")
						c.ConsumesCommandType(MessageA{})
					}
				})

				It("does not return an error", func() {
					_, err := NewIntegrationConfig(handler)
					Expect(err).ShouldNot(HaveOccurred())
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

				_, err := NewIntegrationConfig(handler)
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
				`*fixtures.IntegrationMessageHandler.Configure() did not call IntegrationConfigurer.Identity()`,
				func(c dogma.IntegrationConfigurer) {
					c.ConsumesCommandType(MessageA{})
					c.ProducesEventType(MessageE{})
				},
			),
			Entry(
				"when the handler configures multiple identities",
				`*fixtures.IntegrationMessageHandler.Configure() has already called IntegrationConfigurer.Identity("<name>", "<key>")`,
				func(c dogma.IntegrationConfigurer) {
					c.Identity("<name>", "<key>")
					c.Identity("<other>", "<key>")
					c.ConsumesCommandType(MessageA{})
					c.ProducesEventType(MessageE{})
				},
			),
			Entry(
				"when the handler configures an invalid name",
				`*fixtures.IntegrationMessageHandler.Configure() called IntegrationConfigurer.Identity() with an invalid name "\t \n", names must be non-empty, printable UTF-8 strings with no whitespace`,
				func(c dogma.IntegrationConfigurer) {
					c.Identity("\t \n", "<key>")
					c.ConsumesCommandType(MessageA{})
					c.ProducesEventType(MessageE{})
				},
			),
			Entry(
				"when the handler configures an invalid key",
				`*fixtures.IntegrationMessageHandler.Configure() called IntegrationConfigurer.Identity() with an invalid key "\t \n", keys must be non-empty, printable UTF-8 strings with no whitespace`,
				func(c dogma.IntegrationConfigurer) {
					c.Identity("<name>", "\t \n")
					c.ConsumesCommandType(MessageA{})
					c.ProducesEventType(MessageE{})
				},
			),
			Entry(
				"when the handler does not configure any consumed command types",
				`*fixtures.IntegrationMessageHandler.Configure() did not call IntegrationConfigurer.ConsumesCommandType()`,
				func(c dogma.IntegrationConfigurer) {
					c.Identity("<name>", "<key>")
					c.ProducesEventType(MessageE{})
				},
			),
			Entry(
				"when the handler configures the same consumed command type multiple times",
				`*fixtures.IntegrationMessageHandler.Configure() has already called IntegrationConfigurer.ConsumesCommandType(fixtures.MessageA)`,
				func(c dogma.IntegrationConfigurer) {
					c.Identity("<name>", "<key>")
					c.ConsumesCommandType(MessageA{})
					c.ConsumesCommandType(MessageA{})
					c.ProducesEventType(MessageE{})
				},
			),
			Entry(
				"when the handler configures the same produced event type multiple times",
				`*fixtures.IntegrationMessageHandler.Configure() has already called IntegrationConfigurer.ProducesEventType(fixtures.MessageE)`,
				func(c dogma.IntegrationConfigurer) {
					c.Identity("<name>", "<key>")
					c.ConsumesCommandType(MessageA{})
					c.ProducesEventType(MessageE{})
					c.ProducesEventType(MessageE{})
				},
			),
		)
	})
})
