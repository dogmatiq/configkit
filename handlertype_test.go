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
			Expect(AggregateHandlerType.IsConsumerOf(message.CommandRole)).To(BeTrue())
			Expect(AggregateHandlerType.IsConsumerOf(message.EventRole)).To(BeFalse())
			Expect(AggregateHandlerType.IsConsumerOf(message.TimeoutRole)).To(BeFalse())

			Expect(ProcessHandlerType.IsConsumerOf(message.CommandRole)).To(BeFalse())
			Expect(ProcessHandlerType.IsConsumerOf(message.EventRole)).To(BeTrue())
			Expect(ProcessHandlerType.IsConsumerOf(message.TimeoutRole)).To(BeTrue())

			Expect(IntegrationHandlerType.IsConsumerOf(message.CommandRole)).To(BeTrue())
			Expect(IntegrationHandlerType.IsConsumerOf(message.EventRole)).To(BeFalse())
			Expect(IntegrationHandlerType.IsConsumerOf(message.TimeoutRole)).To(BeFalse())

			Expect(ProjectionHandlerType.IsConsumerOf(message.CommandRole)).To(BeFalse())
			Expect(ProjectionHandlerType.IsConsumerOf(message.EventRole)).To(BeTrue())
			Expect(ProjectionHandlerType.IsConsumerOf(message.TimeoutRole)).To(BeFalse())
		})
	})

	Describe("func IsProducerOf()", func() {
		It("returns the expected values", func() {
			Expect(AggregateHandlerType.IsProducerOf(message.CommandRole)).To(BeFalse())
			Expect(AggregateHandlerType.IsProducerOf(message.EventRole)).To(BeTrue())
			Expect(AggregateHandlerType.IsProducerOf(message.TimeoutRole)).To(BeFalse())

			Expect(ProcessHandlerType.IsProducerOf(message.CommandRole)).To(BeTrue())
			Expect(ProcessHandlerType.IsProducerOf(message.EventRole)).To(BeFalse())
			Expect(ProcessHandlerType.IsProducerOf(message.TimeoutRole)).To(BeTrue())

			Expect(IntegrationHandlerType.IsProducerOf(message.CommandRole)).To(BeFalse())
			Expect(IntegrationHandlerType.IsProducerOf(message.EventRole)).To(BeTrue())
			Expect(IntegrationHandlerType.IsProducerOf(message.TimeoutRole)).To(BeFalse())

			Expect(ProjectionHandlerType.IsProducerOf(message.CommandRole)).To(BeFalse())
			Expect(ProjectionHandlerType.IsProducerOf(message.EventRole)).To(BeFalse())
			Expect(ProjectionHandlerType.IsProducerOf(message.TimeoutRole)).To(BeFalse())
		})
	})

	Describe("func Consumes()", func() {
		It("returns the expected values", func() {
			Expect(AggregateHandlerType.Consumes()).To(ConsistOf(
				message.CommandRole,
			))

			Expect(ProcessHandlerType.Consumes()).To(ConsistOf(
				message.EventRole,
				message.TimeoutRole,
			))

			Expect(IntegrationHandlerType.Consumes()).To(ConsistOf(
				message.CommandRole,
			))

			Expect(ProjectionHandlerType.Consumes()).To(ConsistOf(
				message.EventRole,
			))
		})
	})

	Describe("func Produces()", func() {
		It("returns the expected values", func() {
			Expect(AggregateHandlerType.Produces()).To(ConsistOf(
				message.EventRole,
			))

			Expect(ProcessHandlerType.Produces()).To(ConsistOf(
				message.CommandRole,
				message.TimeoutRole,
			))

			Expect(IntegrationHandlerType.Produces()).To(ConsistOf(
				message.EventRole,
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
			Expect(ConsumersOf(message.CommandRole)).To(ConsistOf(
				AggregateHandlerType,
				IntegrationHandlerType,
			))

			Expect(ConsumersOf(message.EventRole)).To(ConsistOf(
				ProcessHandlerType,
				ProjectionHandlerType,
			))

			Expect(ConsumersOf(message.TimeoutRole)).To(ConsistOf(
				ProcessHandlerType,
			))
		})
	})

	Describe("func ProducersOf()", func() {
		It("returns the expected values", func() {
			Expect(ProducersOf(message.CommandRole)).To(ConsistOf(
				ProcessHandlerType,
			))

			Expect(ProducersOf(message.EventRole)).To(ConsistOf(
				AggregateHandlerType,
				IntegrationHandlerType,
			))

			Expect(ProducersOf(message.TimeoutRole)).To(ConsistOf(
				ProcessHandlerType,
			))
		})
	})
})
