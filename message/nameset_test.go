package message_test

import (
	. "github.com/dogmatiq/configkit/fixtures"
	. "github.com/dogmatiq/configkit/message"
	. "github.com/dogmatiq/dogma/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ NameCollection = NameSet{}

var _ = Describe("type NameSet", func() {
	Describe("func NewNameSet()", func() {
		It("returns a set containing the given names", func() {
			Expect(NewNameSet(
				MessageATypeName,
				MessageBTypeName,
			)).To(Equal(NameSet{
				MessageATypeName: struct{}{},
				MessageBTypeName: struct{}{},
			}))
		})
	})

	Describe("func NamesOf()", func() {
		It("returns a set containing the names of the given messages", func() {
			Expect(NamesOf(
				MessageA1,
				MessageB1,
			)).To(Equal(NameSet{
				MessageATypeName: struct{}{},
				MessageBTypeName: struct{}{},
			}))
		})
	})

	Describe("func Has()", func() {
		set := NamesOf(
			MessageA1,
			MessageB1,
		)

		It("returns true if the name is in the set", func() {
			Expect(
				set.Has(MessageATypeName),
			).To(BeTrue())
		})

		It("returns false if the name is not in the set", func() {
			Expect(
				set.Has(MessageCTypeName),
			).To(BeFalse())
		})
	})

	Describe("func HasM()", func() {
		set := NamesOf(
			MessageA1,
			MessageB1,
		)

		It("returns true if the name is in the set", func() {
			Expect(
				set.HasM(MessageA1),
			).To(BeTrue())
		})

		It("returns false if the name is not in the set", func() {
			Expect(
				set.HasM(MessageC1),
			).To(BeFalse())
		})
	})

	Describe(("func Add()"), func() {
		It("adds the name to the set", func() {
			s := NamesOf()
			s.Add(MessageATypeName)

			Expect(
				s.Has(MessageATypeName),
			).To(BeTrue())
		})

		It("returns true if the name is not already in the set", func() {
			s := NamesOf()

			Expect(
				s.Add(MessageATypeName),
			).To(BeTrue())
		})

		It("returns false if the name is already in the set", func() {
			s := NamesOf()
			s.Add(MessageATypeName)

			Expect(
				s.Add(MessageATypeName),
			).To(BeFalse())
		})
	})

	Describe("func AddM()", func() {
		It("adds the name of the message to the set", func() {
			s := NamesOf()
			s.AddM(MessageA1)

			Expect(
				s.Has(MessageATypeName),
			).To(BeTrue())
		})

		It("returns true if the name is not already in the set", func() {
			s := NamesOf()

			Expect(
				s.AddM(MessageA1),
			).To(BeTrue())
		})

		It("returns false if the name is already in the set", func() {
			s := NamesOf()
			s.Add(MessageATypeName)

			Expect(
				s.AddM(MessageA1),
			).To(BeFalse())
		})
	})

	Describe("func Remove()", func() {
		It("removes the name from the set", func() {
			s := NamesOf(MessageA1)
			s.Remove(MessageATypeName)

			Expect(
				s.Has(MessageATypeName),
			).To(BeFalse())
		})

		It("returns true if the name is already in the set", func() {
			s := NamesOf()
			s.Add(MessageATypeName)

			Expect(
				s.Remove(MessageATypeName),
			).To(BeTrue())
		})

		It("returns false if the name is not already in the set", func() {
			s := NamesOf()

			Expect(
				s.Remove(MessageATypeName),
			).To(BeFalse())
		})
	})

	Describe("func RemoveM()", func() {
		It("removes the name of the message from the set", func() {
			s := NamesOf(MessageA1)
			s.RemoveM(MessageA1)

			Expect(
				s.Has(MessageATypeName),
			).To(BeFalse())
		})

		It("returns true if the name is already in the set", func() {
			s := NamesOf()
			s.Add(MessageATypeName)

			Expect(
				s.RemoveM(MessageA1),
			).To(BeTrue())
		})

		It("returns false if the name is not already in the set", func() {
			s := NamesOf()

			Expect(
				s.RemoveM(MessageA1),
			).To(BeFalse())
		})
	})

	Describe("func Each()", func() {
		s := NewNameSet(
			MessageATypeName,
			MessageBTypeName,
		)

		It("calls fn for each name in the container", func() {
			var names []Name

			all := s.Each(func(n Name) bool {
				names = append(names, n)
				return true
			})

			Expect(names).To(ConsistOf(MessageATypeName, MessageBTypeName))
			Expect(all).To(BeTrue())
		})

		It("stops iterating if fn returns false", func() {
			count := 0

			all := s.Each(func(n Name) bool {
				count++
				return false
			})

			Expect(count).To(BeNumerically("==", 1))
			Expect(all).To(BeFalse())
		})
	})
})
