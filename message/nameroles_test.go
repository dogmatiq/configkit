package message_test

import (
	. "github.com/dogmatiq/configkit/message"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ NameCollection = NameRoles{}

var _ = Describe("type NameRoles", func() {
	Describe("func Has()", func() {
		nr := NameRoles{
			NameFor[CommandStub[TypeA]](): CommandRole,
			NameFor[EventStub[TypeA]]():   EventRole,
		}

		It("returns true if the name is in the map", func() {
			Expect(
				nr.Has(NameFor[CommandStub[TypeA]]()),
			).To(BeTrue())
		})

		It("returns false if the name is not in the map", func() {
			Expect(
				nr.Has(NameFor[CommandStub[TypeX]]()),
			).To(BeFalse())
		})
	})

	Describe("func HasM()", func() {
		nr := NameRoles{
			NameFor[CommandStub[TypeA]](): CommandRole,
			NameFor[EventStub[TypeA]]():   EventRole,
		}

		It("returns true if the name is in the map", func() {
			Expect(
				nr.HasM(CommandA1),
			).To(BeTrue())
		})

		It("returns false if the name is not in the map", func() {
			Expect(
				nr.HasM(CommandX1),
			).To(BeFalse())
		})
	})

	Describe("func Add()", func() {
		It("adds the name to the map", func() {
			n := NameFor[CommandStub[TypeA]]()
			nr := NameRoles{}

			nr.Add(n, CommandRole)

			Expect(
				nr.Has(n),
			).To(BeTrue())
		})

		It("returns true if the name is not already in the map", func() {
			nr := NameRoles{}

			Expect(
				nr.Add(NameFor[CommandStub[TypeA]](), CommandRole),
			).To(BeTrue())
		})

		It("returns false if the name is already in the map", func() {
			n := NameFor[CommandStub[TypeA]]()
			nr := NameRoles{}

			nr.Add(n, CommandRole)

			Expect(
				nr.Add(n, EventRole),
			).To(BeFalse())

			Expect(
				nr[n],
			).To(Equal(CommandRole))
		})
	})

	Describe("func AddM()", func() {
		It("adds the name of the message to the map", func() {
			nr := NameRoles{}
			nr.AddM(CommandA1, CommandRole)

			Expect(
				nr.Has(NameFor[CommandStub[TypeA]]()),
			).To(BeTrue())
		})

		It("returns true if the name is not already in the map", func() {
			nr := NameRoles{}

			Expect(
				nr.AddM(CommandA1, CommandRole),
			).To(BeTrue())
		})

		It("returns false if the name is already in the map", func() {
			nr := NameRoles{}
			nr.AddM(CommandA1, CommandRole)

			Expect(
				nr.AddM(CommandA1, EventRole),
			).To(BeFalse())

			Expect(
				nr[NameFor[CommandStub[TypeA]]()],
			).To(Equal(CommandRole))
		})
	})

	Describe("func Remove()", func() {
		It("removes the name from the set", func() {
			n := NameFor[CommandStub[TypeA]]()
			nr := NameRoles{
				n: CommandRole,
			}

			nr.Remove(n)

			Expect(
				nr.Has(n),
			).To(BeFalse())
		})

		It("returns true if the name is already in the set", func() {
			n := NameFor[CommandStub[TypeA]]()
			nr := NameRoles{
				n: CommandRole,
			}

			Expect(
				nr.Remove(n),
			).To(BeTrue())
		})

		It("returns false if the name is not already in the set", func() {
			nr := NameRoles{}

			Expect(
				nr.Remove(NameFor[CommandStub[TypeA]]()),
			).To(BeFalse())
		})
	})

	Describe("func RemoveM()", func() {
		It("removes the name of the message from the set", func() {
			n := NameFor[CommandStub[TypeA]]()
			nr := NameRoles{
				n: CommandRole,
			}

			nr.RemoveM(CommandA1)

			Expect(
				nr.Has(n),
			).To(BeFalse())
		})

		It("returns true if the name is already in the set", func() {
			nr := NameRoles{
				NameFor[CommandStub[TypeA]](): CommandRole,
			}

			Expect(
				nr.RemoveM(CommandA1),
			).To(BeTrue())
		})

		It("returns false if the name is not already in the set", func() {
			nr := NameRoles{}

			Expect(
				nr.RemoveM(CommandA1),
			).To(BeFalse())
		})
	})

	Describe("func IsEqual()", func() {
		DescribeTable(
			"returns true if the sets are equivalent",
			func(a, b NameRoles) {
				Expect(a.IsEqual(b)).To(BeTrue())
			},
			Entry(
				"equivalent",
				NameRoles{
					NameFor[CommandStub[TypeA]](): CommandRole,
					NameFor[EventStub[TypeA]]():   EventRole,
				},
				NameRoles{
					NameFor[CommandStub[TypeA]](): CommandRole,
					NameFor[EventStub[TypeA]]():   EventRole,
				},
			),
			Entry(
				"nil and empty",
				NameRoles{},
				NameRoles(nil),
			),
		)

		DescribeTable(
			"returns false if the sets are not equivalent",
			func(b NameRoles) {
				a := NameRoles{
					NameFor[CommandStub[TypeA]](): CommandRole,
					NameFor[EventStub[TypeA]]():   EventRole,
				}
				Expect(a.IsEqual(b)).To(BeFalse())
			},
			Entry(
				"subset",
				NameRoles{
					NameFor[CommandStub[TypeA]](): CommandRole,
				},
			),
			Entry(
				"superset",
				NameRoles{
					NameFor[CommandStub[TypeA]](): CommandRole,
					NameFor[EventStub[TypeA]]():   EventRole,
					NameFor[TimeoutStub[TypeA]](): TimeoutRole,
				},
			),
			Entry(
				"same-length, disjoint type",
				NameRoles{
					NameFor[CommandStub[TypeA]](): CommandRole,
					NameFor[EventStub[TypeB]]():   EventRole,
				},
			),
			Entry(
				"same-length, disjoint role",
				NameRoles{
					NameFor[CommandStub[TypeA]](): CommandRole,
					NameFor[EventStub[TypeA]]():   TimeoutRole,
				},
			),
		)
	})

	Describe("func Len()", func() {
		It("returns the number of names in the collection", func() {
			nr := NameRoles{
				NameFor[CommandStub[TypeA]](): CommandRole,
				NameFor[EventStub[TypeA]]():   EventRole,
			}

			Expect(nr.Len()).To(Equal(2))
		})
	})

	Describe("func Range()", func() {
		nr := NameRoles{
			NameFor[CommandStub[TypeA]](): CommandRole,
			NameFor[EventStub[TypeA]]():   EventRole,
		}

		It("calls fn for each name in the container", func() {
			var names []Name

			all := nr.Range(func(n Name) bool {
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

			all := nr.Range(func(n Name) bool {
				count++
				return false
			})

			Expect(count).To(BeNumerically("==", 1))
			Expect(all).To(BeFalse())
		})
	})

	Describe("func RangeByRole()", func() {
		nr := NameRoles{
			NameFor[CommandStub[TypeA]](): CommandRole,
			NameFor[CommandStub[TypeB]](): CommandRole,
			NameFor[EventStub[TypeA]]():   EventRole,
		}

		It("calls fn for each name in the container with the given role", func() {
			var names []Name

			all := nr.RangeByRole(
				CommandRole,
				func(n Name) bool {
					names = append(names, n)
					return true
				},
			)

			Expect(names).To(ConsistOf(
				NameFor[CommandStub[TypeA]](),
				NameFor[CommandStub[TypeB]](),
			))
			Expect(all).To(BeTrue())
		})

		It("stops iterating if fn returns false", func() {
			count := 0

			all := nr.RangeByRole(
				CommandRole,
				func(n Name) bool {
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
			nr := NameRoles{
				NameFor[CommandStub[TypeA]](): CommandRole,
				NameFor[CommandStub[TypeB]](): CommandRole,
				NameFor[EventStub[TypeA]]():   EventRole,
			}

			subset := nr.FilterByRole(CommandRole)

			Expect(subset).To(Equal(
				NameRoles{
					NameFor[CommandStub[TypeA]](): CommandRole,
					NameFor[CommandStub[TypeB]](): CommandRole,
				},
			))
		})
	})
})
