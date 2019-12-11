package validation_test

import (
	. "github.com/dogmatiq/configkit/internal/validation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Error", func() {
	Describe("func Error()", func() {
		It("returns the error message", func() {
			err := Error("<message>")
			Expect(err.Error()).To(Equal("<message>"))
		})
	})
})

var _ = Describe("func Panicf()", func() {
	It("panics", func() {
		Expect(func() {
			Panicf("<%s>", "message")
		}).To(Panic())
	})
})

var _ = Describe("func Errorf()", func() {
	It("returns a validation error", func() {
		err := Errorf("<%s>", "message")
		Expect(err).To(Equal(Error("<message>")))
	})
})
