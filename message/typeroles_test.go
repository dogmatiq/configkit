package message_test

import (
	. "github.com/dogmatiq/configkit/fixtures"
	. "github.com/dogmatiq/configkit/message"
	. "github.com/dogmatiq/dogma/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ TypeCollection = TypeRoles{}

var _ = Describe("type TypeRoles", func() {
	Describe("func Has()", func() {
		tr := TypeRoles{
			MessageAType: CommandRole,
			MessageBType: EventRole,
		}

		It("returns true if the type is in the map", func() {
			Expect(
				tr.Has(MessageAType),
			).To(BeTrue())
		})

		It("returns false if the type is not in the map", func() {
			Expect(
				tr.Has(MessageCType),
			).To(BeFalse())
		})
	})

	Describe("func HasM()", func() {
		tr := TypeRoles{
			MessageAType: CommandRole,
			MessageBType: EventRole,
		}

		It("returns true if the type is in the map", func() {
			Expect(
				tr.HasM(MessageA1),
			).To(BeTrue())
		})

		It("returns false if the type is not in the map", func() {
			Expect(
				tr.HasM(MessageC1),
			).To(BeFalse())
		})
	})

	Describe("func Add()", func() {
		It("adds the type to the map", func() {
			tr := TypeRoles{}
			tr.Add(MessageAType, CommandRole)

			Expect(
				tr.Has(MessageAType),
			).To(BeTrue())
		})

		It("returns true if the type is not already in the map", func() {
			tr := TypeRoles{}

			Expect(
				tr.Add(MessageAType, CommandRole),
			).To(BeTrue())
		})

		It("returns false if the type is already in the map", func() {
			tr := TypeRoles{}
			tr.Add(MessageAType, CommandRole)

			Expect(
				tr.Add(MessageAType, EventRole),
			).To(BeFalse())

			Expect(
				tr[MessageAType],
			).To(Equal(CommandRole))
		})
	})

	Describe("func AddM()", func() {
		It("adds the type of the message to the map", func() {
			tr := TypeRoles{}
			tr.AddM(MessageA1, CommandRole)

			Expect(
				tr.Has(MessageAType),
			).To(BeTrue())
		})

		It("returns true if the type is not already in the map", func() {
			tr := TypeRoles{}

			Expect(
				tr.AddM(MessageA1, CommandRole),
			).To(BeTrue())
		})

		It("returns false if the type is already in the map", func() {
			tr := TypeRoles{}
			tr.AddM(MessageA1, CommandRole)

			Expect(
				tr.AddM(MessageA1, EventRole),
			).To(BeFalse())

			Expect(
				tr[MessageAType],
			).To(Equal(CommandRole))
		})
	})

	Describe("func Remove()", func() {
		It("removes the type from the set", func() {
			tr := TypeRoles{MessageAType: CommandRole}
			tr.Remove(MessageAType)

			Expect(
				tr.Has(MessageAType),
			).To(BeFalse())
		})

		It("returns true if the type is already in the set", func() {
			tr := TypeRoles{MessageAType: CommandRole}

			Expect(
				tr.Remove(MessageAType),
			).To(BeTrue())
		})

		It("returns false if the type is not already in the set", func() {
			tr := TypeRoles{}

			Expect(
				tr.Remove(MessageAType),
			).To(BeFalse())
		})
	})

	Describe("func RemoveM()", func() {
		It("removes the type of the message from the set", func() {
			tr := TypeRoles{MessageAType: CommandRole}
			tr.RemoveM(MessageA1)

			Expect(
				tr.Has(MessageAType),
			).To(BeFalse())
		})

		It("returns true if the type is already in the set", func() {
			tr := TypeRoles{MessageAType: CommandRole}

			Expect(
				tr.RemoveM(MessageA1),
			).To(BeTrue())
		})

		It("returns false if the type is not already in the set", func() {
			tr := TypeRoles{}

			Expect(
				tr.RemoveM(MessageA1),
			).To(BeFalse())
		})
	})

	Describe("func Each()", func() {
		tr := TypeRoles{
			MessageAType: CommandRole,
			MessageBType: EventRole,
		}

		It("calls fn for each type in the container", func() {
			var types []Type

			all := tr.Each(func(t Type) bool {
				types = append(types, t)
				return true
			})

			Expect(types).To(ConsistOf(MessageAType, MessageBType))
			Expect(all).To(BeTrue())
		})

		It("stops iterating if fn returns false", func() {
			count := 0

			all := tr.Each(func(t Type) bool {
				count++
				return false
			})

			Expect(count).To(BeNumerically("==", 1))
			Expect(all).To(BeFalse())
		})
	})
})
