package message_test

import (
	. "github.com/dogmatiq/configkit/fixtures"
	. "github.com/dogmatiq/configkit/message"
	. "github.com/dogmatiq/dogma/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ TypeCollection = TypeSet{}

var _ = Describe("type TypeSet", func() {
	Describe("func NewTypeSet()", func() {
		It("returns a set containing the given types", func() {
			Expect(NewTypeSet(
				MessageAType,
				MessageBType,
			)).To(Equal(TypeSet{
				MessageAType: struct{}{},
				MessageBType: struct{}{},
			}))
		})
	})

	Describe("func TypesOf()", func() {
		It("returns a set containing the types of the given messages", func() {
			Expect(TypesOf(
				MessageA1,
				MessageB1,
			)).To(Equal(TypeSet{
				MessageAType: struct{}{},
				MessageBType: struct{}{},
			}))
		})
	})

	Describe("func Has()", func() {
		set := TypesOf(
			MessageA1,
			MessageB1,
		)

		It("returns true if the type is in the set", func() {
			Expect(
				set.Has(MessageAType),
			).To(BeTrue())
		})

		It("returns false if the type is not in the set", func() {
			Expect(
				set.Has(MessageCType),
			).To(BeFalse())
		})
	})

	Describe("func HasM()", func() {
		set := TypesOf(
			MessageA1,
			MessageB1,
		)

		It("returns true if the type is in the set", func() {
			Expect(
				set.HasM(MessageA1),
			).To(BeTrue())
		})

		It("returns false if the type is not in the set", func() {
			Expect(
				set.HasM(MessageC1),
			).To(BeFalse())
		})
	})

	Describe(("func Add()"), func() {
		It("adds the type to the set", func() {
			s := TypesOf()
			s.Add(MessageAType)

			Expect(
				s.Has(MessageAType),
			).To(BeTrue())
		})

		It("returns true if the type is not already in the set", func() {
			s := TypesOf()

			Expect(
				s.Add(MessageAType),
			).To(BeTrue())
		})

		It("returns false if the type is already in the set", func() {
			s := TypesOf()
			s.Add(MessageAType)

			Expect(
				s.Add(MessageAType),
			).To(BeFalse())
		})
	})

	Describe("func AddM()", func() {
		It("adds the type of the message to the set", func() {
			s := TypesOf()
			s.AddM(MessageA1)

			Expect(
				s.Has(MessageAType),
			).To(BeTrue())
		})

		It("returns true if the type is not already in the set", func() {
			s := TypesOf()

			Expect(
				s.AddM(MessageA1),
			).To(BeTrue())
		})

		It("returns false if the type is already in the set", func() {
			s := TypesOf()
			s.Add(MessageAType)

			Expect(
				s.AddM(MessageA1),
			).To(BeFalse())
		})
	})

	Describe("func Remove()", func() {
		It("removes the type from the set", func() {
			s := TypesOf(MessageA1)
			s.Remove(MessageAType)

			Expect(
				s.Has(MessageAType),
			).To(BeFalse())
		})

		It("returns true if the type is already in the set", func() {
			s := TypesOf()
			s.Add(MessageAType)

			Expect(
				s.Remove(MessageAType),
			).To(BeTrue())
		})

		It("returns false if the type is not already in the set", func() {
			s := TypesOf()

			Expect(
				s.Remove(MessageAType),
			).To(BeFalse())
		})
	})

	Describe("func RemoveM()", func() {
		It("removes the type of the message from the set", func() {
			s := TypesOf(MessageA1)
			s.RemoveM(MessageA1)

			Expect(
				s.Has(MessageAType),
			).To(BeFalse())
		})

		It("returns true if the type is already in the set", func() {
			s := TypesOf()
			s.Add(MessageAType)

			Expect(
				s.RemoveM(MessageA1),
			).To(BeTrue())
		})

		It("returns false if the type is not already in the set", func() {
			s := TypesOf()

			Expect(
				s.RemoveM(MessageA1),
			).To(BeFalse())
		})
	})

	Describe("func Each()", func() {
		s := NewTypeSet(
			MessageAType,
			MessageBType,
		)

		It("calls fn for each type in the container", func() {
			var types []Type

			all := s.Each(func(t Type) bool {
				types = append(types, t)
				return true
			})

			Expect(types).To(ConsistOf(MessageAType, MessageBType))
			Expect(all).To(BeTrue())
		})

		It("stops iterating if fn returns false", func() {
			count := 0

			all := s.Each(func(t Type) bool {
				count++
				return false
			})

			Expect(count).To(BeNumerically("==", 1))
			Expect(all).To(BeFalse())
		})
	})
})
