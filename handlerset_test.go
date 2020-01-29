package configkit_test

import (
	"context"
	"errors"

	. "github.com/dogmatiq/configkit"
	cfixtures "github.com/dogmatiq/configkit/fixtures"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures" // can't dot-import due to conflicts
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("type HandlerSet", func() {
	var set *HandlerSet

	aggregate := FromAggregate(&fixtures.AggregateMessageHandler{
		ConfigureFunc: func(c dogma.AggregateConfigurer) {
			c.Identity("<agg-name>", "<agg-key>")
			c.ConsumesCommandType(fixtures.MessageC{})
			c.ProducesEventType(fixtures.MessageE{})
		},
	})

	projection := FromProjection(&fixtures.ProjectionMessageHandler{
		ConfigureFunc: func(c dogma.ProjectionConfigurer) {
			c.Identity("<proj-name>", "<proj-key>")
			c.ConsumesEventType(fixtures.MessageE{})
		},
	})

	BeforeEach(func() {
		set = &HandlerSet{}
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

			h, ok := set.ByIdentity(MustNewIdentity("<agg-name>", "<agg-key>"))
			Expect(ok).To(BeTrue())
			Expect(h).To(Equal(aggregate))
		})

		It("returns false if no such handler is in the set", func() {
			_, ok := set.ByIdentity(MustNewIdentity("<agg-name>", "<agg-key>"))
			Expect(ok).To(BeFalse())
		})
	})

	Describe("func ByName()", func() {
		It("returns the handler with the given name", func() {
			set.Add(aggregate)

			h, ok := set.ByName("<agg-name>")
			Expect(ok).To(BeTrue())
			Expect(h).To(Equal(aggregate))
		})

		It("returns false if no such handler is in the set", func() {
			_, ok := set.ByName("<agg-name>")
			Expect(ok).To(BeFalse())
		})
	})

	Describe("func ByKey()", func() {
		It("returns the handler with the given key", func() {
			set.Add(aggregate)

			h, ok := set.ByKey("<agg-key>")
			Expect(ok).To(BeTrue())
			Expect(h).To(Equal(aggregate))
		})

		It("returns false if no such handler is in the set", func() {
			_, ok := set.ByKey("<agg-key>")
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

	Describe("func IsEqual()", func() {
		DescribeTable(
			"returns true if the sets are equivalent",
			func(a, b HandlerSet) {
				Expect(a.IsEqual(b)).To(BeTrue())
			},
			Entry(
				"equivalent",
				NewHandlerSet(aggregate, projection),
				NewHandlerSet(aggregate, projection),
			),
			Entry(
				"nil and empty",
				HandlerSet{},
				HandlerSet(nil),
			),
		)

		DescribeTable(
			"returns false if the sets are not equivalent",
			func(b HandlerSet) {
				a := NewHandlerSet(aggregate, projection)

				Expect(a.IsEqual(b)).To(BeFalse())
			},
			Entry(
				"subset",
				NewHandlerSet(aggregate),
			),
			Entry(
				"superset",
				NewHandlerSet(aggregate, projection, FromIntegration(&fixtures.IntegrationMessageHandler{
					ConfigureFunc: func(c dogma.IntegrationConfigurer) {
						c.Identity("<int-name>", "<int-key>")
						c.ConsumesCommandType(fixtures.MessageX{})
					},
				})),
			),
			Entry(
				"same-length, disjoint handler",
				NewHandlerSet(aggregate, FromProjection(&fixtures.ProjectionMessageHandler{
					ConfigureFunc: func(c dogma.ProjectionConfigurer) {
						c.Identity("<proj-name>", "<proj-key>")
						c.ConsumesEventType(fixtures.MessageF{}) // diff
					},
				})),
			),
		)
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

	Describe("func AcceptVisitor()", func() {
		BeforeEach(func() {
			set.Add(aggregate)
			set.Add(projection)
		})

		It("visits each handler in the set", func() {
			var visited []Handler

			err := set.AcceptVisitor(
				context.Background(),
				&cfixtures.Visitor{
					VisitAggregateFunc: func(_ context.Context, cfg Aggregate) error {
						Expect(cfg).To(BeIdenticalTo(aggregate))
						visited = append(visited, cfg)
						return nil
					},
					VisitProjectionFunc: func(_ context.Context, cfg Projection) error {
						Expect(cfg).To(BeIdenticalTo(projection))
						visited = append(visited, cfg)
						return nil
					},
				},
			)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(visited).To(ConsistOf(aggregate, projection))
		})

		It("returns an error if one of the handlers fails", func() {
			err := set.AcceptVisitor(
				context.Background(),
				&cfixtures.Visitor{
					VisitProjectionFunc: func(_ context.Context, cfg Projection) error {
						return errors.New("<error>")
					},
				},
			)

			Expect(err).To(MatchError("<error>"))
		})
	})

	Context("each functions", func() {
		var (
			aggregate1, aggregate2     Aggregate
			process1, process2         Process
			integration1, integration2 Integration
			projection1, projection2   Projection
		)

		BeforeEach(func() {
			aggregate1 = FromAggregate(&fixtures.AggregateMessageHandler{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("<agg1-name>", "<agg1-key>")
					c.ConsumesCommandType(fixtures.MessageC{})
					c.ProducesEventType(fixtures.MessageD{})
				},
			})

			aggregate2 = FromAggregate(&fixtures.AggregateMessageHandler{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("<agg2-name>", "<agg2-key>")
					c.ConsumesCommandType(fixtures.MessageC{})
					c.ProducesEventType(fixtures.MessageD{})
				},
			})

			process1 = FromProcess(&fixtures.ProcessMessageHandler{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("<proc1-name>", "<proc1-key>")
					c.ConsumesEventType(fixtures.MessageE{})
					c.ProducesCommandType(fixtures.MessageC{})
				},
			})

			process2 = FromProcess(&fixtures.ProcessMessageHandler{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("<proc2-name>", "<proc2-key>")
					c.ConsumesEventType(fixtures.MessageE{})
					c.ProducesCommandType(fixtures.MessageC{})
				},
			})

			integration1 = FromIntegration(&fixtures.IntegrationMessageHandler{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("<int1-name>", "<int1-key>")
					c.ConsumesCommandType(fixtures.MessageC{})
					c.ProducesEventType(fixtures.MessageD{})
				},
			})

			integration2 = FromIntegration(&fixtures.IntegrationMessageHandler{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("<int2-name>", "<int2-key>")
					c.ConsumesCommandType(fixtures.MessageC{})
					c.ProducesEventType(fixtures.MessageD{})
				},
			})

			projection1 = FromProjection(&fixtures.ProjectionMessageHandler{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("<proj1-name>", "<proj1-key>")
					c.ConsumesEventType(fixtures.MessageE{})
				},
			})

			projection2 = FromProjection(&fixtures.ProjectionMessageHandler{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("<proj2-name>", "<proj2-key>")
					c.ConsumesEventType(fixtures.MessageE{})
				},
			})

			set.Add(aggregate1)
			set.Add(aggregate2)
			set.Add(process1)
			set.Add(process2)
			set.Add(integration1)
			set.Add(integration2)
			set.Add(projection1)
			set.Add(projection2)
		})

		Describe("func RangeAggregates()", func() {
			It("calls fn for each aggregate in the set", func() {
				var names []string

				all := set.RangeAggregates(func(h Aggregate) bool {
					names = append(names, h.Identity().Name)
					return true
				})

				Expect(names).To(ConsistOf("<agg1-name>", "<agg2-name>"))
				Expect(all).To(BeTrue())
			})

			It("stops iterating if fn returns false", func() {
				count := 0

				all := set.RangeAggregates(func(h Aggregate) bool {
					count++
					return false
				})

				Expect(count).To(BeNumerically("==", 1))
				Expect(all).To(BeFalse())
			})
		})

		Describe("func RangeProcesses()", func() {
			It("calls fn for each process in the set", func() {
				var names []string

				all := set.RangeProcesses(func(h Process) bool {
					names = append(names, h.Identity().Name)
					return true
				})

				Expect(names).To(ConsistOf("<proc1-name>", "<proc2-name>"))
				Expect(all).To(BeTrue())
			})

			It("stops iterating if fn returns false", func() {
				count := 0

				all := set.RangeProcesses(func(h Process) bool {
					count++
					return false
				})

				Expect(count).To(BeNumerically("==", 1))
				Expect(all).To(BeFalse())
			})
		})

		Describe("func RangeIntegrations()", func() {
			It("calls fn for each integration in the set", func() {
				var names []string

				all := set.RangeIntegrations(func(h Integration) bool {
					names = append(names, h.Identity().Name)
					return true
				})

				Expect(names).To(ConsistOf("<int1-name>", "<int2-name>"))
				Expect(all).To(BeTrue())
			})

			It("stops iterating if fn returns false", func() {
				count := 0

				all := set.RangeIntegrations(func(h Integration) bool {
					count++
					return false
				})

				Expect(count).To(BeNumerically("==", 1))
				Expect(all).To(BeFalse())
			})
		})

		Describe("func RangeProjections()", func() {
			It("calls fn for each projection in the set", func() {
				var names []string

				all := set.RangeProjections(func(h Projection) bool {
					names = append(names, h.Identity().Name)
					return true
				})

				Expect(names).To(ConsistOf("<proj1-name>", "<proj2-name>"))
				Expect(all).To(BeTrue())
			})

			It("stops iterating if fn returns false", func() {
				count := 0

				all := set.RangeProjections(func(h Projection) bool {
					count++
					return false
				})

				Expect(count).To(BeNumerically("==", 1))
				Expect(all).To(BeFalse())
			})
		})
	})
})

