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

	Describe("func IntersectionT()", func() {
		It("returns an empty set if no sets are given", func() {
			Expect(IntersectionT()).To(BeEmpty())
		})

		It("returns the original set if a single set is given", func() {
			a := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
			Expect(IntersectionT(a)).To(Equal(a))
		})

		It("returns the original set for identical sets", func() {
			a := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
			b := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
			c := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
			Expect(IntersectionT(a, b, c)).To(Equal(a))
		})

		It("returns an empty set for disjoint sets", func() {
			a := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
			b := TypesOf(fixtures.MessageC1, fixtures.MessageD1) // disjoint to a
			c := TypesOf(fixtures.MessageC1, fixtures.MessageD1) // same as c
			Expect(IntersectionT(a, b, c)).To(BeEmpty())
		})

		It("returns the intersection", func() {
			a := TypesOf(fixtures.MessageA1, fixtures.MessageB1, fixtures.MessageC1)
			b := TypesOf(fixtures.MessageB1, fixtures.MessageC1, fixtures.MessageD1)
			c := TypesOf(fixtures.MessageC1, fixtures.MessageD1, fixtures.MessageE1)
			Expect(IntersectionT(a, b, c)).To(Equal(TypesOf(fixtures.MessageC1)))
		})
	})

	Describe("func UnionT()", func() {
		It("returns an empty set if no sets are given", func() {
			Expect(UnionT()).To(BeEmpty())
		})

		It("returns the original set if a single set is given", func() {
			a := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
			Expect(UnionT(a)).To(Equal(a))
		})

		It("returns the original set for identical sets", func() {
			a := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
			b := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
			c := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
			Expect(UnionT(a, b, c)).To(Equal(a))
		})

		It("returns the union", func() {
			a := TypesOf(fixtures.MessageA1, fixtures.MessageB1, fixtures.MessageC1)
			b := TypesOf(fixtures.MessageB1, fixtures.MessageC1, fixtures.MessageD1)
			c := TypesOf(fixtures.MessageC1, fixtures.MessageD1, fixtures.MessageE1)

			Expect(UnionT(a, b, c)).To(Equal(TypesOf(
				fixtures.MessageA1,
				fixtures.MessageB1,
				fixtures.MessageC1,
				fixtures.MessageD1,
				fixtures.MessageE1,
			)))
		})
	})

	Describe("func DiffT()", func() {
		It("returns an empty set for identical sets", func() {
			a := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
			b := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
			Expect(DiffT(a, b)).To(BeEmpty())
		})

		It("returns an the original set for disjoint sets", func() {
			a := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
			b := TypesOf(fixtures.MessageC1, fixtures.MessageD1)
			Expect(DiffT(a, b)).To(Equal(a))
		})

		It("returns the diff", func() {
			a := TypesOf(fixtures.MessageA1, fixtures.MessageB1, fixtures.MessageC1)
			b := TypesOf(fixtures.MessageB1, fixtures.MessageC1)
			Expect(DiffT(a, b)).To(Equal(TypesOf(fixtures.MessageA1)))
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

	Describe("func IsEqual()", func() {
		DescribeTable(
			"returns true if the sets are equivalent",
			func(a, b TypeSet) {
				Expect(a.IsEqual(b)).To(BeTrue())
			},
			Entry(
				"equivalent",
				NewTypeSet(MessageAType, MessageBType),
				NewTypeSet(MessageAType, MessageBType),
			),
			Entry(
				"nil and empty",
				NewTypeSet(),
				TypeSet(nil),
			),
		)

		DescribeTable(
			"returns false if the sets are not equivalent",
			func(b TypeSet) {
				a := NewTypeSet(
					MessageAType,
					MessageBType,
				)
				Expect(a.IsEqual(b)).To(BeFalse())
			},
			Entry(
				"subset",
				NewTypeSet(MessageAType),
			),
			Entry(
				"superset",
				NewTypeSet(MessageAType, MessageBType, MessageCType),
			),
			Entry(
				"same-length, disjoint",
				NewTypeSet(MessageAType, MessageCType),
			),
		)
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
