package api

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
)

var _ = Describe("func RegisterServer()", func() {
	It("panics if one of the applications can not be marshaled", func() {
		Expect(func() {
			s := grpc.NewServer()
			RegisterServer(s, &application{})
		}).To(Panic())
	})
})
