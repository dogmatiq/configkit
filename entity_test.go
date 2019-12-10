package configkit_test

import (
	. "github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/fixtures" // can't dot-import due to conflicts
	"github.com/dogmatiq/configkit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("type EntityMessageNames", func() {
	Describe("func IsEqual()", func() {
		It("returns true if the sets are equivalent", func() {
			a := EntityMessageNames{
				Roles: message.NameRoles{
					fixtures.MessageATypeName: message.CommandRole,
					fixtures.MessageBTypeName: message.EventRole,
				},
				Produced: message.NameRoles{
					fixtures.MessageBTypeName: message.EventRole,
				},
				Consumed: message.NameRoles{
					fixtures.MessageATypeName: message.CommandRole,
				},
			}

			b := EntityMessageNames{
				Roles: message.NameRoles{
					fixtures.MessageATypeName: message.CommandRole,
					fixtures.MessageBTypeName: message.EventRole,
				},
				Produced: message.NameRoles{
					fixtures.MessageBTypeName: message.EventRole,
				},
				Consumed: message.NameRoles{
					fixtures.MessageATypeName: message.CommandRole,
				},
			}

			Expect(a.IsEqual(b)).To(BeTrue())
		})

		DescribeTable(
			"returns false if the sets are not equivalent",
			func(b EntityMessageNames) {
				a := EntityMessageNames{
					Roles: message.NameRoles{
						fixtures.MessageATypeName: message.CommandRole,
						fixtures.MessageBTypeName: message.EventRole,
					},
					Produced: message.NameRoles{
						fixtures.MessageBTypeName: message.EventRole,
					},
					Consumed: message.NameRoles{
						fixtures.MessageATypeName: message.CommandRole,
					},
				}
				Expect(a.IsEqual(b)).To(BeFalse())
			},
			Entry(
				"roles differ",
				EntityMessageNames{
					Roles: message.NameRoles{
						fixtures.MessageATypeName: message.CommandRole,
						fixtures.MessageBTypeName: message.EventRole,
						fixtures.MessageCTypeName: message.TimeoutRole, // diff
					},
					Produced: message.NameRoles{
						fixtures.MessageBTypeName: message.EventRole,
					},
					Consumed: message.NameRoles{
						fixtures.MessageATypeName: message.CommandRole,
					},
				},
			),
			Entry(
				"produced messages differ",
				EntityMessageNames{
					Roles: message.NameRoles{
						fixtures.MessageATypeName: message.CommandRole,
						fixtures.MessageBTypeName: message.EventRole,
					},
					Produced: message.NameRoles{
						fixtures.MessageBTypeName: message.EventRole,
						fixtures.MessageCTypeName: message.TimeoutRole, // diff
					},
					Consumed: message.NameRoles{
						fixtures.MessageATypeName: message.CommandRole,
					},
				},
			),
			Entry(
				"consumed messages differ",
				EntityMessageNames{
					Roles: message.NameRoles{
						fixtures.MessageATypeName: message.CommandRole,
						fixtures.MessageBTypeName: message.EventRole,
					},
					Produced: message.NameRoles{
						fixtures.MessageBTypeName: message.EventRole,
					},
					Consumed: message.NameRoles{
						fixtures.MessageATypeName: message.CommandRole,
						fixtures.MessageCTypeName: message.TimeoutRole, // diff
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
					fixtures.MessageAType: message.CommandRole,
					fixtures.MessageBType: message.EventRole,
				},
				Produced: message.TypeRoles{
					fixtures.MessageBType: message.EventRole,
				},
				Consumed: message.TypeRoles{
					fixtures.MessageAType: message.CommandRole,
				},
			}

			b := EntityMessageTypes{
				Roles: message.TypeRoles{
					fixtures.MessageAType: message.CommandRole,
					fixtures.MessageBType: message.EventRole,
				},
				Produced: message.TypeRoles{
					fixtures.MessageBType: message.EventRole,
				},
				Consumed: message.TypeRoles{
					fixtures.MessageAType: message.CommandRole,
				},
			}

			Expect(a.IsEqual(b)).To(BeTrue())
		})

		DescribeTable(
			"returns false if the sets are not equivalent",
			func(b EntityMessageTypes) {
				a := EntityMessageTypes{
					Roles: message.TypeRoles{
						fixtures.MessageAType: message.CommandRole,
						fixtures.MessageBType: message.EventRole,
					},
					Produced: message.TypeRoles{
						fixtures.MessageBType: message.EventRole,
					},
					Consumed: message.TypeRoles{
						fixtures.MessageAType: message.CommandRole,
					},
				}
				Expect(a.IsEqual(b)).To(BeFalse())
			},
			Entry(
				"roles differ",
				EntityMessageTypes{
					Roles: message.TypeRoles{
						fixtures.MessageAType: message.CommandRole,
						fixtures.MessageBType: message.EventRole,
						fixtures.MessageCType: message.TimeoutRole, // diff
					},
					Produced: message.TypeRoles{
						fixtures.MessageBType: message.EventRole,
					},
					Consumed: message.TypeRoles{
						fixtures.MessageAType: message.CommandRole,
					},
				},
			),
			Entry(
				"produced messages differ",
				EntityMessageTypes{
					Roles: message.TypeRoles{
						fixtures.MessageAType: message.CommandRole,
						fixtures.MessageBType: message.EventRole,
					},
					Produced: message.TypeRoles{
						fixtures.MessageBType: message.EventRole,
						fixtures.MessageCType: message.TimeoutRole, // diff
					},
					Consumed: message.TypeRoles{
						fixtures.MessageAType: message.CommandRole,
					},
				},
			),
			Entry(
				"consumed messages differ",
				EntityMessageTypes{
					Roles: message.TypeRoles{
						fixtures.MessageAType: message.CommandRole,
						fixtures.MessageBType: message.EventRole,
					},
					Produced: message.TypeRoles{
						fixtures.MessageBType: message.EventRole,
					},
					Consumed: message.TypeRoles{
						fixtures.MessageAType: message.CommandRole,
						fixtures.MessageCType: message.TimeoutRole, // diff
					},
				},
			),
		)
	})
})
