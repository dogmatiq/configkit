package configkit_test

import (
	. "github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type HandlerType", func() {
	Describe("func MustValidate()", func() {
		It("does not panic when the type is valid", func() {
			AggregateHandlerType.MustValidate()
			ProcessHandlerType.MustValidate()
			IntegrationHandlerType.MustValidate()
			ProjectionHandlerType.MustValidate()
		})

		It("panics when the type is not valid", func() {
			Expect(func() {
				HandlerType("<invalid>").MustValidate()
			}).To(Panic())
		})
	})

	Describe("func Is()", func() {
		It("returns true when the type is in the given set", func() {
			Expect(AggregateHandlerType.Is(AggregateHandlerType, ProcessHandlerType)).To(BeTrue())
		})

		It("returns false when the type is not in the given set", func() {
			Expect(IntegrationHandlerType.Is(AggregateHandlerType, ProcessHandlerType)).To(BeFalse())
		})
	})

	Describe("func MustBe()", func() {
		It("does not panic when the type is in the given set", func() {
			AggregateHandlerType.MustBe(AggregateHandlerType, ProcessHandlerType)
		})

		It("panics when the type is not in the given set", func() {
			Expect(func() {
				IntegrationHandlerType.MustBe(AggregateHandlerType, ProcessHandlerType)
			}).To(Panic())
		})
	})

	Describe("func MustNotBe()", func() {
		It("does not panic when the type is not in the given set", func() {
			IntegrationHandlerType.MustNotBe(AggregateHandlerType, ProcessHandlerType)
		})

		It("panics when the type is in the given set", func() {
			Expect(func() {
				AggregateHandlerType.MustNotBe(AggregateHandlerType, ProcessHandlerType)
			}).To(Panic())
		})
	})

	Describe("func IsConsumerOf()", func() {
		It("returns the expected values", func() {
			Expect(AggregateHandlerType.IsConsumerOf(message.CommandKind)).To(BeTrue())
			Expect(AggregateHandlerType.IsConsumerOf(message.EventKind)).To(BeFalse())
			Expect(AggregateHandlerType.IsConsumerOf(message.TimeoutKind)).To(BeFalse())

			Expect(ProcessHandlerType.IsConsumerOf(message.CommandKind)).To(BeFalse())
			Expect(ProcessHandlerType.IsConsumerOf(message.EventKind)).To(BeTrue())
			Expect(ProcessHandlerType.IsConsumerOf(message.TimeoutKind)).To(BeTrue())

			Expect(IntegrationHandlerType.IsConsumerOf(message.CommandKind)).To(BeTrue())
			Expect(IntegrationHandlerType.IsConsumerOf(message.EventKind)).To(BeFalse())
			Expect(IntegrationHandlerType.IsConsumerOf(message.TimeoutKind)).To(BeFalse())

			Expect(ProjectionHandlerType.IsConsumerOf(message.CommandKind)).To(BeFalse())
			Expect(ProjectionHandlerType.IsConsumerOf(message.EventKind)).To(BeTrue())
			Expect(ProjectionHandlerType.IsConsumerOf(message.TimeoutKind)).To(BeFalse())
		})
	})

	Describe("func IsProducerOf()", func() {
		It("returns the expected values", func() {
			Expect(AggregateHandlerType.IsProducerOf(message.CommandKind)).To(BeFalse())
			Expect(AggregateHandlerType.IsProducerOf(message.EventKind)).To(BeTrue())
			Expect(AggregateHandlerType.IsProducerOf(message.TimeoutKind)).To(BeFalse())

			Expect(ProcessHandlerType.IsProducerOf(message.CommandKind)).To(BeTrue())
			Expect(ProcessHandlerType.IsProducerOf(message.EventKind)).To(BeFalse())
			Expect(ProcessHandlerType.IsProducerOf(message.TimeoutKind)).To(BeTrue())

			Expect(IntegrationHandlerType.IsProducerOf(message.CommandKind)).To(BeFalse())
			Expect(IntegrationHandlerType.IsProducerOf(message.EventKind)).To(BeTrue())
			Expect(IntegrationHandlerType.IsProducerOf(message.TimeoutKind)).To(BeFalse())

			Expect(ProjectionHandlerType.IsProducerOf(message.CommandKind)).To(BeFalse())
			Expect(ProjectionHandlerType.IsProducerOf(message.EventKind)).To(BeFalse())
			Expect(ProjectionHandlerType.IsProducerOf(message.TimeoutKind)).To(BeFalse())
		})
	})

	Describe("func Consumes()", func() {
		It("returns the expected values", func() {
			Expect(AggregateHandlerType.Consumes()).To(ConsistOf(
				message.CommandKind,
			))

			Expect(ProcessHandlerType.Consumes()).To(ConsistOf(
				message.EventKind,
				message.TimeoutKind,
			))

			Expect(IntegrationHandlerType.Consumes()).To(ConsistOf(
				message.CommandKind,
			))

			Expect(ProjectionHandlerType.Consumes()).To(ConsistOf(
				message.EventKind,
			))
		})
	})

	Describe("func Produces()", func() {
		It("returns the expected values", func() {
			Expect(AggregateHandlerType.Produces()).To(ConsistOf(
				message.EventKind,
			))

			Expect(ProcessHandlerType.Produces()).To(ConsistOf(
				message.CommandKind,
				message.TimeoutKind,
			))

			Expect(IntegrationHandlerType.Produces()).To(ConsistOf(
				message.EventKind,
			))

			Expect(ProjectionHandlerType.Produces()).To(BeEmpty())
		})
	})

	Describe("func ShortString()", func() {
		It("returns the type value as a short string", func() {
			Expect(AggregateHandlerType.ShortString()).To(Equal("AGG"))
			Expect(ProcessHandlerType.ShortString()).To(Equal("PRC"))
			Expect(IntegrationHandlerType.ShortString()).To(Equal("INT"))
			Expect(ProjectionHandlerType.ShortString()).To(Equal("PRJ"))
		})
	})

	Describe("func String()", func() {
		It("returns the type value as a string", func() {
			Expect(AggregateHandlerType.String()).To(Equal("aggregate"))
			Expect(ProcessHandlerType.String()).To(Equal("process"))
			Expect(IntegrationHandlerType.String()).To(Equal("integration"))
			Expect(ProjectionHandlerType.String()).To(Equal("projection"))
			Expect(HandlerType("<invalid>").String()).To(Equal("<invalid handler type: <invalid>>"))
		})
	})

	Describe("func MarshalText()", func() {
		It("marshals the type to text", func() {
			Expect(AggregateHandlerType.MarshalText()).To(Equal([]byte("aggregate")))
			Expect(ProcessHandlerType.MarshalText()).To(Equal([]byte("process")))
			Expect(IntegrationHandlerType.MarshalText()).To(Equal([]byte("integration")))
			Expect(ProjectionHandlerType.MarshalText()).To(Equal([]byte("projection")))
		})

		It("returns an error if the type is invalid", func() {
			_, err := HandlerType("<invalid>").MarshalText()
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func UnmarshalText()", func() {
		It("unmarshals the type from text", func() {
			var t HandlerType

			err := t.UnmarshalText([]byte("aggregate"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(t).To(Equal(AggregateHandlerType))

			err = t.UnmarshalText([]byte("process"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(t).To(Equal(ProcessHandlerType))

			err = t.UnmarshalText([]byte("integration"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(t).To(Equal(IntegrationHandlerType))

			err = t.UnmarshalText([]byte("projection"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(t).To(Equal(ProjectionHandlerType))
		})

		It("returns an error if the data is invalid", func() {
			var t HandlerType

			err := t.UnmarshalText([]byte("<invalid>"))
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func MarshalBinary()", func() {
		It("marshals the type to binary", func() {
			Expect(AggregateHandlerType.MarshalBinary()).To(Equal([]byte("A")))
			Expect(ProcessHandlerType.MarshalBinary()).To(Equal([]byte("P")))
			Expect(IntegrationHandlerType.MarshalBinary()).To(Equal([]byte("I")))
			Expect(ProjectionHandlerType.MarshalBinary()).To(Equal([]byte("R")))
		})

		It("returns an error if the type is invalid", func() {
			_, err := HandlerType("<invalid>").MarshalBinary()
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func UnmarshalBinary()", func() {
		It("unmarshals the type from binary", func() {
			var t HandlerType

			err := t.UnmarshalBinary([]byte("A"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(t).To(Equal(AggregateHandlerType))

			err = t.UnmarshalBinary([]byte("P"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(t).To(Equal(ProcessHandlerType))

			err = t.UnmarshalBinary([]byte("I"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(t).To(Equal(IntegrationHandlerType))

			err = t.UnmarshalBinary([]byte("R"))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(t).To(Equal(ProjectionHandlerType))
		})

		It("returns an error if the data is the wrong length", func() {
			var t HandlerType

			err := t.UnmarshalBinary([]byte("<invalid>"))
			Expect(err).Should(HaveOccurred())
		})

		It("returns an error if the data does not contain a valid type", func() {
			var t HandlerType

			err := t.UnmarshalBinary([]byte("X"))
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("func ConsumersOf()", func() {
		It("returns the expected values", func() {
			Expect(ConsumersOf(message.CommandKind)).To(ConsistOf(
				AggregateHandlerType,
				IntegrationHandlerType,
			))

			Expect(ConsumersOf(message.EventKind)).To(ConsistOf(
				ProcessHandlerType,
				ProjectionHandlerType,
			))

			Expect(ConsumersOf(message.TimeoutKind)).To(ConsistOf(
				ProcessHandlerType,
			))
		})
	})

	Describe("func ProducersOf()", func() {
		It("returns the expected values", func() {
			Expect(ProducersOf(message.CommandKind)).To(ConsistOf(
				ProcessHandlerType,
			))

			Expect(ProducersOf(message.EventKind)).To(ConsistOf(
				AggregateHandlerType,
				IntegrationHandlerType,
			))

			Expect(ProducersOf(message.TimeoutKind)).To(ConsistOf(
				ProcessHandlerType,
			))
		})
	})
})
