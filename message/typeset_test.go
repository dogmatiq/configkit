package message_test

import (
	. "github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma/fixtures"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ TypeCollection = TypeSet{}

var _ = Describe("type TypeSet", func() {
	Describe("func NewTypeSet()", func() {
		It("returns a set containing the given types", func() {
			Expect(NewTypeSet(
				TypeFor[CommandStub[TypeA]](),
				TypeFor[EventStub[TypeA]](),
			)).To(Equal(TypeSet{
				TypeFor[CommandStub[TypeA]](): struct{}{},
				TypeFor[EventStub[TypeA]]():   struct{}{},
			}))
		})
	})

	Describe("func TypesOf()", func() {
		It("returns a set containing the types of the given messages", func() {
			Expect(TypesOf(
				CommandA1,
				EventA1,
			)).To(Equal(TypeSet{
				TypeFor[CommandStub[TypeA]](): struct{}{},
				TypeFor[EventStub[TypeA]]():   struct{}{},
			}))
		})
	})

	Describe("func IntersectionT()", func() {
		It("returns an empty set if no sets are given", func() {
			Expect(IntersectionT()).To(BeEmpty())
		})

		It("returns the original set if a single set is given", func() {
			a := TypesOf(CommandA1, EventA1)
			Expect(IntersectionT(a)).To(Equal(a))
		})

		It("returns the original set for identical sets", func() {
			a := TypesOf(CommandA1, EventA1)
			b := TypesOf(CommandA1, EventA1)
			c := TypesOf(CommandA1, EventA1)
			Expect(IntersectionT(a, b, c)).To(Equal(a))
		})

		It("returns an empty set for disjoint sets", func() {
			a := TypesOf(CommandA1, EventA1)
			b := TypesOf(CommandB1, EventB1) // disjoint to a
			c := TypesOf(CommandB1, EventB1) // same as c
			Expect(IntersectionT(a, b, c)).To(BeEmpty())
		})

		It("returns the intersection", func() {
			a := TypesOf(CommandA1, EventA1, CommandB1)
			b := TypesOf(EventA1, CommandB1, EventB1)
			c := TypesOf(CommandB1, EventB1, TimeoutA1)
			Expect(IntersectionT(a, b, c)).To(Equal(TypesOf(CommandB1)))
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
		set := NewTypeSet(
			TypeFor[CommandStub[TypeA]](),
			TypeFor[EventStub[TypeA]](),
		)

		It("returns true if the type is in the set", func() {
			Expect(
				set.Has(TypeFor[CommandStub[TypeA]]()),
			).To(BeTrue())
		})

		It("returns false if the type is not in the set", func() {
			Expect(
				set.Has(TypeFor[CommandStub[TypeX]]()),
			).To(BeFalse())
		})
	})

	Describe("func HasM()", func() {
		set := TypesOf(
			CommandA1,
			EventA1,
		)

		It("returns true if the type is in the set", func() {
			Expect(
				set.HasM(CommandA1),
			).To(BeTrue())
		})

		It("returns false if the type is not in the set", func() {
			Expect(
				set.HasM(CommandX1),
			).To(BeFalse())
		})
	})

	Describe(("func Add()"), func() {
		It("adds the type to the set", func() {
			t := TypeFor[CommandStub[TypeA]]()
			s := NewTypeSet()

			s.Add(t)

			Expect(
				s.Has(t),
			).To(BeTrue())
		})

		It("returns true if the type is not already in the set", func() {
			s := NewTypeSet()

			Expect(
				s.Add(TypeFor[CommandStub[TypeA]]()),
			).To(BeTrue())
		})

		It("returns false if the type is already in the set", func() {
			t := TypeFor[CommandStub[TypeA]]()
			s := NewTypeSet()

			s.Add(t)

			Expect(
				s.Add(t),
			).To(BeFalse())
		})
	})

	Describe("func AddM()", func() {
		It("adds the type of the message to the set", func() {
			s := NewTypeSet()

			s.AddM(CommandA1)

			Expect(
				s.Has(TypeFor[CommandStub[TypeA]]()),
			).To(BeTrue())
		})

		It("returns true if the type is not already in the set", func() {
			s := NewTypeSet()

			Expect(
				s.AddM(CommandA1),
			).To(BeTrue())
		})

		It("returns false if the type is already in the set", func() {
			s := NewTypeSet()
			s.Add(TypeFor[CommandStub[TypeA]]())

			Expect(
				s.AddM(CommandA1),
			).To(BeFalse())
		})
	})

	Describe("func Remove()", func() {
		It("removes the type from the set", func() {
			t := TypeFor[CommandStub[TypeA]]()
			s := NewTypeSet(t)

			s.Remove(t)

			Expect(
				s.Has(t),
			).To(BeFalse())
		})

		It("returns true if the type is already in the set", func() {
			t := TypeFor[CommandStub[TypeA]]()
			s := NewTypeSet()

			s.Add(t)

			Expect(
				s.Remove(t),
			).To(BeTrue())
		})

		It("returns false if the type is not already in the set", func() {
			s := NewTypeSet()

			Expect(
				s.Remove(TypeFor[CommandStub[TypeA]]()),
			).To(BeFalse())
		})
	})

	Describe("func RemoveM()", func() {
		It("removes the type of the message from the set", func() {
			t := TypeFor[CommandStub[TypeA]]()
			s := NewTypeSet(t)

			s.RemoveM(CommandA1)

			Expect(
				s.Has(t),
			).To(BeFalse())
		})

		It("returns true if the type is already in the set", func() {
			s := NewTypeSet()

			s.Add(TypeFor[CommandStub[TypeA]]())

			Expect(
				s.RemoveM(CommandA1),
			).To(BeTrue())
		})

		It("returns false if the type is not already in the set", func() {
			s := NewTypeSet()

			Expect(
				s.RemoveM(CommandA1),
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
				NewTypeSet(
					TypeFor[CommandStub[TypeA]](),
					TypeFor[EventStub[TypeA]](),
				),
				NewTypeSet(
					TypeFor[CommandStub[TypeA]](),
					TypeFor[EventStub[TypeA]](),
				),
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
					TypeFor[CommandStub[TypeA]](),
					TypeFor[EventStub[TypeA]](),
				)
				Expect(a.IsEqual(b)).To(BeFalse())
			},
			Entry(
				"subset",
				NewTypeSet(
					TypeFor[CommandStub[TypeA]](),
				),
			),
			Entry(
				"superset",
				NewTypeSet(
					TypeFor[CommandStub[TypeA]](),
					TypeFor[EventStub[TypeA]](),
					TypeFor[TimeoutStub[TypeA]](),
				),
			),
			Entry(
				"same-length, disjoint",
				NewTypeSet(
					TypeFor[CommandStub[TypeA]](),
					TypeFor[TimeoutStub[TypeA]](),
				),
			),
		)
	})

	Describe("func Len()", func() {
		It("returns the number of types in the collection", func() {
			s := NewTypeSet(
				TypeFor[CommandStub[TypeA]](),
				TypeFor[EventStub[TypeA]](),
			)

			Expect(s.Len()).To(Equal(2))
		})
	})

	Describe("func Range()", func() {
		s := NewTypeSet(
			TypeFor[CommandStub[TypeA]](),
			TypeFor[EventStub[TypeA]](),
		)

		It("calls fn for each type in the container", func() {
			var types []Type

			all := s.Range(func(t Type) bool {
				types = append(types, t)
				return true
			})

			Expect(types).To(ConsistOf(
				TypeFor[CommandStub[TypeA]](),
				TypeFor[EventStub[TypeA]](),
			))
			Expect(all).To(BeTrue())
		})

		It("stops iterating if fn returns false", func() {
			count := 0

			all := s.Range(func(t Type) bool {
				count++
				return false
			})

			Expect(count).To(BeNumerically("==", 1))
			Expect(all).To(BeFalse())
		})
	})
})
