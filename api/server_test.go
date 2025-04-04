package api_test

import (
	. "github.com/dogmatiq/configkit/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func NewServer()", func() {
	It("panics if one of the applications can not be marshaled", func() {
		Expect(func() {
			NewServer(nil)
		}).To(Panic())
	})
})
