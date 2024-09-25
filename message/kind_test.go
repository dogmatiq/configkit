package message_test

import (
	. "github.com/dogmatiq/configkit/message"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Kind", func() {
	Describe("func KindOf()", func() {
		It("returns the kind", func() {
			Expect(KindOf(CommandA1)).To(Equal(CommandKind))
			Expect(KindOf(EventA1)).To(Equal(EventKind))
			Expect(KindOf(TimeoutA1)).To(Equal(TimeoutKind))
		})

		It("panics if the message is nil", func() {
			Expect(func() {
				KindOf(nil)
			}).To(Panic())
		})
	})

	Describe("func KindFor()", func() {
		It("returns the kind", func() {
			Expect(KindFor[CommandStub[TypeA]]()).To(Equal(CommandKind))
			Expect(KindFor[EventStub[TypeA]]()).To(Equal(EventKind))
			Expect(KindFor[TimeoutStub[TypeA]]()).To(Equal(TimeoutKind))
		})
	})

	Describe("func Symbol()", func() {
		It("returns the symbol representation of the kind", func() {
			Expect(CommandKind.Symbol()).To(Equal("?"))
			Expect(EventKind.Symbol()).To(Equal("!"))
			Expect(TimeoutKind.Symbol()).To(Equal("@"))
		})
	})

	Describe("func String()", func() {
		It("returns the string representation of the kind", func() {
			Expect(CommandKind.String()).To(Equal("command"))
			Expect(EventKind.String()).To(Equal("event"))
			Expect(TimeoutKind.String()).To(Equal("timeout"))
		})
	})
})
