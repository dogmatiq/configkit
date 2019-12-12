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
	Describe("func RoleOf()", func() {
		It("returns the role of a produced message", func() {
			m := EntityMessageNames{
				Produced: message.NameRoles{
					fixtures.MessageATypeName: message.CommandRole,
				},
			}

			r, ok := m.RoleOf(fixtures.MessageATypeName)
			Expect(ok).To(BeTrue())
			Expect(r).To(Equal(message.CommandRole))
		})

		It("returns the role of a consumed message", func() {
			m := EntityMessageNames{
				Consumed: message.NameRoles{
					fixtures.MessageATypeName: message.CommandRole,
				},
			}

			r, ok := m.RoleOf(fixtures.MessageATypeName)
			Expect(ok).To(BeTrue())
			Expect(r).To(Equal(message.CommandRole))
		})

		It("returns false is the message is neither produced or consumed", func() {
			m := EntityMessageNames{}

			_, ok := m.RoleOf(fixtures.MessageATypeName)
			Expect(ok).To(BeFalse())
		})
	})

	Describe("func All()", func() {
		It("returns the union of the produced and consumed messages", func() {
			m := EntityMessageNames{
				Produced: message.NameRoles{
					fixtures.MessageCTypeName: message.CommandRole,
				},
				Consumed: message.NameRoles{
					fixtures.MessageETypeName: message.EventRole,
				},
			}

			Expect(m.All()).To(Equal(
				message.NameRoles{
					fixtures.MessageCTypeName: message.CommandRole,
					fixtures.MessageETypeName: message.EventRole,
				},
			))
		})
	})

	Describe("func IsEqual()", func() {
		It("returns true if the sets are equivalent", func() {
			a := EntityMessageNames{
				Produced: message.NameRoles{
					fixtures.MessageBTypeName: message.EventRole,
				},
				Consumed: message.NameRoles{
					fixtures.MessageATypeName: message.CommandRole,
				},
			}

			b := EntityMessageNames{
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
				"produced messages differ",
				EntityMessageNames{
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
	Describe("func RoleOf()", func() {
		It("returns the role of a produced message", func() {
			m := EntityMessageTypes{
				Produced: message.TypeRoles{
					fixtures.MessageAType: message.CommandRole,
				},
			}

			r, ok := m.RoleOf(fixtures.MessageAType)
			Expect(ok).To(BeTrue())
			Expect(r).To(Equal(message.CommandRole))
		})

		It("returns the role of a consumed message", func() {
			m := EntityMessageTypes{
				Consumed: message.TypeRoles{
					fixtures.MessageAType: message.CommandRole,
				},
			}

			r, ok := m.RoleOf(fixtures.MessageAType)
			Expect(ok).To(BeTrue())
			Expect(r).To(Equal(message.CommandRole))
		})

		It("returns false is the message is neither produced or consumed", func() {
			m := EntityMessageTypes{}

			_, ok := m.RoleOf(fixtures.MessageAType)
			Expect(ok).To(BeFalse())
		})
	})

	Describe("func All()", func() {
		It("returns the union of the produced and consumed messages", func() {
			m := EntityMessageTypes{
				Produced: message.TypeRoles{
					fixtures.MessageCType: message.CommandRole,
				},
				Consumed: message.TypeRoles{
					fixtures.MessageEType: message.EventRole,
				},
			}

			Expect(m.All()).To(Equal(
				message.TypeRoles{
					fixtures.MessageCType: message.CommandRole,
					fixtures.MessageEType: message.EventRole,
				},
			))
		})
	})

	Describe("func IsEqual()", func() {
		It("returns true if the sets are equivalent", func() {
			a := EntityMessageTypes{
				Produced: message.TypeRoles{
					fixtures.MessageBType: message.EventRole,
				},
				Consumed: message.TypeRoles{
					fixtures.MessageAType: message.CommandRole,
				},
			}

			b := EntityMessageTypes{
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
				"produced messages differ",
				EntityMessageTypes{
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
