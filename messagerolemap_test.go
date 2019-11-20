package configkit_test

import (
	. "github.com/dogmatiq/configkit"
	configfixtures "github.com/dogmatiq/configkit/fixtures"
	"github.com/dogmatiq/dogma/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ MessageTypeContainer = MessageRoleMap{}

var _ = Describe("type MessageRoleMap", func() {
	Describe("func Has()", func() {
		rm := MessageRoleMap{
			configfixtures.MessageAType: CommandMessageRole,
			configfixtures.MessageBType: EventMessageRole,
		}

		It("returns true if the type is in the map", func() {
			Expect(
				rm.Has(configfixtures.MessageAType),
			).To(BeTrue())
		})

		It("returns false if the type is not in the map", func() {
			Expect(
				rm.Has(configfixtures.MessageCType),
			).To(BeFalse())
		})
	})

	Describe("func HasM()", func() {
		rm := MessageRoleMap{
			configfixtures.MessageAType: CommandMessageRole,
			configfixtures.MessageBType: EventMessageRole,
		}

		It("returns true if the type is in the map", func() {
			Expect(
				rm.HasM(fixtures.MessageA1),
			).To(BeTrue())
		})

		It("returns false if the type is not in the map", func() {
			Expect(
				rm.HasM(fixtures.MessageC1),
			).To(BeFalse())
		})
	})

	Describe("func Add()", func() {
		It("adds the type to the map", func() {
			rm := MessageRoleMap{}
			rm.Add(configfixtures.MessageAType, CommandMessageRole)

			Expect(
				rm.Has(configfixtures.MessageAType),
			).To(BeTrue())
		})

		It("returns true if the type is not already in the map", func() {
			rm := MessageRoleMap{}

			Expect(
				rm.Add(configfixtures.MessageAType, CommandMessageRole),
			).To(BeTrue())
		})

		It("returns false if the type is already in the map", func() {
			rm := MessageRoleMap{}
			rm.Add(configfixtures.MessageAType, CommandMessageRole)

			Expect(
				rm.Add(configfixtures.MessageAType, EventMessageRole),
			).To(BeFalse())

			Expect(
				rm[configfixtures.MessageAType],
			).To(Equal(CommandMessageRole))
		})
	})

	Describe("func AddM()", func() {
		It("adds the type of the message to the map", func() {
			rm := MessageRoleMap{}
			rm.AddM(fixtures.MessageA1, CommandMessageRole)

			Expect(
				rm.Has(configfixtures.MessageAType),
			).To(BeTrue())
		})

		It("returns true if the type is not already in the map", func() {
			rm := MessageRoleMap{}

			Expect(
				rm.AddM(fixtures.MessageA1, CommandMessageRole),
			).To(BeTrue())
		})

		It("returns false if the type is already in the map", func() {
			rm := MessageRoleMap{}
			rm.AddM(fixtures.MessageA1, CommandMessageRole)

			Expect(
				rm.AddM(fixtures.MessageA1, EventMessageRole),
			).To(BeFalse())

			Expect(
				rm[configfixtures.MessageAType],
			).To(Equal(CommandMessageRole))
		})
	})

	Describe("func Remove()", func() {
		It("removes the type from the set", func() {
			rm := MessageRoleMap{configfixtures.MessageAType: CommandMessageRole}
			rm.Remove(configfixtures.MessageAType)

			Expect(
				rm.Has(configfixtures.MessageAType),
			).To(BeFalse())
		})

		It("returns true if the type is already in the set", func() {
			rm := MessageRoleMap{configfixtures.MessageAType: CommandMessageRole}

			Expect(
				rm.Remove(configfixtures.MessageAType),
			).To(BeTrue())
		})

		It("returns false if the type is not already in the set", func() {
			rm := MessageRoleMap{}

			Expect(
				rm.Remove(configfixtures.MessageAType),
			).To(BeFalse())
		})
	})

	Describe("func RemoveM()", func() {
		It("removes the type of the message from the set", func() {
			rm := MessageRoleMap{configfixtures.MessageAType: CommandMessageRole}
			rm.RemoveM(fixtures.MessageA1)

			Expect(
				rm.Has(configfixtures.MessageAType),
			).To(BeFalse())
		})

		It("returns true if the type is already in the set", func() {
			rm := MessageRoleMap{configfixtures.MessageAType: CommandMessageRole}

			Expect(
				rm.RemoveM(fixtures.MessageA1),
			).To(BeTrue())
		})

		It("returns false if the type is not already in the set", func() {
			rm := MessageRoleMap{}

			Expect(
				rm.RemoveM(fixtures.MessageA1),
			).To(BeFalse())
		})
	})

	Describe("func Each()", func() {
		rm := MessageRoleMap{
			configfixtures.MessageAType: CommandMessageRole,
			configfixtures.MessageBType: EventMessageRole,
		}

		It("calls fn for each type in the container", func() {
			var types []MessageType

			all := rm.Each(func(t MessageType) bool {
				types = append(types, t)
				return true
			})

			Expect(types).To(ConsistOf(configfixtures.MessageAType, configfixtures.MessageBType))
			Expect(all).To(BeTrue())
		})

		It("stops iterating if fn returns false", func() {
			count := 0

			all := rm.Each(func(t MessageType) bool {
				count++
				return false
			})

			Expect(count).To(BeNumerically("==", 1))
			Expect(all).To(BeFalse())
		})
	})
})
