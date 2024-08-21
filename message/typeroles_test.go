package message_test

import (
	. "github.com/dogmatiq/configkit/message"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ TypeCollection = TypeRoles{}

var _ = Describe("type TypeRoles", func() {
	Describe("func Has()", func() {
		tr := TypeRoles{
			TypeFor[CommandStub[TypeA]](): CommandRole,
			TypeFor[EventStub[TypeA]]():   EventRole,
		}

		It("returns true if the type is in the map", func() {
			Expect(
				tr.Has(TypeFor[CommandStub[TypeA]]()),
			).To(BeTrue())
		})

		It("returns false if the type is not in the map", func() {
			Expect(
				tr.Has(TypeFor[EventStub[TypeX]]()),
			).To(BeFalse())
		})
	})

	Describe("func HasM()", func() {
		tr := TypeRoles{
			TypeFor[CommandStub[TypeA]](): CommandRole,
			TypeFor[EventStub[TypeA]]():   EventRole,
		}

		It("returns true if the type is in the map", func() {
			Expect(
				tr.HasM(CommandA1),
			).To(BeTrue())
		})

		It("returns false if the type is not in the map", func() {
			Expect(
				tr.HasM(CommandX1),
			).To(BeFalse())
		})
	})

	Describe("func Add()", func() {
		It("adds the type to the map", func() {
			t := TypeFor[CommandStub[TypeA]]()
			tr := TypeRoles{}

			tr.Add(t, CommandRole)

			Expect(
				tr.Has(t),
			).To(BeTrue())
		})

		It("returns true if the type is not already in the map", func() {
			tr := TypeRoles{}

			Expect(
				tr.Add(TypeFor[CommandStub[TypeA]](), CommandRole),
			).To(BeTrue())
		})

		It("returns false if the type is already in the map", func() {
			t := TypeFor[CommandStub[TypeA]]()
			tr := TypeRoles{}

			tr.Add(t, CommandRole)

			Expect(
				tr.Add(t, EventRole),
			).To(BeFalse())

			Expect(
				tr[t],
			).To(Equal(CommandRole))
		})
	})

	Describe("func AddM()", func() {
		It("adds the type of the message to the map", func() {
			tr := TypeRoles{}

			tr.AddM(CommandA1, CommandRole)

			Expect(
				tr.Has(TypeFor[CommandStub[TypeA]]()),
			).To(BeTrue())
		})

		It("returns true if the type is not already in the map", func() {
			tr := TypeRoles{}

			Expect(
				tr.AddM(CommandA1, CommandRole),
			).To(BeTrue())
		})

		It("returns false if the type is already in the map", func() {
			tr := TypeRoles{}
			tr.AddM(CommandA1, CommandRole)

			Expect(
				tr.AddM(CommandA1, EventRole),
			).To(BeFalse())

			Expect(
				tr[TypeFor[CommandStub[TypeA]]()],
			).To(Equal(CommandRole))
		})
	})

	Describe("func Remove()", func() {
		It("removes the type from the set", func() {
			t := TypeFor[CommandStub[TypeA]]()
			tr := TypeRoles{
				t: CommandRole,
			}

			tr.Remove(t)

			Expect(
				tr.Has(t),
			).To(BeFalse())
		})

		It("returns true if the type is already in the set", func() {
			t := TypeFor[CommandStub[TypeA]]()
			tr := TypeRoles{
				t: CommandRole,
			}

			Expect(
				tr.Remove(t),
			).To(BeTrue())
		})

		It("returns false if the type is not already in the set", func() {
			tr := TypeRoles{}

			Expect(
				tr.Remove(TypeFor[CommandStub[TypeA]]()),
			).To(BeFalse())
		})
	})

	Describe("func RemoveM()", func() {
		It("removes the type of the message from the set", func() {
			t := TypeFor[CommandStub[TypeA]]()
			tr := TypeRoles{
				t: CommandRole,
			}

			tr.RemoveM(CommandA1)

			Expect(
				tr.Has(t),
			).To(BeFalse())
		})

		It("returns true if the type is already in the set", func() {
			tr := TypeRoles{
				TypeFor[CommandStub[TypeA]](): CommandRole,
			}

			Expect(
				tr.RemoveM(CommandA1),
			).To(BeTrue())
		})

		It("returns false if the type is not already in the set", func() {
			tr := TypeRoles{}

			Expect(
				tr.RemoveM(CommandA1),
			).To(BeFalse())
		})
	})

	Describe("func IsEqual()", func() {
		DescribeTable(
			"returns true if the sets are equivalent",
			func(a, b TypeRoles) {
				Expect(a.IsEqual(b)).To(BeTrue())
			},
			Entry(
				"equivalent",
				TypeRoles{
					TypeFor[CommandStub[TypeA]](): CommandRole,
					TypeFor[EventStub[TypeA]]():   EventRole,
				},
				TypeRoles{
					TypeFor[CommandStub[TypeA]](): CommandRole,
					TypeFor[EventStub[TypeA]]():   EventRole,
				},
			),
			Entry(
				"nil and empty",
				TypeRoles{},
				TypeRoles(nil),
			),
		)

		DescribeTable(
			"returns false if the sets are not equivalent",
			func(b TypeRoles) {
				a := TypeRoles{
					TypeFor[CommandStub[TypeA]](): CommandRole,
					TypeFor[EventStub[TypeA]]():   EventRole,
				}
				Expect(a.IsEqual(b)).To(BeFalse())
			},
			Entry(
				"subset",
				TypeRoles{
					TypeFor[CommandStub[TypeA]](): CommandRole,
				},
			),
			Entry(
				"superset",
				TypeRoles{
					TypeFor[CommandStub[TypeA]](): CommandRole,
					TypeFor[EventStub[TypeA]]():   EventRole,
					TypeFor[TimeoutStub[TypeA]](): TimeoutRole,
				},
			),
			Entry(
				"same-length, disjoint type",
				TypeRoles{
					TypeFor[CommandStub[TypeA]](): CommandRole,
					TypeFor[TimeoutStub[TypeA]](): EventRole,
				},
			),
			Entry(
				"same-length, disjoint role",
				TypeRoles{
					TypeFor[CommandStub[TypeA]](): CommandRole,
					TypeFor[EventStub[TypeA]]():   TimeoutRole,
				},
			),
		)
	})

	Describe("func Len()", func() {
		It("returns the number of types in the collection", func() {
			tr := TypeRoles{
				TypeFor[CommandStub[TypeA]](): CommandRole,
				TypeFor[EventStub[TypeA]]():   EventRole,
			}

			Expect(tr.Len()).To(Equal(2))
		})
	})

	Describe("func Range()", func() {
		tr := TypeRoles{
			TypeFor[CommandStub[TypeA]](): CommandRole,
			TypeFor[EventStub[TypeA]]():   EventRole,
		}

		It("calls fn for each type in the container", func() {
			var types []Type

			all := tr.Range(func(t Type) bool {
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

			all := tr.Range(func(t Type) bool {
				count++
				return false
			})

			Expect(count).To(BeNumerically("==", 1))
			Expect(all).To(BeFalse())
		})
	})

	Describe("func RangeByRole()", func() {
		tr := TypeRoles{
			TypeFor[CommandStub[TypeA]](): CommandRole,
			TypeFor[CommandStub[TypeB]](): CommandRole,
			TypeFor[EventStub[TypeA]]():   EventRole,
		}

		It("calls fn for each type in the container with the given role", func() {
			var types []Type

			all := tr.RangeByRole(
				CommandRole,
				func(n Type) bool {
					types = append(types, n)
					return true
				},
			)

			Expect(types).To(ConsistOf(
				TypeFor[CommandStub[TypeA]](),
				TypeFor[CommandStub[TypeB]](),
			))
			Expect(all).To(BeTrue())
		})

		It("stops iterating if fn returns false", func() {
			count := 0

			all := tr.RangeByRole(
				CommandRole,
				func(n Type) bool {
					count++
					return false
				},
			)

			Expect(count).To(BeNumerically("==", 1))
			Expect(all).To(BeFalse())
		})
	})

	Describe("func FilterByRole()", func() {
		It("returns a subset containing only the given roles", func() {
			tr := TypeRoles{
				TypeFor[CommandStub[TypeA]](): CommandRole,
				TypeFor[CommandStub[TypeB]](): CommandRole,
				TypeFor[EventStub[TypeA]]():   EventRole,
			}

			subset := tr.FilterByRole(CommandRole)

			Expect(subset).To(Equal(
				TypeRoles{
					TypeFor[CommandStub[TypeA]](): CommandRole,
					TypeFor[CommandStub[TypeB]](): CommandRole,
				},
			))
		})
	})
})
