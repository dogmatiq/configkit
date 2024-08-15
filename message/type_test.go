package message_test

import (
	"reflect"

	. "github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func IsEqualT()", func() {
	It("returns true for identical sets", func() {
		a := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
		b := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
		Expect(IsEqualSetT(a, b)).To(BeTrue())
	})

	It("returns false for disjoint sets", func() {
		a := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
		b := TypesOf(fixtures.MessageC1, fixtures.MessageD1)
		Expect(IsEqualSetT(a, b)).To(BeFalse())
	})

	It("returns false for intersecting sets", func() {
		a := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
		b := TypesOf(fixtures.MessageB1, fixtures.MessageC1)
		Expect(IsEqualSetT(a, b)).To(BeFalse())
	})
})

var _ = Describe("func IsIntersectingT()", func() {
	It("returns true for identical sets", func() {
		a := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
		b := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
		Expect(IsIntersectingT(a, b)).To(BeTrue())
	})

	It("returns false for disjoint sets", func() {
		a := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
		b := TypesOf(fixtures.MessageC1, fixtures.MessageD1)
		Expect(IsIntersectingT(a, b)).To(BeFalse())
	})

	It("returns true for intersecting sets", func() {
		a := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
		b := TypesOf(fixtures.MessageB1, fixtures.MessageC1)
		Expect(IsIntersectingT(a, b)).To(BeTrue())
	})
})

var _ = Describe("func IsSubsetT()", func() {
	It("returns true for identical sets", func() {
		a := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
		b := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
		Expect(IsSubsetT(a, b)).To(BeTrue())
	})

	It("returns true for strict subsets", func() {
		a := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
		b := TypesOf(fixtures.MessageA1, fixtures.MessageB1, fixtures.MessageC1)
		Expect(IsSubsetT(a, b)).To(BeTrue())
	})

	It("returns false for supersets", func() {
		a := TypesOf(fixtures.MessageA1, fixtures.MessageB1, fixtures.MessageC1)
		b := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
		Expect(IsSubsetT(a, b)).To(BeFalse())
	})

	It("returns false for disjoint sets", func() {
		a := TypesOf(fixtures.MessageA1, fixtures.MessageB1)
		b := TypesOf(fixtures.MessageC1, fixtures.MessageD1)
		Expect(IsSubsetT(a, b)).To(BeFalse())
	})
})

var _ = Describe("type Type", func() {
	Describe("func TypeOf()", func() {
		It("returns values that compare as equal for messages of the same type", func() {
			tb := TypeOf(fixtures.MessageA1)
			ta := TypeOf(fixtures.MessageA1)

			Expect(ta).To(Equal(tb))
			Expect(ta == tb).To(BeTrue()) // we're testing == here specifically, hence not using To(Equal())
		})

		It("returns values that do not compare as equal for messages of different types", func() {
			ta := TypeOf(fixtures.MessageA1)
			tb := TypeOf(fixtures.MessageB1)

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
			tb := TypeOf(fixtures.MessageA1)
			ta := TypeFor[fixtures.MessageA]()

			Expect(ta).To(Equal(tb))
			Expect(ta == tb).To(BeTrue()) // we're testing == here specifically, hence not using To(Equal())
		})
	})

	Describe("func TypeFromReflect()", func() {
		It("returns values that compare as equal for messages of the same type", func() {
			tb := TypeFromReflect(reflect.TypeOf(fixtures.MessageA1))
			ta := TypeFromReflect(reflect.TypeOf(fixtures.MessageA1))

			Expect(ta).To(Equal(tb))
			Expect(ta == tb).To(BeTrue()) // we're testing == here specifically, hence not using To(Equal())
		})

		It("returns values that do not compare as equal for messages of different types", func() {
			ta := TypeFromReflect(reflect.TypeOf(fixtures.MessageA1))
			tb := TypeFromReflect(reflect.TypeOf(fixtures.MessageB1))

			Expect(ta).NotTo(Equal(tb))
			Expect(ta != tb).To(BeTrue()) // we're testing != here specifically, hence not using NotTo(Equal())
		})
	})

	Describe("func Name()", func() {
		It("returns the fully qualified type name", func() {
			mt := TypeOf(fixtures.MessageA1)
			Expect(mt.Name()).To(Equal(NameOf(fixtures.MessageA1)))
		})

		It("panics if the type is the zero-value", func() {
			Expect(func() {
				Type{}.Name()
			}).To(Panic())
		})
	})

	Describe("func ReflectType()", func() {
		It("returns the reflect.Type for the message", func() {
			mt := TypeOf(fixtures.MessageA1)
			rt := reflect.TypeOf(fixtures.MessageA1)

			Expect(mt.ReflectType()).To(BeIdenticalTo(rt))
		})

		It("panics if the type is the zero-value", func() {
			Expect(func() {
				Type{}.ReflectType()
			}).To(Panic())
		})
	})

	Describe("func String()", func() {
		It("returns the string representation of the type", func() {
			t := TypeOf(fixtures.MessageA1)
			Expect(t.String()).To(Equal("fixtures.MessageA"))
		})
	})
})
