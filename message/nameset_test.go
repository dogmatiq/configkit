package message_test

import (
	. "github.com/dogmatiq/configkit/fixtures"
	. "github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma/fixtures"
	. "github.com/dogmatiq/dogma/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
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

	Describe("func IntersectionN()", func() {
		It("returns an empty set if no sets are given", func() {
			Expect(IntersectionN()).To(BeEmpty())
		})

		It("returns the original set if a single set is given", func() {
			a := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
			Expect(IntersectionN(a)).To(Equal(a))
		})

		It("returns the original set for identical sets", func() {
			a := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
			b := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
			c := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
			Expect(IntersectionN(a, b, c)).To(Equal(a))
		})

		It("returns an empty set for disjoint sets", func() {
			a := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
			b := NamesOf(fixtures.MessageC1, fixtures.MessageD1) // disjoint to a
			c := NamesOf(fixtures.MessageC1, fixtures.MessageD1) // same as c
			Expect(IntersectionN(a, b, c)).To(BeEmpty())
		})

		It("returns the intersection", func() {
			a := NamesOf(fixtures.MessageA1, fixtures.MessageB1, fixtures.MessageC1)
			b := NamesOf(fixtures.MessageB1, fixtures.MessageC1, fixtures.MessageD1)
			c := NamesOf(fixtures.MessageC1, fixtures.MessageD1, fixtures.MessageE1)
			Expect(IntersectionN(a, b, c)).To(Equal(NamesOf(fixtures.MessageC1)))
		})
	})

	Describe("func UnionN()", func() {
		It("returns an empty set if no sets are given", func() {
			Expect(UnionN()).To(BeEmpty())
		})

		It("returns the original set if a single set is given", func() {
			a := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
			Expect(UnionN(a)).To(Equal(a))
		})

		It("returns the original set for identical sets", func() {
			a := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
			b := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
			c := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
			Expect(UnionN(a, b, c)).To(Equal(a))
		})

		It("returns the union", func() {
			a := NamesOf(fixtures.MessageA1, fixtures.MessageB1, fixtures.MessageC1)
			b := NamesOf(fixtures.MessageB1, fixtures.MessageC1, fixtures.MessageD1)
			c := NamesOf(fixtures.MessageC1, fixtures.MessageD1, fixtures.MessageE1)

			Expect(UnionN(a, b, c)).To(Equal(NamesOf(
				fixtures.MessageA1,
				fixtures.MessageB1,
				fixtures.MessageC1,
				fixtures.MessageD1,
				fixtures.MessageE1,
			)))
		})
	})

	Describe("func DiffN()", func() {
		It("returns an empty set for identical sets", func() {
			a := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
			b := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
			Expect(DiffN(a, b)).To(BeEmpty())
		})

		It("returns an the original set for disjoint sets", func() {
			a := NamesOf(fixtures.MessageA1, fixtures.MessageB1)
			b := NamesOf(fixtures.MessageC1, fixtures.MessageD1)
			Expect(DiffN(a, b)).To(Equal(a))
		})

		It("returns the diff", func() {
			a := NamesOf(fixtures.MessageA1, fixtures.MessageB1, fixtures.MessageC1)
			b := NamesOf(fixtures.MessageB1, fixtures.MessageC1)
			Expect(DiffN(a, b)).To(Equal(NamesOf(fixtures.MessageA1)))
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

	Describe("func IsEqual()", func() {
		DescribeTable(
			"returns true if the sets are equivalent",
			func(a, b NameSet) {
				Expect(a.IsEqual(b)).To(BeTrue())
			},
			Entry(
				"equivalent",
				NewNameSet(MessageATypeName, MessageBTypeName),
				NewNameSet(MessageATypeName, MessageBTypeName),
			),
			Entry(
				"nil and empty",
				NewNameSet(),
				NameSet(nil),
			),
		)

		DescribeTable(
			"returns false if the sets are not equivalent",
			func(b NameSet) {
				a := NewNameSet(
					MessageATypeName,
					MessageBTypeName,
				)
				Expect(a.IsEqual(b)).To(BeFalse())
			},
			Entry(
				"subset",
				NewNameSet(MessageATypeName),
			),
			Entry(
				"superset",
				NewNameSet(MessageATypeName, MessageBTypeName, MessageCTypeName),
			),
			Entry(
				"same-length, disjoint",
				NewNameSet(MessageATypeName, MessageCTypeName),
			),
		)
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
