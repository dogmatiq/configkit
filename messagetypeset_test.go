package configkit_test

import (
	. "github.com/dogmatiq/configkit"
	configfixtures "github.com/dogmatiq/configkit/fixtures"
	"github.com/dogmatiq/dogma/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ MessageTypeContainer = MessageTypeSet{}

var _ = Describe("type MessageTypeSet", func() {
	Describe("func NewMessageTypeSet", func() {
		It("returns a set containing the given types", func() {
			Expect(NewMessageTypeSet(
				configfixtures.MessageAType,
				configfixtures.MessageBType,
			)).To(Equal(MessageTypeSet{
				configfixtures.MessageAType: struct{}{},
				configfixtures.MessageBType: struct{}{},
			}))
		})
	})

	Describe("func MessageTypesOf", func() {
		It("returns a set containing the types of the given messages", func() {
			Expect(MessageTypesOf(
				fixtures.MessageA1,
				fixtures.MessageB1,
			)).To(Equal(MessageTypeSet{
				configfixtures.MessageAType: struct{}{},
				configfixtures.MessageBType: struct{}{},
			}))
		})
	})

	Describe("func Has", func() {
		set := MessageTypesOf(
			fixtures.MessageA1,
			fixtures.MessageB1,
		)

		It("returns true if the type is in the set", func() {
			Expect(
				set.Has(configfixtures.MessageAType),
			).To(BeTrue())
		})

		It("returns false if the type is not in the set", func() {
			Expect(
				set.Has(configfixtures.MessageCType),
			).To(BeFalse())
		})
	})

	Describe("func HasM", func() {
		set := MessageTypesOf(
			fixtures.MessageA1,
			fixtures.MessageB1,
		)

		It("returns true if the type is in the set", func() {
			Expect(
				set.HasM(fixtures.MessageA1),
			).To(BeTrue())
		})

		It("returns false if the type is not in the set", func() {
			Expect(
				set.HasM(fixtures.MessageC1),
			).To(BeFalse())
		})
	})

	Describe("func Add", func() {
		It("adds the type to the set", func() {
			s := MessageTypesOf()
			s.Add(configfixtures.MessageAType)

			Expect(
				s.Has(configfixtures.MessageAType),
			).To(BeTrue())
		})

		It("returns true if the type is not already in the set", func() {
			s := MessageTypesOf()

			Expect(
				s.Add(configfixtures.MessageAType),
			).To(BeTrue())
		})

		It("returns false if the type is already in the set", func() {
			s := MessageTypesOf()
			s.Add(configfixtures.MessageAType)

			Expect(
				s.Add(configfixtures.MessageAType),
			).To(BeFalse())
		})
	})

	Describe("func AddM", func() {
		It("adds the type of the message to the set", func() {
			s := MessageTypesOf()
			s.AddM(fixtures.MessageA1)

			Expect(
				s.Has(configfixtures.MessageAType),
			).To(BeTrue())
		})

		It("returns true if the type is not already in the set", func() {
			s := MessageTypesOf()

			Expect(
				s.AddM(fixtures.MessageA1),
			).To(BeTrue())
		})

		It("returns false if the type is already in the set", func() {
			s := MessageTypesOf()
			s.Add(configfixtures.MessageAType)

			Expect(
				s.AddM(fixtures.MessageA1),
			).To(BeFalse())
		})
	})

	Describe("func Remove", func() {
		It("removes the type from the set", func() {
			s := MessageTypesOf(fixtures.MessageA1)
			s.Remove(configfixtures.MessageAType)

			Expect(
				s.Has(configfixtures.MessageAType),
			).To(BeFalse())
		})

		It("returns true if the type is already in the set", func() {
			s := MessageTypesOf()
			s.Add(configfixtures.MessageAType)

			Expect(
				s.Remove(configfixtures.MessageAType),
			).To(BeTrue())
		})

		It("returns false if the type is not already in the set", func() {
			s := MessageTypesOf()

			Expect(
				s.Remove(configfixtures.MessageAType),
			).To(BeFalse())
		})
	})

	Describe("func RemoveM", func() {
		It("removes the type of the message from the set", func() {
			s := MessageTypesOf(fixtures.MessageA1)
			s.RemoveM(fixtures.MessageA1)

			Expect(
				s.Has(configfixtures.MessageAType),
			).To(BeFalse())
		})

		It("returns true if the type is already in the set", func() {
			s := MessageTypesOf()
			s.Add(configfixtures.MessageAType)

			Expect(
				s.RemoveM(fixtures.MessageA1),
			).To(BeTrue())
		})

		It("returns false if the type is not already in the set", func() {
			s := MessageTypesOf()

			Expect(
				s.RemoveM(fixtures.MessageA1),
			).To(BeFalse())
		})
	})

	Describe("func Each()", func() {
		s := NewMessageTypeSet(
			configfixtures.MessageAType,
			configfixtures.MessageBType,
		)

		It("calls fn for each type in the container", func() {
			var types []MessageType

			all := s.Each(func(t MessageType) bool {
				types = append(types, t)
				return true
			})

			Expect(types).To(ConsistOf(configfixtures.MessageAType, configfixtures.MessageBType))
			Expect(all).To(BeTrue())
		})

		It("stops iterating if fn returns false", func() {
			count := 0

			all := s.Each(func(t MessageType) bool {
				count++
				return false
			})

			Expect(count).To(BeNumerically("==", 1))
			Expect(all).To(BeFalse())
		})
	})
})
