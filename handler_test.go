package configkit_test

import (
	. "github.com/dogmatiq/configkit"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures" // can't dot-import due to conflicts
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func IsHandlerEqual()", func() {
	It("returns true if the two handlers are equivalent", func() {
		h := &fixtures.AggregateMessageHandler{
			ConfigureFunc: func(c dogma.AggregateConfigurer) {
				c.Identity("<name>", aggregateKey)
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ConsumesCommandType(fixtures.MessageB{})
				c.ProducesEventType(fixtures.MessageE{})
			},
		}

		a := FromAggregate(h)
		b := FromAggregate(h)

		Expect(IsHandlerEqual(a, b)).To(BeTrue())
	})

	DescribeTable(
		"returns false if the handlers are not equivalent",
		func(b Handler) {
			h := &fixtures.AggregateMessageHandler{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("<name>", aggregateKey)
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ConsumesCommandType(fixtures.MessageB{})
					c.ProducesEventType(fixtures.MessageE{})
				},
			}

			a := FromAggregate(h)

			Expect(IsHandlerEqual(a, b)).To(BeFalse())
		},
		Entry(
			"type differs",
			FromIntegration(&fixtures.IntegrationMessageHandler{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("<name>", integrationKey)
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageB{}) // diff
					c.ProducesEventType(fixtures.MessageE{})
				},
			}),
		),
		Entry(
			"identity name differs",
			FromAggregate(&fixtures.AggregateMessageHandler{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("<name-different>", aggregateKey) // diff
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ConsumesCommandType(fixtures.MessageB{})
					c.ProducesEventType(fixtures.MessageE{})
				},
			}),
		),
		Entry(
			"identity key differs",
			FromAggregate(&fixtures.AggregateMessageHandler{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("<name>", "799239e7-8c03-48f9-a324-14b7f9b76e30") // diff
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ConsumesCommandType(fixtures.MessageB{})
					c.ProducesEventType(fixtures.MessageE{})
				},
			}),
		),
		Entry(
			"messages differ",
			FromAggregate(&fixtures.AggregateMessageHandler{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("<name>", aggregateKey)
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageB{}) // diff
					c.ProducesEventType(fixtures.MessageE{})
				},
			}),
		),
	)
})
