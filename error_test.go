package configkit_test

import (
	. "github.com/dogmatiq/configkit"
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/dogma/fixtures"
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
	It("does not catch non-config panics", func() {
		Expect(func() {
			NewAggregateConfig(&AggregateMessageHandler{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					panic("<panic>")
				},
			})
		}).To(Panic())
	})
})
