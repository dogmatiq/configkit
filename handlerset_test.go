package configkit_test

import (
	. "github.com/dogmatiq/configkit"
	cfixtures "github.com/dogmatiq/configkit/fixtures"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures" // can't dot-import due to conflicts
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type HandlerSet", func() {
	var (
		set        *HandlerSet
		aggregate  RichAggregate
		projection RichProjection
	)

	BeforeEach(func() {
		set = &HandlerSet{}

		aggregate = FromAggregate(&fixtures.AggregateMessageHandler{
			ConfigureFunc: func(c dogma.AggregateConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesCommandType(fixtures.MessageC{})
				c.ProducesEventType(fixtures.MessageE{})
			},
		})

		projection = FromProjection(&fixtures.ProjectionMessageHandler{
			ConfigureFunc: func(c dogma.ProjectionConfigurer) {
				c.Identity("<proj-name>", "<proj-key>")
				c.ConsumesEventType(fixtures.MessageE{})
			},
		})
	})

	Describe("func NewHandlerSet()", func() {
		It("returns a set containing the given handlers", func() {
			s := NewHandlerSet(aggregate, projection)
			Expect(s).To(HaveLen(2))
			Expect(s.Has(aggregate)).To(BeTrue())
			Expect(s.Has(projection)).To(BeTrue())
		})

		It("panics if the handler identities conflict", func() {
			Expect(func() {
				NewHandlerSet(
					aggregate,
					FromAggregate(aggregate.Handler()),
				)
			}).To(Panic())
		})
	})

	Describe("func Add()", func() {
		It("adds the handler to the set", func() {
			ok := set.Add(aggregate)
			Expect(ok).To(BeTrue())
			Expect(set.Has(aggregate)).To(BeTrue())
		})

		It("does not add the handler if the identity conflicts with another handler", func() {
			another := FromAggregate(aggregate.Handler())
			set.Add(aggregate)

			ok := set.Add(another)
			Expect(ok).To(BeFalse())
			Expect(set.Has(another)).To(BeFalse())
		})
	})

	Describe("func Has()", func() {
		It("returns true if the handler is in the set", func() {
			set.Add(aggregate)
			Expect(set.Has(aggregate)).To(BeTrue())
		})

		It("returns false if the handler is not the set", func() {
			Expect(set.Has(aggregate)).To(BeFalse())
		})

		It("returns false if the a different handler with the same identity is in the set", func() {
			set.Add(aggregate)

			another := FromAggregate(aggregate.Handler())
			Expect(set.Has(another)).To(BeFalse())
		})
	})

	Describe("func ByIdentity()", func() {
		It("returns the handler with the given identity", func() {
			set.Add(aggregate)

			h, ok := set.ByIdentity(MustNewIdentity("<name>", "<key>"))
			Expect(ok).To(BeTrue())
			Expect(h).To(Equal(aggregate))
		})

		It("returns false if no such handler is in the set", func() {
			_, ok := set.ByIdentity(MustNewIdentity("<name>", "<key>"))
			Expect(ok).To(BeFalse())
		})
	})

	Describe("func ByName()", func() {
		It("returns the handler with the given name", func() {
			set.Add(aggregate)

			h, ok := set.ByName("<name>")
			Expect(ok).To(BeTrue())
			Expect(h).To(Equal(aggregate))
		})

		It("returns false if no such handler is in the set", func() {
			_, ok := set.ByName("<name>")
			Expect(ok).To(BeFalse())
		})
	})

	Describe("func ByKey()", func() {
		It("returns the handler with the given key", func() {
			set.Add(aggregate)

			h, ok := set.ByKey("<key>")
			Expect(ok).To(BeTrue())
			Expect(h).To(Equal(aggregate))
		})

		It("returns false if no such handler is in the set", func() {
			_, ok := set.ByKey("<key>")
			Expect(ok).To(BeFalse())
		})
	})

	Describe("func ByType()", func() {
		BeforeEach(func() {
			set.Add(aggregate)
			set.Add(projection)
		})

		It("returns the handlers of the given type", func() {
			subset := set.ByType(AggregateHandlerType)
			Expect(subset).To(HaveLen(1))
			Expect(set.Has(aggregate)).To(BeTrue())
		})
	})

	Describe("func ConsumersOf()", func() {
		BeforeEach(func() {
			set.Add(aggregate)
			set.Add(projection)
		})

		It("returns the handlers that consume the given message", func() {
			subset := set.ConsumersOf(cfixtures.MessageCTypeName)
			Expect(subset).To(HaveLen(1))
			Expect(set.Has(aggregate)).To(BeTrue())
		})
	})

	Describe("func ProducersOf()", func() {
		BeforeEach(func() {
			set.Add(aggregate)
			set.Add(projection)
		})

		It("returns the handlers the produce the given message", func() {
			subset := set.ProducersOf(cfixtures.MessageETypeName)
			Expect(subset).To(HaveLen(1))
			Expect(set.Has(aggregate)).To(BeTrue())
		})
	})

	Describe("func Find()", func() {
		BeforeEach(func() {
			set.Add(aggregate)
			set.Add(projection)
		})

		It("returns a handler that matches the predicate", func() {
			h, ok := set.Find(func(h Handler) bool {
				return h.Identity() == projection.Identity()
			})
			Expect(ok).To(BeTrue())
			Expect(h).To(Equal(projection))
		})

		It("returns false if no handler matches", func() {
			_, ok := set.Find(func(Handler) bool {
				return false
			})
			Expect(ok).To(BeFalse())
		})
	})

	Describe("func Filter()", func() {
		BeforeEach(func() {
			set.Add(aggregate)
			set.Add(projection)
		})

		It("returns handlers that match the predicate", func() {
			subset := set.Filter(func(h Handler) bool {
				return h.Identity() == projection.Identity()
			})
			Expect(subset).To(HaveLen(1))
			Expect(set.Has(projection)).To(BeTrue())
		})
	})
})

