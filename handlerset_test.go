package configkit_test

import (
	"context"
	"errors"

	. "github.com/dogmatiq/configkit"
	cfixtures "github.com/dogmatiq/configkit/fixtures"
	"github.com/dogmatiq/configkit/message"
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
			c.Identity("<agg-name>", aggregateKey)
			c.Routes(
				dogma.HandlesCommand[fixtures.MessageC](),
				dogma.RecordsEvent[fixtures.MessageE](),
			)
		},
	})

	projection := FromProjection(&fixtures.ProjectionMessageHandler{
		ConfigureFunc: func(c dogma.ProjectionConfigurer) {
			c.Identity("<proj-name>", projectionKey)
			c.Routes(
				dogma.HandlesEvent[fixtures.MessageE](),
			)
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

			h, ok := set.ByIdentity(MustNewIdentity("<agg-name>", aggregateKey))
			Expect(ok).To(BeTrue())
			Expect(h).To(Equal(aggregate))
		})

		It("returns false if no such handler is in the set", func() {
			_, ok := set.ByIdentity(MustNewIdentity("<agg-name>", aggregateKey))
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

			h, ok := set.ByKey(aggregateKey)
			Expect(ok).To(BeTrue())
			Expect(h).To(Equal(aggregate))
		})

		It("returns false if no such handler is in the set", func() {
			_, ok := set.ByKey(aggregateKey)
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

	Describe("func MessageNames()", func() {
		BeforeEach(func() {
			set.Add(aggregate)
			set.Add(projection)
		})

		It("returns the messages used by the handlers in the set", func() {
			Expect(set.MessageNames()).To(Equal(
				EntityMessageNames{
					Produced: message.NameRoles{
						cfixtures.MessageETypeName: message.EventRole,
					},
					Consumed: message.NameRoles{
						cfixtures.MessageCTypeName: message.CommandRole,
						cfixtures.MessageETypeName: message.EventRole,
					},
				},
			))
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
						c.Identity("<int-name>", integrationKey)
						c.Routes(
							dogma.HandlesCommand[fixtures.MessageX](),
						)
					},
				})),
			),
			Entry(
				"same-length, disjoint handler",
				NewHandlerSet(aggregate, FromProjection(&fixtures.ProjectionMessageHandler{
					ConfigureFunc: func(c dogma.ProjectionConfigurer) {
						c.Identity("<proj-name>", projectionKey)
						c.Routes(
							dogma.HandlesEvent[fixtures.MessageF](), // diff
						)
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

	Context("type-specific filtering", func() {
		var (
			aggregate1, aggregate2     Aggregate
			process1, process2         Process
			integration1, integration2 Integration
			projection1, projection2   Projection
		)

		BeforeEach(func() {
			aggregate1 = FromAggregate(&fixtures.AggregateMessageHandler{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("<agg1-name>", "ca82de11-1794-486e-a190-78e2443de7dd")
					c.Routes(
						dogma.HandlesCommand[fixtures.MessageC](),
						dogma.RecordsEvent[fixtures.MessageD](),
					)
				},
			})

			aggregate2 = FromAggregate(&fixtures.AggregateMessageHandler{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("<agg2-name>", "3a7ad56e-d9b9-42be-9a16-01f25e572c49")
					c.Routes(
						dogma.HandlesCommand[fixtures.MessageC](),
						dogma.RecordsEvent[fixtures.MessageD](),
					)
				},
			})

			process1 = FromProcess(&fixtures.ProcessMessageHandler{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("<proc1-name>", "5695e728-33ca-4b0f-b063-ff0ff6f48276")
					c.Routes(
						dogma.HandlesEvent[fixtures.MessageE](),
						dogma.ExecutesCommand[fixtures.MessageC](),
					)
				},
			})

			process2 = FromProcess(&fixtures.ProcessMessageHandler{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("<proc2-name>", "3d510ba8-6dca-46bb-bcde-193015867834")
					c.Routes(
						dogma.HandlesEvent[fixtures.MessageE](),
						dogma.ExecutesCommand[fixtures.MessageC](),
					)
				},
			})

			integration1 = FromIntegration(&fixtures.IntegrationMessageHandler{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("<int1-name>", "fdf0059e-8786-42db-a348-caac60d6118a")
					c.Routes(
						dogma.HandlesCommand[fixtures.MessageC](),
						dogma.RecordsEvent[fixtures.MessageD](),
					)
				},
			})

			integration2 = FromIntegration(&fixtures.IntegrationMessageHandler{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("<int2-name>", "34d19336-95c5-47c7-b36e-4e90b24b1b83")
					c.Routes(
						dogma.HandlesCommand[fixtures.MessageC](),
						dogma.RecordsEvent[fixtures.MessageD](),
					)
				},
			})

			projection1 = FromProjection(&fixtures.ProjectionMessageHandler{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("<proj1-name>", "ee9bd355-e9fd-413a-ac83-7182ea76cb89")
					c.Routes(
						dogma.HandlesEvent[fixtures.MessageE](),
					)
				},
			})

			projection2 = FromProjection(&fixtures.ProjectionMessageHandler{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("<proj2-name>", "0fe3a5e3-b8e1-4ba0-8c90-f002f0a842f9")
					c.Routes(
						dogma.HandlesEvent[fixtures.MessageE](),
					)
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

		Describe("func Aggregates()", func() {
			It("returns a slice containing the aggregates", func() {
				Expect(set.Aggregates()).To(ConsistOf(aggregate1, aggregate2))
			})
		})

		Describe("func Processes()", func() {
			It("returns a slice containing the processes", func() {
				Expect(set.Processes()).To(ConsistOf(process1, process2))
			})
		})

		Describe("func Integrations()", func() {
			It("returns a slice containing the integrations", func() {
				Expect(set.Integrations()).To(ConsistOf(integration1, integration2))
			})
		})

		Describe("func Projections()", func() {
			It("returns a slice containing the projections", func() {
				Expect(set.Projections()).To(ConsistOf(projection1, projection2))
			})
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
			c.Identity("<agg-name>", aggregateKey)
			c.Routes(
				dogma.HandlesCommand[fixtures.MessageC](),
				dogma.RecordsEvent[fixtures.MessageE](),
			)
		},
	})

	projection := FromProjection(&fixtures.ProjectionMessageHandler{
		ConfigureFunc: func(c dogma.ProjectionConfigurer) {
			c.Identity("<proj-name>", projectionKey)
			c.Routes(
				dogma.HandlesEvent[fixtures.MessageE](),
			)
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

			h, ok := set.ByIdentity(MustNewIdentity("<agg-name>", aggregateKey))
			Expect(ok).To(BeTrue())
			Expect(h).To(Equal(aggregate))
		})

		It("returns false if no such handler is in the set", func() {
			_, ok := set.ByIdentity(MustNewIdentity("<agg-name>", aggregateKey))
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

			h, ok := set.ByKey(aggregateKey)
			Expect(ok).To(BeTrue())
			Expect(h).To(Equal(aggregate))
		})

		It("returns false if no such handler is in the set", func() {
			_, ok := set.ByKey(aggregateKey)
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

	Describe("func MessageTypes()", func() {
		BeforeEach(func() {
			set.Add(aggregate)
			set.Add(projection)
		})

		It("returns the messages used by the handlers in the set", func() {
			Expect(set.MessageTypes()).To(Equal(
				EntityMessageTypes{
					Produced: message.TypeRoles{
						cfixtures.MessageEType: message.EventRole,
					},
					Consumed: message.TypeRoles{
						cfixtures.MessageCType: message.CommandRole,
						cfixtures.MessageEType: message.EventRole,
					},
				},
			))
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
						c.Identity("<int-name>", integrationKey)
						c.Routes(
							dogma.HandlesCommand[fixtures.MessageX](),
						)
					},
				})),
			),
			Entry(
				"same-length, disjoint handler",
				NewRichHandlerSet(aggregate, FromProjection(&fixtures.ProjectionMessageHandler{
					ConfigureFunc: func(c dogma.ProjectionConfigurer) {
						c.Identity("<proj-name>", projectionKey)
						c.Routes(
							dogma.HandlesEvent[fixtures.MessageF](), // diff
						)
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

	Context("type-specific filtering", func() {
		var (
			aggregate1, aggregate2     RichAggregate
			process1, process2         RichProcess
			integration1, integration2 RichIntegration
			projection1, projection2   RichProjection
		)

		BeforeEach(func() {
			aggregate1 = FromAggregate(&fixtures.AggregateMessageHandler{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("<agg1-name>", "648c035d-2a6a-49e6-8968-044bec062fed")
					c.Routes(
						dogma.HandlesCommand[fixtures.MessageC](),
						dogma.RecordsEvent[fixtures.MessageD](),
					)
				},
			})

			aggregate2 = FromAggregate(&fixtures.AggregateMessageHandler{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("<agg2-name>", "e465c85d-4ac0-4aed-8054-665a86b9ef4e")
					c.Routes(
						dogma.HandlesCommand[fixtures.MessageC](),
						dogma.RecordsEvent[fixtures.MessageD](),
					)
				},
			})

			process1 = FromProcess(&fixtures.ProcessMessageHandler{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("<proc1-name>", "71a4111b-ee0d-4df1-a059-d8bb94dc3e77")
					c.Routes(
						dogma.HandlesEvent[fixtures.MessageE](),
						dogma.ExecutesCommand[fixtures.MessageC](),
					)
				},
			})

			process2 = FromProcess(&fixtures.ProcessMessageHandler{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("<proc2-name>", "3b4ce9af-ca54-4c77-a8e7-285267f73c82")
					c.Routes(
						dogma.HandlesEvent[fixtures.MessageE](),
						dogma.ExecutesCommand[fixtures.MessageC](),
					)
				},
			})

			integration1 = FromIntegration(&fixtures.IntegrationMessageHandler{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("<int1-name>", "22857e0c-7990-4dfe-9cd0-40d6dd160aaf")
					c.Routes(
						dogma.HandlesCommand[fixtures.MessageC](),
						dogma.RecordsEvent[fixtures.MessageD](),
					)
				},
			})

			integration2 = FromIntegration(&fixtures.IntegrationMessageHandler{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("<int2-name>", "26ae7db1-7a81-407d-ac08-52e35f7765d1")
					c.Routes(
						dogma.HandlesCommand[fixtures.MessageC](),
						dogma.RecordsEvent[fixtures.MessageD](),
					)
				},
			})

			projection1 = FromProjection(&fixtures.ProjectionMessageHandler{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("<proj1-name>", "400a4609-9e00-4ccd-8436-3ad9ef073f5d")
					c.Routes(
						dogma.HandlesEvent[fixtures.MessageE](),
					)
				},
			})

			projection2 = FromProjection(&fixtures.ProjectionMessageHandler{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("<proj2-name>", "2d09f134-8971-4dff-8b84-b0e3c279ca88")
					c.Routes(
						dogma.HandlesEvent[fixtures.MessageE](),
					)
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

		Describe("func Aggregates()", func() {
			It("returns a slice containing the aggregates", func() {
				Expect(set.Aggregates()).To(ConsistOf(aggregate1, aggregate2))
			})
		})

		Describe("func Processes()", func() {
			It("returns a slice containing the processes", func() {
				Expect(set.Processes()).To(ConsistOf(process1, process2))
			})
		})

		Describe("func Integrations()", func() {
			It("returns a slice containing the integrations", func() {
				Expect(set.Integrations()).To(ConsistOf(integration1, integration2))
			})
		})

		Describe("func Projections()", func() {
			It("returns a slice containing the projections", func() {
				Expect(set.Projections()).To(ConsistOf(projection1, projection2))
			})
		})

		Describe("func RangeAggregates()", func() {
			It("calls fn for each aggregate in the set", func() {
				var names []string

				all := set.RangeAggregates(func(h RichAggregate) bool {
					names = append(names, h.Identity().Name)
					return true
				})

				Expect(names).To(ConsistOf("<agg1-name>", "<agg2-name>"))
				Expect(all).To(BeTrue())
			})

			It("stops iterating if fn returns false", func() {
				count := 0

				all := set.RangeAggregates(func(h RichAggregate) bool {
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

				all := set.RangeProcesses(func(h RichProcess) bool {
					names = append(names, h.Identity().Name)
					return true
				})

				Expect(names).To(ConsistOf("<proc1-name>", "<proc2-name>"))
				Expect(all).To(BeTrue())
			})

			It("stops iterating if fn returns false", func() {
				count := 0

				all := set.RangeProcesses(func(h RichProcess) bool {
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

				all := set.RangeIntegrations(func(h RichIntegration) bool {
					names = append(names, h.Identity().Name)
					return true
				})

				Expect(names).To(ConsistOf("<int1-name>", "<int2-name>"))
				Expect(all).To(BeTrue())
			})

			It("stops iterating if fn returns false", func() {
				count := 0

				all := set.RangeIntegrations(func(h RichIntegration) bool {
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

				all := set.RangeProjections(func(h RichProjection) bool {
					names = append(names, h.Identity().Name)
					return true
				})

				Expect(names).To(ConsistOf("<proj1-name>", "<proj2-name>"))
				Expect(all).To(BeTrue())
			})

			It("stops iterating if fn returns false", func() {
				count := 0

				all := set.RangeProjections(func(h RichProjection) bool {
					count++
					return false
				})

				Expect(count).To(BeNumerically("==", 1))
				Expect(all).To(BeFalse())
			})
		})
	})
})
