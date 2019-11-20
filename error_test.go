package configkit

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type ValidationError", func() {
	Describe("func Error", func() {
		It("returns the error message", func() {
			err := ValidationError("<message>")
			Expect(err.Error()).To(Equal("<message>"))
		})
	})
})

var _ = Describe("func catch", func() {
	It("returns an error if the function panics with a ValidatioNError", func() {
		err := catch(func() {
			panicf("<message>")
		})
		Expect(err).To(Equal(ValidationError("<message>")))
	})

	It("returns nil if no panic occurs", func() {
		err := catch(func() {})
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("does not catch non-config panics", func() {
		Expect(func() {
			catch(func() {
				panic("<panic>")
			})
		}).To(Panic())
	})
})