var _ = Describe("type RichHandlerSet", func() {
	var set *RichHandlerSet

	aggregate := FromAggregate(&fixtures.AggregateMessageHandler{
		ConfigureFunc: func(c dogma.AggregateConfigurer) {
			c.Identity("<agg-name>", "<agg-key>")
			c.ConsumesCommandType(fixtures.MessageC{})
			c.ProducesEventType(fixtures.MessageE{})
		},
	})

	projection := FromProjection(&fixtures.ProjectionMessageHandler{
		ConfigureFunc: func(c dogma.ProjectionConfigurer) {
			c.Identity("<proj-name>", "<proj-key>")
			c.ConsumesEventType(fixtures.MessageE{})
		},
	})

	BeforeEach(func() {
		set = &RichHandlerSet{}
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

			h, ok := set.ByIdentity(MustNewIdentity("<agg-name>", "<agg-key>"))
			Expect(ok).To(BeTrue())
			Expect(h).To(Equal(aggregate))
		})

		It("returns false if no such handler is in the set", func() {
			_, ok := set.ByIdentity(MustNewIdentity("<agg-name>", "<agg-key>"))
			Expect(ok).To(BeFalse())
		})
	})

	Describe("func ByName()", func() {
		It("returns the handler with the given name", func() {
			set.Add(aggregate)

			h, ok := set.ByName("<agg-name>")
			Expect(ok).To(BeTrue())
			Expect(h).To(Equal(aggregate))
		})

		It("returns false if no such handler is in the set", func() {
			_, ok := set.ByName("<agg-name>")
			Expect(ok).To(BeFalse())
		})
	})

	Describe("func ByKey()", func() {
		It("returns the handler with the given key", func() {
			set.Add(aggregate)

			h, ok := set.ByKey("<agg-key>")
			Expect(ok).To(BeTrue())
			Expect(h).To(Equal(aggregate))
		})

		It("returns false if no such handler is in the set", func() {
			_, ok := set.ByKey("<agg-key>")
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

	Describe("func IsEqual()", func() {
		DescribeTable(
			"returns true if the sets are equivalent",
			func(a, b RichHandlerSet) {
				Expect(a.IsEqual(b)).To(BeTrue())
			},
			Entry(
				"equivalent",
				NewRichHandlerSet(aggregate, projection),
				NewRichHandlerSet(aggregate, projection),
			),
			Entry(
				"nil and empty",
				RichHandlerSet{},
				RichHandlerSet(nil),
			),
		)

		DescribeTable(
			"returns false if the sets are not equivalent",
			func(b RichHandlerSet) {
				a := NewRichHandlerSet(aggregate, projection)

				Expect(a.IsEqual(b)).To(BeFalse())
			},
			Entry(
				"subset",
				NewRichHandlerSet(aggregate),
			),
			Entry(
				"superset",
				NewRichHandlerSet(aggregate, projection, FromIntegration(&fixtures.IntegrationMessageHandler{
					ConfigureFunc: func(c dogma.IntegrationConfigurer) {
						c.Identity("<int-name>", "<int-key>")
						c.ConsumesCommandType(fixtures.MessageX{})
					},
				})),
			),
			Entry(
				"same-length, disjoint handler",
				NewRichHandlerSet(aggregate, FromProjection(&fixtures.ProjectionMessageHandler{
					ConfigureFunc: func(c dogma.ProjectionConfigurer) {
						c.Identity("<proj-name>", "<proj-key>")
						c.ConsumesEventType(fixtures.MessageF{}) // diff
					},
				})),
			),
		)
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

	Describe("func AcceptRichVisitor()", func() {
		BeforeEach(func() {
			set.Add(aggregate)
			set.Add(projection)
		})

		It("visits each handler in the set", func() {
			var visited []RichHandler

			err := set.AcceptRichVisitor(
				context.Background(),
				&cfixtures.RichVisitor{
					VisitRichAggregateFunc: func(_ context.Context, cfg RichAggregate) error {
						Expect(cfg).To(BeIdenticalTo(aggregate))
						visited = append(visited, cfg)
						return nil
					},
					VisitRichProjectionFunc: func(_ context.Context, cfg RichProjection) error {
						Expect(cfg).To(BeIdenticalTo(projection))
						visited = append(visited, cfg)
						return nil
					},
				},
			)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(visited).To(ConsistOf(aggregate, projection))
		})

		It("returns an error if one of the handlers fails", func() {
			err := set.AcceptRichVisitor(
				context.Background(),
				&cfixtures.RichVisitor{
					VisitRichProjectionFunc: func(_ context.Context, cfg RichProjection) error {
						return errors.New("<error>")
					},
				},
			)

			Expect(err).To(MatchError("<error>"))
		})
	})
})