var _ = Describe("type RichHandlerSet", func() {
	var (
		set        *RichHandlerSet
		aggregate  RichAggregate
		projection RichProjection
	)

	BeforeEach(func() {
		set = &RichHandlerSet{}

		aggregate = FromAggregate(&fixtures.AggregateMessageHandler{
			ConfigureFunc: func(c dogma.AggregateConfigurer) {
				c.Identity("<name>", "<key>")
				c.ConsumesCommandType(fixtures.MessageC{})
				c.ProducesEventType(fixtures.MessageE{})
			},
		})

		projection = FromProjection(&fixtures.ProjectionMessageHandler{
			ConfigureFunc: func(c dogma.ProjectionConfigurer) {
				c.Identity("<proj-name>", "<proj-key>")
				c.ConsumesEventType(fixtures.MessageE{})
			},
		})
	})

	Describe("func NewRichHandlerSet()", func() {
		It("returns a set containing the given handlers", func() {
			s := NewRichHandlerSet(aggregate, projection)
			Expect(s).To(HaveLen(2))
			Expect(s.Has(aggregate)).To(BeTrue())
			Expect(s.Has(projection)).To(BeTrue())
		})

		It("panics if the handler identities conflict", func() {
			Expect(func() {
				NewRichHandlerSet(
					aggregate,
					FromAggregate(aggregate.Handler()),
				)
			}).To(Panic())
		})
	})

	Describe("func Add()", func() {
		It("adds the handler to the set", func() {
			ok := set.Add(aggregate)
			Expect(ok).To(BeTrue())
			Expect(set.Has(aggregate)).To(BeTrue())
		})

		It("does not add the handler if the identity conflicts with another handler", func() {
			another := FromAggregate(aggregate.Handler())
			set.Add(aggregate)

			ok := set.Add(another)
			Expect(ok).To(BeFalse())
			Expect(set.Has(another)).To(BeFalse())
		})
	})

	Describe("func Has()", func() {
		It("returns true if the handler is in the set", func() {
			set.Add(aggregate)
			Expect(set.Has(aggregate)).To(BeTrue())
		})

		It("returns false if the handler is not the set", func() {
			Expect(set.Has(aggregate)).To(BeFalse())
		})

		It("returns false if the a different handler with the same identity is in the set", func() {
			set.Add(aggregate)

			another := FromAggregate(aggregate.Handler())
			Expect(set.Has(another)).To(BeFalse())
		})
	})

	Describe("func ByIdentity()", func() {
		It("returns the handler with the given identity", func() {
			set.Add(aggregate)

			h, ok := set.ByIdentity(MustNewIdentity("<name>", "<key>"))
			Expect(ok).To(BeTrue())
			Expect(h).To(Equal(aggregate))
		})

		It("returns false if no such handler is in the set", func() {
			_, ok := set.ByIdentity(MustNewIdentity("<name>", "<key>"))
			Expect(ok).To(BeFalse())
		})
	})

	Describe("func ByName()", func() {
		It("returns the handler with the given name", func() {
			set.Add(aggregate)

			h, ok := set.ByName("<name>")
			Expect(ok).To(BeTrue())
			Expect(h).To(Equal(aggregate))
		})

		It("returns false if no such handler is in the set", func() {
			_, ok := set.ByName("<name>")
			Expect(ok).To(BeFalse())
		})
	})

	Describe("func ByKey()", func() {
		It("returns the handler with the given key", func() {
			set.Add(aggregate)

			h, ok := set.ByKey("<key>")
			Expect(ok).To(BeTrue())
			Expect(h).To(Equal(aggregate))
		})

		It("returns false if no such handler is in the set", func() {
			_, ok := set.ByKey("<key>")
			Expect(ok).To(BeFalse())
		})
	})

	Describe("func ByType()", func() {
		BeforeEach(func() {
			set.Add(aggregate)
			set.Add(projection)
		})

		It("returns the handlers of the given type", func() {
			subset := set.ByType(AggregateHandlerType)
			Expect(subset).To(HaveLen(1))
			Expect(set.Has(aggregate)).To(BeTrue())
		})
	})

	Describe("func ConsumersOf()", func() {
		BeforeEach(func() {
			set.Add(aggregate)
			set.Add(projection)
		})

		It("returns the handlers that consume the given message", func() {
			subset := set.ConsumersOf(cfixtures.MessageCType)
			Expect(subset).To(HaveLen(1))
			Expect(set.Has(aggregate)).To(BeTrue())
		})
	})

	Describe("func ProducersOf()", func() {
		BeforeEach(func() {
			set.Add(aggregate)
			set.Add(projection)
		})

		It("returns the handlers the produce the given message", func() {
			subset := set.ProducersOf(cfixtures.MessageEType)
			Expect(subset).To(HaveLen(1))
			Expect(set.Has(aggregate)).To(BeTrue())
		})
	})

	Describe("func Find()", func() {
		BeforeEach(func() {
			set.Add(aggregate)
			set.Add(projection)
		})

		It("returns a handler that matches the predicate", func() {
			h, ok := set.Find(func(h RichHandler) bool {
				return h.Identity() == projection.Identity()
			})
			Expect(ok).To(BeTrue())
			Expect(h).To(Equal(projection))
		})

		It("returns false if no handler matches", func() {
			_, ok := set.Find(func(RichHandler) bool {
				return false
			})
			Expect(ok).To(BeFalse())
		})
	})

	Describe("func Filter()", func() {
		BeforeEach(func() {
			set.Add(aggregate)
			set.Add(projection)
		})

		It("returns handlers that match the predicate", func() {
			subset := set.Filter(func(h RichHandler) bool {
				return h.Identity() == projection.Identity()
			})
			Expect(subset).To(HaveLen(1))
			Expect(set.Has(projection)).To(BeTrue())
		})
	})
})
