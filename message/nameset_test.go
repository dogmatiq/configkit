package message_test

import (
	. "github.com/dogmatiq/configkit/message"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ NameCollection = NameSet{}

var _ = Describe("type NameSet", func() {
	Describe("func NewNameSet()", func() {
		It("returns a set containing the given names", func() {
			Expect(NewNameSet(
				NameFor[CommandStub[TypeA]](),
				NameFor[EventStub[TypeA]](),
			)).To(Equal(NameSet{
				NameFor[CommandStub[TypeA]](): struct{}{},
				NameFor[EventStub[TypeA]]():   struct{}{},
			}))
		})
	})

	Describe("func NamesOf()", func() {
		It("returns a set containing the names of the given messages", func() {
			Expect(NamesOf(
				CommandA1,
				EventA1,
			)).To(Equal(NameSet{
				NameFor[CommandStub[TypeA]](): struct{}{},
				NameFor[EventStub[TypeA]]():   struct{}{},
			}))
		})
	})

	Describe("func IntersectionN()", func() {
		It("returns an empty set if no sets are given", func() {
			Expect(IntersectionN()).To(BeEmpty())
		})

		It("returns the original set if a single set is given", func() {
			a := NamesOf(CommandA1, EventA1)
			Expect(IntersectionN(a)).To(Equal(a))
		})

		It("returns the original set for identical sets", func() {
			a := NamesOf(CommandA1, EventA1)
			b := NamesOf(CommandA1, EventA1)
			c := NamesOf(CommandA1, EventA1)
			Expect(IntersectionN(a, b, c)).To(Equal(a))
		})

		It("returns an empty set for disjoint sets", func() {
			a := NamesOf(CommandA1, EventA1)
			b := NamesOf(CommandB1, EventB1) // disjoint to a
			c := NamesOf(CommandB1, EventB1) // same as c
			Expect(IntersectionN(a, b, c)).To(BeEmpty())
		})

		It("returns the intersection", func() {
			a := NamesOf(CommandA1, EventA1, CommandB1)
			b := NamesOf(EventA1, CommandB1, EventB1)
			c := NamesOf(CommandB1, EventB1, TimeoutA1)
			Expect(IntersectionN(a, b, c)).To(Equal(NamesOf(CommandB1)))
		})
	})

	Describe("func UnionN()", func() {
		It("returns an empty set if no sets are given", func() {
			Expect(UnionN()).To(BeEmpty())
		})

		It("returns the original set if a single set is given", func() {
			a := NamesOf(CommandA1, EventA1)
			Expect(UnionN(a)).To(Equal(a))
		})

		It("returns the original set for identical sets", func() {
			a := NamesOf(CommandA1, EventA1)
			b := NamesOf(CommandA1, EventA1)
			c := NamesOf(CommandA1, EventA1)
			Expect(UnionN(a, b, c)).To(Equal(a))
		})

		It("returns the union", func() {
			a := NamesOf(CommandA1, EventA1, CommandB1)
			b := NamesOf(EventA1, CommandB1, EventB1)
			c := NamesOf(CommandB1, EventB1, TimeoutA1)

			Expect(UnionN(a, b, c)).To(Equal(NamesOf(
				CommandA1,
				EventA1,
				CommandB1,
				EventB1,
				TimeoutA1,
			)))
		})
	})

	Describe("func DiffN()", func() {
		It("returns an empty set for identical sets", func() {
			a := NamesOf(CommandA1, EventA1)
			b := NamesOf(CommandA1, EventA1)
			Expect(DiffN(a, b)).To(BeEmpty())
		})

		It("returns an the original set for disjoint sets", func() {
			a := NamesOf(CommandA1, EventA1)
			b := NamesOf(CommandB1, EventB1)
			Expect(DiffN(a, b)).To(Equal(a))
		})

		It("returns the diff", func() {
			a := NamesOf(CommandA1, EventA1, CommandB1)
			b := NamesOf(EventA1, CommandB1)
			Expect(DiffN(a, b)).To(Equal(NamesOf(CommandA1)))
		})
	})

	Describe("func Has()", func() {
		set := NewNameSet(
			NameFor[CommandStub[TypeA]](),
			NameFor[EventStub[TypeA]](),
		)

		It("returns true if the name is in the set", func() {
			Expect(
				set.Has(NameFor[CommandStub[TypeA]]()),
			).To(BeTrue())
		})

		It("returns false if the name is not in the set", func() {
			Expect(
				set.Has(NameFor[CommandStub[TypeX]]()),
			).To(BeFalse())
		})
	})

	Describe("func HasM()", func() {
		set := NewNameSet(
			NameFor[CommandStub[TypeA]](),
			NameFor[EventStub[TypeA]](),
		)

		It("returns true if the name is in the set", func() {
			Expect(
				set.HasM(CommandA1),
			).To(BeTrue())
		})

		It("returns false if the name is not in the set", func() {
			Expect(
				set.HasM(CommandX1),
			).To(BeFalse())
		})
	})

	Describe(("func Add()"), func() {
		It("adds the name to the set", func() {
			n := NameFor[CommandStub[TypeA]]()
			s := NewNameSet()

			s.Add(n)

			Expect(
				s.Has(n),
			).To(BeTrue())
		})

		It("returns true if the name is not already in the set", func() {
			s := NewNameSet()

			Expect(
				s.Add(NameFor[CommandStub[TypeA]]()),
			).To(BeTrue())
		})

		It("returns false if the name is already in the set", func() {
			n := NameFor[CommandStub[TypeA]]()
			s := NewNameSet()

			s.Add(n)

			Expect(
				s.Add(n),
			).To(BeFalse())
		})
	})

	Describe("func AddM()", func() {
		It("adds the name of the message to the set", func() {
			s := NewNameSet()
			s.AddM(CommandA1)

			Expect(
				s.Has(NameFor[CommandStub[TypeA]]()),
			).To(BeTrue())
		})

		It("returns true if the name is not already in the set", func() {
			s := NewNameSet()

			Expect(
				s.AddM(CommandA1),
			).To(BeTrue())
		})

		It("returns false if the name is already in the set", func() {
			s := NewNameSet()

			s.Add(NameFor[CommandStub[TypeA]]())

			Expect(
				s.AddM(CommandA1),
			).To(BeFalse())
		})
	})

	Describe("func Remove()", func() {
		It("removes the name from the set", func() {
			n := NameFor[CommandStub[TypeA]]()
			s := NewNameSet(n)

			s.Remove(n)

			Expect(
				s.Has(n),
			).To(BeFalse())
		})

		It("returns true if the name is already in the set", func() {
			n := NameFor[CommandStub[TypeA]]()
			s := NewNameSet()

			s.Add(n)

			Expect(
				s.Remove(n),
			).To(BeTrue())
		})

		It("returns false if the name is not already in the set", func() {
			s := NewNameSet()

			Expect(
				s.Remove(NameFor[CommandStub[TypeA]]()),
			).To(BeFalse())
		})
	})

	Describe("func RemoveM()", func() {
		It("removes the name of the message from the set", func() {
			n := NameFor[CommandStub[TypeA]]()
			s := NewNameSet(n)

			s.RemoveM(CommandA1)

			Expect(
				s.Has(n),
			).To(BeFalse())
		})

		It("returns true if the name is already in the set", func() {
			s := NewNameSet()

			s.Add(NameFor[CommandStub[TypeA]]())

			Expect(
				s.RemoveM(CommandA1),
			).To(BeTrue())
		})

		It("returns false if the name is not already in the set", func() {
			s := NewNameSet()

			Expect(
				s.RemoveM(CommandA1),
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
				NewNameSet(
					NameFor[CommandStub[TypeA]](),
					NameFor[EventStub[TypeA]](),
				),
				NewNameSet(
					NameFor[CommandStub[TypeA]](),
					NameFor[EventStub[TypeA]](),
				),
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
					NameFor[CommandStub[TypeA]](),
					NameFor[EventStub[TypeA]](),
				)
				Expect(a.IsEqual(b)).To(BeFalse())
			},
			Entry(
				"subset",
				NewNameSet(
					NameFor[CommandStub[TypeA]](),
				),
			),
			Entry(
				"superset",
				NewNameSet(
					NameFor[CommandStub[TypeA]](),
					NameFor[EventStub[TypeA]](),
					NameFor[TimeoutStub[TypeA]](),
				),
			),
			Entry(
				"same-length, disjoint",
				NewNameSet(
					NameFor[CommandStub[TypeA]](),
					NameFor[TimeoutStub[TypeA]](),
				),
			),
		)
	})

	Describe("func Len()", func() {
		It("returns the number of names in the collection", func() {
			s := NewNameSet(
				NameFor[CommandStub[TypeA]](),
				NameFor[EventStub[TypeA]](),
			)

			Expect(s.Len()).To(Equal(2))
		})
	})

	Describe("func Range()", func() {
		s := NewNameSet(
			NameFor[CommandStub[TypeA]](),
			NameFor[EventStub[TypeA]](),
		)

		It("calls fn for each name in the container", func() {
			var names []Name

			all := s.Range(func(n Name) bool {
				names = append(names, n)
				return true
			})

			Expect(names).To(ConsistOf(
				NameFor[CommandStub[TypeA]](),
				NameFor[EventStub[TypeA]](),
			))
			Expect(all).To(BeTrue())
		})

		It("stops iterating if fn returns false", func() {
			count := 0

			all := s.Range(func(n Name) bool {
				count++
				return false
			})

			Expect(count).To(BeNumerically("==", 1))
			Expect(all).To(BeFalse())
		})
	})
})
