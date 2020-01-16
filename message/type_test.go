package message_test

import (
	"reflect"

	. "github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/dogma/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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
