package configkit_test

import (
	. "github.com/dogmatiq/configkit"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures" // can't dot-import due to conflicts
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
					dogma.HandlesCommand[fixtures.MessageA](),
					dogma.HandlesCommand[fixtures.MessageB](),
					dogma.RecordsEvent[fixtures.MessageE](),
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
						dogma.HandlesCommand[fixtures.MessageA](),
						dogma.HandlesCommand[fixtures.MessageB](),
						dogma.RecordsEvent[fixtures.MessageE](),
					)
				},
			}

			a := FromAggregate(h)

			Expect(IsHandlerEqual(a, b)).To(BeFalse())
		},
		Entry(
			"type differs",
			FromIntegration(&IntegrationMessageHandlerStub{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("<name>", integrationKey)
					c.Routes(
						dogma.HandlesCommand[fixtures.MessageA](),
						dogma.RecordsEvent[fixtures.MessageB](), // diff
						dogma.RecordsEvent[fixtures.MessageE](),
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
						dogma.HandlesCommand[fixtures.MessageA](),
						dogma.HandlesCommand[fixtures.MessageB](),
						dogma.RecordsEvent[fixtures.MessageE](),
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
						dogma.HandlesCommand[fixtures.MessageA](),
						dogma.HandlesCommand[fixtures.MessageB](),
						dogma.RecordsEvent[fixtures.MessageE](),
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
						dogma.HandlesCommand[fixtures.MessageA](),
						dogma.HandlesCommand[fixtures.MessageB](),
						dogma.RecordsEvent[fixtures.MessageE](),
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
						dogma.HandlesCommand[fixtures.MessageA](),
						dogma.RecordsEvent[fixtures.MessageB](), // diff
						dogma.RecordsEvent[fixtures.MessageE](),
					)
				},
			}),
		),
	)
})
