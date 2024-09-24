package message_test

import (
	"reflect"

	. "github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func IsEqualT()", func() {
	It("returns true for identical sets", func() {
		a := TypesOf(CommandA1, EventA1)
		b := TypesOf(CommandA1, EventA1)
		Expect(IsEqualSetT(a, b)).To(BeTrue())
	})

	It("returns false for disjoint sets", func() {
		a := TypesOf(CommandA1, EventA1)
		b := TypesOf(CommandB1, EventB1)
		Expect(IsEqualSetT(a, b)).To(BeFalse())
	})

	It("returns false for intersecting sets", func() {
		a := TypesOf(CommandA1, EventA1)
		b := TypesOf(EventA1, CommandB1)
		Expect(IsEqualSetT(a, b)).To(BeFalse())
	})
})

var _ = Describe("func IsIntersectingT()", func() {
	It("returns true for identical sets", func() {
		a := TypesOf(CommandA1, EventA1)
		b := TypesOf(CommandA1, EventA1)
		Expect(IsIntersectingT(a, b)).To(BeTrue())
	})

	It("returns false for disjoint sets", func() {
		a := TypesOf(CommandA1, EventA1)
		b := TypesOf(CommandB1, EventB1)
		Expect(IsIntersectingT(a, b)).To(BeFalse())
	})

	It("returns true for intersecting sets", func() {
		a := TypesOf(CommandA1, EventA1)
		b := TypesOf(EventA1, CommandB1)
		Expect(IsIntersectingT(a, b)).To(BeTrue())
	})
})

var _ = Describe("func IsSubsetT()", func() {
	It("returns true for identical sets", func() {
		a := TypesOf(CommandA1, EventA1)
		b := TypesOf(CommandA1, EventA1)
		Expect(IsSubsetT(a, b)).To(BeTrue())
	})

	It("returns true for strict subsets", func() {
		a := TypesOf(CommandA1, EventA1)
		b := TypesOf(CommandA1, EventA1, CommandB1)
		Expect(IsSubsetT(a, b)).To(BeTrue())
	})

	It("returns false for supersets", func() {
		a := TypesOf(CommandA1, EventA1, CommandB1)
		b := TypesOf(CommandA1, EventA1)
		Expect(IsSubsetT(a, b)).To(BeFalse())
	})

	It("returns false for disjoint sets", func() {
		a := TypesOf(CommandA1, EventA1)
		b := TypesOf(CommandB1, EventB1)
		Expect(IsSubsetT(a, b)).To(BeFalse())
	})
})

var _ = Describe("type Type", func() {
	Describe("func TypeOf()", func() {
		It("returns values that compare as equal for messages of the same type", func() {
			tb := TypeOf(CommandA1)
			ta := TypeOf(CommandA1)

			Expect(ta).To(Equal(tb))
			Expect(ta == tb).To(BeTrue()) // we're testing == here specifically, hence not using To(Equal())
		})

		It("returns values that do not compare as equal for messages of different types", func() {
			ta := TypeOf(CommandA1)
			tb := TypeOf(EventA1)

			Expect(ta).NotTo(Equal(tb))
			Expect(ta != tb).To(BeTrue()) // we're testing != here specifically, hence not using NotTo(Equal())
		})

		It("panics if the message is nil", func() {
			Expect(func() {
				TypeOf(nil)
			}).To(Panic())
		})
	})

	Describe("func TypeFor()", func() {
		It("returns the type", func() {
			ta := TypeFor[CommandStub[TypeA]]()
			tb := TypeOf(CommandA1)

			Expect(ta).To(Equal(tb))
			Expect(ta == tb).To(BeTrue()) // we're testing == here specifically, hence not using To(Equal())
		})
	})

	Describe("func TypeFromReflect()", func() {
		It("returns values that compare as equal for messages of the same type", func() {
			tb := TypeFromReflect(reflect.TypeOf(CommandA1))
			ta := TypeFromReflect(reflect.TypeOf(CommandA1))

			Expect(ta).To(Equal(tb))
			Expect(ta == tb).To(BeTrue()) // we're testing == here specifically, hence not using To(Equal())
		})

		It("returns values that do not compare as equal for messages of different types", func() {
			ta := TypeFromReflect(reflect.TypeOf(CommandA1))
			tb := TypeFromReflect(reflect.TypeOf(EventA1))

			Expect(ta).NotTo(Equal(tb))
			Expect(ta != tb).To(BeTrue()) // we're testing != here specifically, hence not using NotTo(Equal())
		})
	})

	Describe("func Name()", func() {
		It("returns the fully qualified type name", func() {
			mt := TypeOf(CommandA1)
			Expect(mt.Name()).To(Equal(NameOf(CommandA1)))
		})

		It("panics if the type is the zero-value", func() {
			Expect(func() {
				Type{}.Name()
			}).To(Panic())
		})
	})

	Describe("func Role()", func() {
		It("returns the role of the message", func() {
			Expect(TypeOf(CommandA1).Role()).To(Equal(CommandRole))
			Expect(TypeOf(EventA1).Role()).To(Equal(EventRole))
			Expect(TypeOf(TimeoutA1).Role()).To(Equal(TimeoutRole))
		})
	})

	Describe("func ReflectType()", func() {
		It("returns the reflect.Type for the message", func() {
			mt := TypeOf(CommandA1)
			rt := reflect.TypeOf(CommandA1)

			Expect(mt.ReflectType()).To(BeIdenticalTo(rt))
		})

		It("panics if the type is the zero-value", func() {
			Expect(func() {
				Type{}.ReflectType()
			}).To(Panic())
		})
	})

	Describe("func String()", func() {
		type NonGenericCommand struct {
			CommandStub[TypeA]
		}

		DescribeTable(
			"it returns the string representation of the type",
			func(v dogma.Message, expected string) {
				t := TypeOf(v)
				Expect(t.String()).To(Equal(expected))
			},
			Entry("generic type", CommandA1, "stubs.CommandStub[TypeA]"),
			Entry("pointer to generic type", &CommandA1, "*stubs.CommandStub[TypeA]"),
			Entry("non-generic type", NonGenericCommand{}, "message_test.NonGenericCommand"),
			Entry("pointer to non-generic type", &NonGenericCommand{}, "*message_test.NonGenericCommand"),
		)
	})
})
