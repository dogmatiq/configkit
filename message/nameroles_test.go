package message_test

import (
	. "github.com/dogmatiq/configkit/fixtures"
	. "github.com/dogmatiq/configkit/message"
	. "github.com/dogmatiq/dogma/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ NameCollection = NameRoles{}

var _ = Describe("type NameRoles", func() {
	Describe("func Has()", func() {
		nr := NameRoles{
			MessageATypeName: CommandRole,
			MessageBTypeName: EventRole,
		}

		It("returns true if the name is in the map", func() {
			Expect(
				nr.Has(MessageATypeName),
			).To(BeTrue())
		})

		It("returns false if the name is not in the map", func() {
			Expect(
				nr.Has(MessageCTypeName),
			).To(BeFalse())
		})
	})

	Describe("func HasM()", func() {
		nr := NameRoles{
			MessageATypeName: CommandRole,
			MessageBTypeName: EventRole,
		}

		It("returns true if the name is in the map", func() {
			Expect(
				nr.HasM(MessageA1),
			).To(BeTrue())
		})

		It("returns false if the name is not in the map", func() {
			Expect(
				nr.HasM(MessageC1),
			).To(BeFalse())
		})
	})

	Describe("func Add()", func() {
		It("adds the name to the map", func() {
			nr := NameRoles{}
			nr.Add(MessageATypeName, CommandRole)

			Expect(
				nr.Has(MessageATypeName),
			).To(BeTrue())
		})

		It("returns true if the name is not already in the map", func() {
			nr := NameRoles{}

			Expect(
				nr.Add(MessageATypeName, CommandRole),
			).To(BeTrue())
		})

		It("returns false if the name is already in the map", func() {
			nr := NameRoles{}
			nr.Add(MessageATypeName, CommandRole)

			Expect(
				nr.Add(MessageATypeName, EventRole),
			).To(BeFalse())

			Expect(
				nr[MessageATypeName],
			).To(Equal(CommandRole))
		})
	})

	Describe("func AddM()", func() {
		It("adds the name of the message to the map", func() {
			nr := NameRoles{}
			nr.AddM(MessageA1, CommandRole)

			Expect(
				nr.Has(MessageATypeName),
			).To(BeTrue())
		})

		It("returns true if the name is not already in the map", func() {
			nr := NameRoles{}

			Expect(
				nr.AddM(MessageA1, CommandRole),
			).To(BeTrue())
		})

		It("returns false if the name is already in the map", func() {
			nr := NameRoles{}
			nr.AddM(MessageA1, CommandRole)

			Expect(
				nr.AddM(MessageA1, EventRole),
			).To(BeFalse())

			Expect(
				nr[MessageATypeName],
			).To(Equal(CommandRole))
		})
	})

	Describe("func Remove()", func() {
		It("removes the name from the set", func() {
			nr := NameRoles{MessageATypeName: CommandRole}
			nr.Remove(MessageATypeName)

			Expect(
				nr.Has(MessageATypeName),
			).To(BeFalse())
		})

		It("returns true if the name is already in the set", func() {
			nr := NameRoles{MessageATypeName: CommandRole}

			Expect(
				nr.Remove(MessageATypeName),
			).To(BeTrue())
		})

		It("returns false if the name is not already in the set", func() {
			nr := NameRoles{}

			Expect(
				nr.Remove(MessageATypeName),
			).To(BeFalse())
		})
	})

	Describe("func RemoveM()", func() {
		It("removes the name of the message from the set", func() {
			nr := NameRoles{MessageATypeName: CommandRole}
			nr.RemoveM(MessageA1)

			Expect(
				nr.Has(MessageATypeName),
			).To(BeFalse())
		})

		It("returns true if the name is already in the set", func() {
			nr := NameRoles{MessageATypeName: CommandRole}

			Expect(
				nr.RemoveM(MessageA1),
			).To(BeTrue())
		})

		It("returns false if the name is not already in the set", func() {
			nr := NameRoles{}

			Expect(
				nr.RemoveM(MessageA1),
			).To(BeFalse())
		})
	})

	Describe("func IsEqual()", func() {
		DescribeTable(
			"returns true if the sets are equivalent",
			func(a, b NameRoles) {
				Expect(a.IsEqual(b)).To(BeTrue())
			},
			Entry(
				"equivalent",
				NameRoles{
					MessageATypeName: CommandRole,
					MessageBTypeName: EventRole,
				},
				NameRoles{
					MessageATypeName: CommandRole,
					MessageBTypeName: EventRole,
				},
			),
			Entry(
				"nil and empty",
				NameRoles{},
				NameRoles(nil),
			),
		)

		DescribeTable(
			"returns false if the sets are not equivalent",
			func(b NameRoles) {
				a := NameRoles{
					MessageATypeName: CommandRole,
					MessageBTypeName: EventRole,
				}
				Expect(a.IsEqual(b)).To(BeFalse())
			},
			Entry(
				"subset",
				NameRoles{
					MessageATypeName: CommandRole,
				},
			),
			Entry(
				"superset",
				NameRoles{
					MessageATypeName: CommandRole,
					MessageBTypeName: EventRole,
					MessageCTypeName: TimeoutRole,
				},
			),
			Entry(
				"same-length, disjoint type",
				NameRoles{
					MessageATypeName: CommandRole,
					MessageCTypeName: EventRole,
				},
			),
			Entry(
				"same-length, disjoint role",
				NameRoles{
					MessageATypeName: CommandRole,
					MessageBTypeName: TimeoutRole,
				},
			),
		)
	})

	Describe("func Range()", func() {
		nr := NameRoles{
			MessageATypeName: CommandRole,
			MessageBTypeName: EventRole,
		}

		It("calls fn for each name in the container", func() {
			var names []Name

			all := nr.Range(func(n Name) bool {
				names = append(names, n)
				return true
			})

			Expect(names).To(ConsistOf(MessageATypeName, MessageBTypeName))
			Expect(all).To(BeTrue())
		})

		It("stops iterating if fn returns false", func() {
			count := 0

			all := nr.Range(func(n Name) bool {
				count++
				return false
			})

			Expect(count).To(BeNumerically("==", 1))
			Expect(all).To(BeFalse())
		})
	})

	Describe("func RangeByRole()", func() {
		nr := NameRoles{
			MessageATypeName: CommandRole,
			MessageBTypeName: EventRole,
			MessageCTypeName: CommandRole,
		}

		It("calls fn for each name in the container with the given role", func() {
			var names []Name

			all := nr.RangeByRole(
				CommandRole,
				func(n Name) bool {
					names = append(names, n)
					return true
				},
			)

			Expect(names).To(ConsistOf(MessageATypeName, MessageCTypeName))
			Expect(all).To(BeTrue())
		})

		It("stops iterating if fn returns false", func() {
			count := 0

			all := nr.RangeByRole(
				CommandRole,
				func(n Name) bool {
					count++
					return false
				},
			)

			Expect(count).To(BeNumerically("==", 1))
			Expect(all).To(BeFalse())
		})
	})
})
