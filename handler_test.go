package configkit_test

import (
	. "github.com/dogmatiq/configkit"
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func IsHandlerEqual()", func() {
	It("returns true if the two handlers are equivalent", func() {
		h := &AggregateMessageHandlerStub{
			ConfigureFunc: func(c dogma.AggregateConfigurer) {
				c.Identity("<name>", aggregateKey)
				c.Routes(
					dogma.HandlesCommand[*CommandStub[TypeA]](),
					dogma.HandlesCommand[*CommandStub[TypeB]](),
					dogma.RecordsEvent[*EventStub[TypeA]](),
				)
			},
		}

		a := FromAggregate(h)
		b := FromAggregate(h)

		Expect(IsHandlerEqual(a, b)).To(BeTrue())
	})

	DescribeTable(
		"returns false if the handlers are not equivalent",
		func(b Handler) {
			h := &AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("<name>", aggregateKey)
					c.Routes(
						dogma.HandlesCommand[*CommandStub[TypeA]](),
						dogma.HandlesCommand[*CommandStub[TypeB]](),
						dogma.RecordsEvent[*EventStub[TypeA]](),
					)
				},
			}

			a := FromAggregate(h)

			Expect(IsHandlerEqual(a, b)).To(BeFalse())
		},
		Entry(
			"handler type differs",
			FromIntegration(&IntegrationMessageHandlerStub{ // diff
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("<name>", integrationKey)
					c.Routes(
						dogma.HandlesCommand[*CommandStub[TypeA]](),
						dogma.HandlesCommand[*CommandStub[TypeB]](),
						dogma.RecordsEvent[*EventStub[TypeA]](),
					)
				},
			}),
		),
		Entry(
			"identity name differs",
			FromAggregate(&AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("<name-different>", aggregateKey) // diff
					c.Routes(
						dogma.HandlesCommand[*CommandStub[TypeA]](),
						dogma.HandlesCommand[*CommandStub[TypeB]](),
						dogma.RecordsEvent[*EventStub[TypeA]](),
					)
				},
			}),
		),
		Entry(
			"identity key differs",
			FromAggregate(&AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("<name>", "799239e7-8c03-48f9-a324-14b7f9b76e30") // diff
					c.Routes(
						dogma.HandlesCommand[*CommandStub[TypeA]](),
						dogma.HandlesCommand[*CommandStub[TypeB]](),
						dogma.RecordsEvent[*EventStub[TypeA]](),
					)
				},
			}),
		),
		Entry(
			"disabled state differs",
			FromAggregate(&AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("<name>", aggregateKey)
					c.Routes(
						dogma.HandlesCommand[*CommandStub[TypeA]](),
						dogma.HandlesCommand[*CommandStub[TypeB]](),
						dogma.RecordsEvent[*EventStub[TypeA]](),
					)
					c.Disable()
				},
			}),
		),
		Entry(
			"messages differ",
			FromAggregate(&AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("<name>", aggregateKey)
					c.Routes(
						dogma.HandlesCommand[*CommandStub[TypeA]](),
						dogma.HandlesCommand[*CommandStub[TypeC]](), // diff
						dogma.RecordsEvent[*EventStub[TypeA]](),
					)
				},
			}),
		),
	)
})
