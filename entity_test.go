package configkit_test

import (
	. "github.com/dogmatiq/configkit"
	cfixtures "github.com/dogmatiq/configkit/fixtures" // can't dot-import due to conflicts
	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures" // can't dot-import due to conflicts
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("type EntityMessageNames", func() {
	Describe("func IsEqual()", func() {
		It("returns true if the sets are equivalent", func() {
			a := EntityMessageNames{
				Roles: message.NameRoles{
					cfixtures.MessageATypeName: message.CommandRole,
					cfixtures.MessageBTypeName: message.EventRole,
				},
				Produced: message.NameRoles{
					cfixtures.MessageBTypeName: message.EventRole,
				},
				Consumed: message.NameRoles{
					cfixtures.MessageATypeName: message.CommandRole,
				},
			}

			b := EntityMessageNames{
				Roles: message.NameRoles{
					cfixtures.MessageATypeName: message.CommandRole,
					cfixtures.MessageBTypeName: message.EventRole,
				},
				Produced: message.NameRoles{
					cfixtures.MessageBTypeName: message.EventRole,
				},
				Consumed: message.NameRoles{
					cfixtures.MessageATypeName: message.CommandRole,
				},
			}

			Expect(a.IsEqual(b)).To(BeTrue())
		})

		DescribeTable(
			"returns false if the sets are not equivalent",
			func(b EntityMessageNames) {
				a := EntityMessageNames{
					Roles: message.NameRoles{
						cfixtures.MessageATypeName: message.CommandRole,
						cfixtures.MessageBTypeName: message.EventRole,
					},
					Produced: message.NameRoles{
						cfixtures.MessageBTypeName: message.EventRole,
					},
					Consumed: message.NameRoles{
						cfixtures.MessageATypeName: message.CommandRole,
					},
				}
				Expect(a.IsEqual(b)).To(BeFalse())
			},
			Entry(
				"roles differ",
				EntityMessageNames{
					Roles: message.NameRoles{
						cfixtures.MessageATypeName: message.CommandRole,
						cfixtures.MessageBTypeName: message.EventRole,
						cfixtures.MessageCTypeName: message.TimeoutRole, // diff
					},
					Produced: message.NameRoles{
						cfixtures.MessageBTypeName: message.EventRole,
					},
					Consumed: message.NameRoles{
						cfixtures.MessageATypeName: message.CommandRole,
					},
				},
			),
			Entry(
				"produced messages differ",
				EntityMessageNames{
					Roles: message.NameRoles{
						cfixtures.MessageATypeName: message.CommandRole,
						cfixtures.MessageBTypeName: message.EventRole,
					},
					Produced: message.NameRoles{
						cfixtures.MessageBTypeName: message.EventRole,
						cfixtures.MessageCTypeName: message.TimeoutRole, // diff
					},
					Consumed: message.NameRoles{
						cfixtures.MessageATypeName: message.CommandRole,
					},
				},
			),
			Entry(
				"consumed messages differ",
				EntityMessageNames{
					Roles: message.NameRoles{
						cfixtures.MessageATypeName: message.CommandRole,
						cfixtures.MessageBTypeName: message.EventRole,
					},
					Produced: message.NameRoles{
						cfixtures.MessageBTypeName: message.EventRole,
					},
					Consumed: message.NameRoles{
						cfixtures.MessageATypeName: message.CommandRole,
						cfixtures.MessageCTypeName: message.TimeoutRole, // diff
					},
				},
			),
		)
	})
})

var _ = Describe("type EntityMessageTypes", func() {
	Describe("func IsEqual()", func() {
		It("returns true if the sets are equivalent", func() {
			a := EntityMessageTypes{
				Roles: message.TypeRoles{
					cfixtures.MessageAType: message.CommandRole,
					cfixtures.MessageBType: message.EventRole,
				},
				Produced: message.TypeRoles{
					cfixtures.MessageBType: message.EventRole,
				},
				Consumed: message.TypeRoles{
					cfixtures.MessageAType: message.CommandRole,
				},
			}

			b := EntityMessageTypes{
				Roles: message.TypeRoles{
					cfixtures.MessageAType: message.CommandRole,
					cfixtures.MessageBType: message.EventRole,
				},
				Produced: message.TypeRoles{
					cfixtures.MessageBType: message.EventRole,
				},
				Consumed: message.TypeRoles{
					cfixtures.MessageAType: message.CommandRole,
				},
			}

			Expect(a.IsEqual(b)).To(BeTrue())
		})

		DescribeTable(
			"returns false if the sets are not equivalent",
			func(b EntityMessageTypes) {
				a := EntityMessageTypes{
					Roles: message.TypeRoles{
						cfixtures.MessageAType: message.CommandRole,
						cfixtures.MessageBType: message.EventRole,
					},
					Produced: message.TypeRoles{
						cfixtures.MessageBType: message.EventRole,
					},
					Consumed: message.TypeRoles{
						cfixtures.MessageAType: message.CommandRole,
					},
				}
				Expect(a.IsEqual(b)).To(BeFalse())
			},
			Entry(
				"roles differ",
				EntityMessageTypes{
					Roles: message.TypeRoles{
						cfixtures.MessageAType: message.CommandRole,
						cfixtures.MessageBType: message.EventRole,
						cfixtures.MessageCType: message.TimeoutRole, // diff
					},
					Produced: message.TypeRoles{
						cfixtures.MessageBType: message.EventRole,
					},
					Consumed: message.TypeRoles{
						cfixtures.MessageAType: message.CommandRole,
					},
				},
			),
			Entry(
				"produced messages differ",
				EntityMessageTypes{
					Roles: message.TypeRoles{
						cfixtures.MessageAType: message.CommandRole,
						cfixtures.MessageBType: message.EventRole,
					},
					Produced: message.TypeRoles{
						cfixtures.MessageBType: message.EventRole,
						cfixtures.MessageCType: message.TimeoutRole, // diff
					},
					Consumed: message.TypeRoles{
						cfixtures.MessageAType: message.CommandRole,
					},
				},
			),
			Entry(
				"consumed messages differ",
				EntityMessageTypes{
					Roles: message.TypeRoles{
						cfixtures.MessageAType: message.CommandRole,
						cfixtures.MessageBType: message.EventRole,
					},
					Produced: message.TypeRoles{
						cfixtures.MessageBType: message.EventRole,
					},
					Consumed: message.TypeRoles{
						cfixtures.MessageAType: message.CommandRole,
						cfixtures.MessageCType: message.TimeoutRole, // diff
					},
				},
			),
		)
	})
})

var _ = Describe("func IsEqual()", func() {
	It("returns true if the two entities are equivalent", func() {
		h := &fixtures.AggregateMessageHandler{
			ConfigureFunc: func(c dogma.AggregateConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ConsumesCommandType(fixtures.MessageB{})
				c.ProducesEventType(fixtures.MessageE{})
			},
		}

		a := FromAggregate(h)
		b := FromAggregate(h)

		Expect(IsEqual(a, b)).To(BeTrue())
	})

	DescribeTable(
		"returns false if the entities are not equivalent",
		func(b Entity) {
			h := &fixtures.AggregateMessageHandler{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("<name>", "<key>")
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ConsumesCommandType(fixtures.MessageB{})
					c.ProducesEventType(fixtures.MessageE{})
				},
			}

			a := FromAggregate(h)

			Expect(IsEqual(a, b)).To(BeFalse())
		},
		Entry(
			"type differs",
			FromIntegration(&fixtures.IntegrationMessageHandler{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("<name>", "<key>")
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
					c.Identity("<name-different>", "<key>") // diff
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
					c.Identity("<name>", "<key-different>") // diff
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
					c.Identity("<name>", "<key>")
					c.ConsumesCommandType(fixtures.MessageA{})
					c.ProducesEventType(fixtures.MessageB{}) // diff
					c.ProducesEventType(fixtures.MessageE{})
				},
			}),
		),
	)
})
