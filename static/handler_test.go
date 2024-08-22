package static_test

import (
	"os"

	"github.com/dogmatiq/configkit"
	cfixtures "github.com/dogmatiq/configkit/fixtures"
	"github.com/dogmatiq/configkit/message"
	. "github.com/dogmatiq/configkit/static"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// matchIdentities matches the given identities to those of the handlers in the
// handler set.
func matchIdentities(
	hs configkit.HandlerSet,
	identities ...configkit.Identity,
) {
	for _, identity := range identities {
		_, ok := hs.ByIdentity(identity)
		Expect(ok).To(BeTrue())
	}
}

var _ = Describe("func FromPackages() (handler analysis)", func() {
	When("the application contains a single handler of each type", func() {
		It("returns a single configuration for each handler type", func() {
			apps := FromDir("testdata/handlers/single")
			Expect(apps).To(HaveLen(1))
			Expect(apps[0].Handlers().Aggregates()).To(HaveLen(1))
			Expect(apps[0].Handlers().Processes()).To(HaveLen(1))
			Expect(apps[0].Handlers().Projections()).To(HaveLen(1))
			Expect(apps[0].Handlers().Integrations()).To(HaveLen(1))

			aggregate := apps[0].Handlers().Aggregates()[0]
			Expect(aggregate.Identity()).To(
				Equal(
					configkit.Identity{
						Name: "<aggregate>",
						Key:  "ef16c9d1-d7b6-4c99-a0e7-a59218e544fc",
					},
				),
			)
			Expect(aggregate.TypeName()).To(
				Equal(
					"github.com/dogmatiq/configkit/static/testdata/handlers/single.AggregateHandler",
				),
			)
			Expect(aggregate.HandlerType()).To(Equal(configkit.AggregateHandlerType))

			Expect(aggregate.MessageNames()).To(Equal(
				configkit.EntityMessageNames{
					Consumed: message.NameRoles{
						cfixtures.MessageATypeName: message.CommandRole,
						cfixtures.MessageBTypeName: message.CommandRole,
					},
					Produced: message.NameRoles{
						cfixtures.MessageCTypeName: message.EventRole,
						cfixtures.MessageDTypeName: message.EventRole,
					},
				},
			))

			process := apps[0].Handlers().Processes()[0]
			Expect(process.Identity()).To(
				Equal(
					configkit.Identity{
						Name: "<process>",
						Key:  "5e839b73-170b-42c0-bf41-8feee4b5a583",
					},
				),
			)
			Expect(process.TypeName()).To(
				Equal(
					"github.com/dogmatiq/configkit/static/testdata/handlers/single.ProcessHandler",
				),
			)
			Expect(process.HandlerType()).To(Equal(configkit.ProcessHandlerType))

			Expect(process.MessageNames()).To(Equal(
				configkit.EntityMessageNames{
					Consumed: message.NameRoles{
						cfixtures.MessageATypeName: message.EventRole,
						cfixtures.MessageBTypeName: message.EventRole,
						cfixtures.MessageETypeName: message.TimeoutRole,
						cfixtures.MessageFTypeName: message.TimeoutRole,
					},
					Produced: message.NameRoles{
						cfixtures.MessageCTypeName: message.CommandRole,
						cfixtures.MessageDTypeName: message.CommandRole,
						cfixtures.MessageETypeName: message.TimeoutRole,
						cfixtures.MessageFTypeName: message.TimeoutRole,
					},
				},
			))

			projection := apps[0].Handlers().Projections()[0]
			Expect(projection.Identity()).To(
				Equal(
					configkit.Identity{
						Name: "<projection>",
						Key:  "823e61d3-ace1-469d-b0a6-778e84c0a508",
					},
				),
			)
			Expect(projection.TypeName()).To(
				Equal(
					"github.com/dogmatiq/configkit/static/testdata/handlers/single.ProjectionHandler",
				),
			)
			Expect(projection.HandlerType()).To(Equal(configkit.ProjectionHandlerType))

			Expect(projection.MessageNames()).To(Equal(
				configkit.EntityMessageNames{
					Consumed: message.NameRoles{
						cfixtures.MessageATypeName: message.EventRole,
						cfixtures.MessageBTypeName: message.EventRole,
					},
					Produced: message.NameRoles{},
				},
			))

			integration := apps[0].Handlers().Integrations()[0]
			Expect(integration.Identity()).To(
				Equal(
					configkit.Identity{
						Name: "<integration>",
						Key:  "099b5b8d-9e04-422f-bcc3-bb0d451158c7",
					},
				),
			)
			Expect(integration.TypeName()).To(
				Equal(
					"github.com/dogmatiq/configkit/static/testdata/handlers/single.IntegrationHandler",
				),
			)
			Expect(integration.HandlerType()).To(Equal(configkit.IntegrationHandlerType))

			Expect(integration.MessageNames()).To(Equal(
				configkit.EntityMessageNames{
					Consumed: message.NameRoles{
						cfixtures.MessageATypeName: message.CommandRole,
						cfixtures.MessageBTypeName: message.CommandRole,
					},
					Produced: message.NameRoles{
						cfixtures.MessageCTypeName: message.EventRole,
						cfixtures.MessageDTypeName: message.EventRole,
					},
				},
			))
		})
		When("a handler is a type alias", func() {
			goDbg := os.Getenv("GODEBUG")

			BeforeEach(func() {
				// Set the GODEBUG environment variable to enable type alias
				// support.
				//
				// TODO: Remove this setting once the package is
				// migrated to Go 1.23 that generates `types.Alias` types by
				// default.
				os.Setenv("GODEBUG", "gotypesalias=1")
			})

			AfterEach(func() {
				os.Setenv("GODEBUG", goDbg)
			})

			It("returns a single configuration for each handler type", func() {
				apps := FromDir("testdata/handlers/typealias")
				Expect(apps).To(HaveLen(1))
				Expect(apps[0].Handlers().Aggregates()).To(HaveLen(1))
				Expect(apps[0].Handlers().Processes()).To(HaveLen(1))
				Expect(apps[0].Handlers().Projections()).To(HaveLen(1))
				Expect(apps[0].Handlers().Integrations()).To(HaveLen(1))

				aggregate := apps[0].Handlers().Aggregates()[0]
				Expect(aggregate.Identity()).To(
					Equal(
						configkit.Identity{
							Name: "<aggregate>",
							Key:  "92623de9-c9cf-42f3-8338-33c50eeb06fb",
						},
					),
				)
				Expect(aggregate.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/typealias.AggregateHandlerAlias",
					),
				)
				Expect(aggregate.HandlerType()).To(Equal(configkit.AggregateHandlerType))

				Expect(aggregate.MessageNames()).To(Equal(
					configkit.EntityMessageNames{
						Consumed: message.NameRoles{
							cfixtures.MessageATypeName: message.CommandRole,
							cfixtures.MessageBTypeName: message.CommandRole,
						},
						Produced: message.NameRoles{
							cfixtures.MessageCTypeName: message.EventRole,
							cfixtures.MessageDTypeName: message.EventRole,
						},
					},
				))

				process := apps[0].Handlers().Processes()[0]
				Expect(process.Identity()).To(
					Equal(
						configkit.Identity{
							Name: "<process>",
							Key:  "ad9d6955-893a-4d8d-a26e-e25886b113b2",
						},
					),
				)
				Expect(process.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/typealias.ProcessHandlerAlias",
					),
				)
				Expect(process.HandlerType()).To(Equal(configkit.ProcessHandlerType))

				Expect(process.MessageNames()).To(Equal(
					configkit.EntityMessageNames{
						Consumed: message.NameRoles{
							cfixtures.MessageATypeName: message.EventRole,
							cfixtures.MessageBTypeName: message.EventRole,
							cfixtures.MessageETypeName: message.TimeoutRole,
							cfixtures.MessageFTypeName: message.TimeoutRole,
						},
						Produced: message.NameRoles{
							cfixtures.MessageCTypeName: message.CommandRole,
							cfixtures.MessageDTypeName: message.CommandRole,
							cfixtures.MessageETypeName: message.TimeoutRole,
							cfixtures.MessageFTypeName: message.TimeoutRole,
						},
					},
				))

				projection := apps[0].Handlers().Projections()[0]
				Expect(projection.Identity()).To(
					Equal(
						configkit.Identity{
							Name: "<projection>",
							Key:  "d012b7ed-3c4b-44db-9276-7bbc90fb54fd",
						},
					),
				)
				Expect(projection.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/typealias.ProjectionHandlerAlias",
					),
				)
				Expect(projection.HandlerType()).To(Equal(configkit.ProjectionHandlerType))

				Expect(projection.MessageNames()).To(Equal(
					configkit.EntityMessageNames{
						Consumed: message.NameRoles{
							cfixtures.MessageATypeName: message.EventRole,
							cfixtures.MessageBTypeName: message.EventRole,
						},
						Produced: message.NameRoles{},
					},
				))

				integration := apps[0].Handlers().Integrations()[0]
				Expect(integration.Identity()).To(
					Equal(
						configkit.Identity{
							Name: "<integration>",
							Key:  "4d8cd3f5-21dc-475b-a8dc-80138adde3f2",
						},
					),
				)
				Expect(integration.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/typealias.IntegrationHandlerAlias",
					),
				)
				Expect(integration.HandlerType()).To(Equal(configkit.IntegrationHandlerType))

				Expect(integration.MessageNames()).To(Equal(
					configkit.EntityMessageNames{
						Consumed: message.NameRoles{
							cfixtures.MessageATypeName: message.CommandRole,
							cfixtures.MessageBTypeName: message.CommandRole,
						},
						Produced: message.NameRoles{
							cfixtures.MessageCTypeName: message.EventRole,
							cfixtures.MessageDTypeName: message.EventRole,
						},
					},
				))
			})
		})

		When("messages are passed to the *Configurer.Routes() method", func() {
			It("includes messages passed as args to *Configurer.Routes() method only", func() {
				apps := FromDir("testdata/handlers/only-routes-args")
				Expect(apps).To(HaveLen(1))
				Expect(apps[0].Handlers().Aggregates()).To(HaveLen(1))
				Expect(apps[0].Handlers().Processes()).To(HaveLen(1))
				Expect(apps[0].Handlers().Projections()).To(HaveLen(1))
				Expect(apps[0].Handlers().Integrations()).To(HaveLen(1))

				aggregate := apps[0].Handlers().Aggregates()[0]
				Expect(aggregate.Identity()).To(
					Equal(
						configkit.Identity{
							Name: "<aggregate>",
							Key:  "dcfdd034-e374-478b-8faa-bc688ff59f1f",
						},
					),
				)
				Expect(aggregate.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/only-routes-args.AggregateHandler",
					),
				)
				Expect(aggregate.HandlerType()).To(Equal(configkit.AggregateHandlerType))

				Expect(aggregate.MessageNames()).To(Equal(
					configkit.EntityMessageNames{
						Consumed: message.NameRoles{
							cfixtures.MessageATypeName: message.CommandRole,
							cfixtures.MessageBTypeName: message.CommandRole,
						},
						Produced: message.NameRoles{
							cfixtures.MessageCTypeName: message.EventRole,
							cfixtures.MessageDTypeName: message.EventRole,
						},
					},
				))

				process := apps[0].Handlers().Processes()[0]
				Expect(process.Identity()).To(
					Equal(
						configkit.Identity{
							Name: "<process>",
							Key:  "24c61438-e7ae-4d54-8e28-2fc6e848c948",
						},
					),
				)
				Expect(process.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/only-routes-args.ProcessHandler",
					),
				)
				Expect(process.HandlerType()).To(Equal(configkit.ProcessHandlerType))

				Expect(process.MessageNames()).To(Equal(
					configkit.EntityMessageNames{
						Consumed: message.NameRoles{
							cfixtures.MessageATypeName: message.EventRole,
							cfixtures.MessageBTypeName: message.EventRole,
							cfixtures.MessageETypeName: message.TimeoutRole,
							cfixtures.MessageFTypeName: message.TimeoutRole,
						},
						Produced: message.NameRoles{
							cfixtures.MessageCTypeName: message.CommandRole,
							cfixtures.MessageDTypeName: message.CommandRole,
							cfixtures.MessageETypeName: message.TimeoutRole,
							cfixtures.MessageFTypeName: message.TimeoutRole,
						},
					},
				))

				projection := apps[0].Handlers().Projections()[0]
				Expect(projection.Identity()).To(
					Equal(
						configkit.Identity{
							Name: "<projection>",
							Key:  "6b9acb05-cd77-4342-bf10-b3de9d2d5bba",
						},
					),
				)
				Expect(projection.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/only-routes-args.ProjectionHandler",
					),
				)
				Expect(projection.HandlerType()).To(Equal(configkit.ProjectionHandlerType))

				Expect(projection.MessageNames()).To(Equal(
					configkit.EntityMessageNames{
						Consumed: message.NameRoles{
							cfixtures.MessageATypeName: message.EventRole,
							cfixtures.MessageBTypeName: message.EventRole,
						},
						Produced: message.NameRoles{},
					},
				))

				integration := apps[0].Handlers().Integrations()[0]
				Expect(integration.Identity()).To(
					Equal(
						configkit.Identity{
							Name: "<integration>",
							Key:  "ac391765-da58-4e7c-a478-e4725eb2b0e9",
						},
					),
				)
				Expect(integration.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/only-routes-args.IntegrationHandler",
					),
				)
				Expect(integration.HandlerType()).To(Equal(configkit.IntegrationHandlerType))

				Expect(integration.MessageNames()).To(Equal(
					configkit.EntityMessageNames{
						Consumed: message.NameRoles{
							cfixtures.MessageATypeName: message.CommandRole,
							cfixtures.MessageBTypeName: message.CommandRole,
						},
						Produced: message.NameRoles{
							cfixtures.MessageCTypeName: message.EventRole,
							cfixtures.MessageDTypeName: message.EventRole,
						},
					},
				))
			})
		})

		When("messages are passed to the *Configurer.Routes() method as a dynamically populated splice", func() {
			It("returns a single configuration for each handler type", func() {
				apps := FromDir("testdata/handlers/dynamic-routes")
				Expect(apps).To(HaveLen(1))
				Expect(apps[0].Handlers().Aggregates()).To(HaveLen(1))
				Expect(apps[0].Handlers().Processes()).To(HaveLen(1))
				Expect(apps[0].Handlers().Projections()).To(HaveLen(1))
				Expect(apps[0].Handlers().Integrations()).To(HaveLen(1))

				aggregate := apps[0].Handlers().Aggregates()[0]
				Expect(aggregate.Identity()).To(
					Equal(
						configkit.Identity{
							Name: "<aggregate>",
							Key:  "3876b4e5-8759-4c0b-bf0b-03ef49777e5c",
						},
					),
				)
				Expect(aggregate.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/dynamic-routes.AggregateHandler",
					),
				)
				Expect(aggregate.HandlerType()).To(Equal(configkit.AggregateHandlerType))

				Expect(aggregate.MessageNames()).To(Equal(
					configkit.EntityMessageNames{
						Consumed: message.NameRoles{
							cfixtures.MessageATypeName: message.CommandRole,
							cfixtures.MessageBTypeName: message.CommandRole,
						},
						Produced: message.NameRoles{
							cfixtures.MessageCTypeName: message.EventRole,
							cfixtures.MessageDTypeName: message.EventRole,
						},
					},
				))

				process := apps[0].Handlers().Processes()[0]
				Expect(process.Identity()).To(
					Equal(
						configkit.Identity{
							Name: "<process>",
							Key:  "d00131d2-f99c-44e3-be11-104e69b04f77",
						},
					),
				)
				Expect(process.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/dynamic-routes.ProcessHandler",
					),
				)
				Expect(process.HandlerType()).To(Equal(configkit.ProcessHandlerType))

				Expect(process.MessageNames()).To(Equal(
					configkit.EntityMessageNames{
						Consumed: message.NameRoles{
							cfixtures.MessageATypeName: message.EventRole,
							cfixtures.MessageBTypeName: message.EventRole,
							cfixtures.MessageETypeName: message.TimeoutRole,
							cfixtures.MessageFTypeName: message.TimeoutRole,
						},
						Produced: message.NameRoles{
							cfixtures.MessageCTypeName: message.CommandRole,
							cfixtures.MessageDTypeName: message.CommandRole,
							cfixtures.MessageETypeName: message.TimeoutRole,
							cfixtures.MessageFTypeName: message.TimeoutRole,
						},
					},
				))

				projection := apps[0].Handlers().Projections()[0]
				Expect(projection.Identity()).To(
					Equal(
						configkit.Identity{
							Name: "<projection>",
							Key:  "6f20c336-b740-4249-a80b-d94d3bdce750",
						},
					),
				)
				Expect(projection.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/dynamic-routes.ProjectionHandler",
					),
				)
				Expect(projection.HandlerType()).To(Equal(configkit.ProjectionHandlerType))

				Expect(projection.MessageNames()).To(Equal(
					configkit.EntityMessageNames{
						Consumed: message.NameRoles{
							cfixtures.MessageATypeName: message.EventRole,
							cfixtures.MessageBTypeName: message.EventRole,
						},
						Produced: message.NameRoles{},
					},
				))

				integration := apps[0].Handlers().Integrations()[0]
				Expect(integration.Identity()).To(
					Equal(
						configkit.Identity{
							Name: "<integration>",
							Key:  "3a06b7da-1079-4e4b-a6a6-064c62241918",
						},
					),
				)
				Expect(integration.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/dynamic-routes.IntegrationHandler",
					),
				)
				Expect(integration.HandlerType()).To(Equal(configkit.IntegrationHandlerType))

				Expect(integration.MessageNames()).To(Equal(
					configkit.EntityMessageNames{
						Consumed: message.NameRoles{
							cfixtures.MessageATypeName: message.CommandRole,
							cfixtures.MessageBTypeName: message.CommandRole,
						},
						Produced: message.NameRoles{
							cfixtures.MessageCTypeName: message.EventRole,
							cfixtures.MessageDTypeName: message.EventRole,
						},
					},
				))
			})
		})

		When("messages are passed to the *Configurer.Routes() method in conditional branches", func() {
			It("returns messages populated in every conditional branch", func() {
				apps := FromDir("testdata/handlers/conditional-branches")
				Expect(apps).To(HaveLen(1))
				Expect(apps[0].Handlers().Aggregates()).To(HaveLen(1))
				Expect(apps[0].Handlers().Processes()).To(HaveLen(1))
				Expect(apps[0].Handlers().Projections()).To(HaveLen(1))
				Expect(apps[0].Handlers().Integrations()).To(HaveLen(1))

				aggregate := apps[0].Handlers().Aggregates()[0]
				Expect(aggregate.Identity()).To(
					Equal(
						configkit.Identity{
							Name: "<aggregate>",
							Key:  "c3b4b3c7-fbe6-4789-9358-e4f45b154d31",
						},
					),
				)
				Expect(aggregate.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/conditional-branches.AggregateHandler",
					),
				)
				Expect(aggregate.HandlerType()).To(Equal(configkit.AggregateHandlerType))

				Expect(aggregate.MessageNames()).To(Equal(
					configkit.EntityMessageNames{
						Consumed: message.NameRoles{
							cfixtures.MessageATypeName: message.CommandRole,
							cfixtures.MessageBTypeName: message.CommandRole,
						},
						Produced: message.NameRoles{
							cfixtures.MessageCTypeName: message.EventRole,
							cfixtures.MessageDTypeName: message.EventRole,
						},
					},
				))

				process := apps[0].Handlers().Processes()[0]
				Expect(process.Identity()).To(
					Equal(
						configkit.Identity{
							Name: "<process>",
							Key:  "f754da79-205b-4d65-889f-0d8ae86e3def",
						},
					),
				)
				Expect(process.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/conditional-branches.ProcessHandler",
					),
				)
				Expect(process.HandlerType()).To(Equal(configkit.ProcessHandlerType))

				Expect(process.MessageNames()).To(Equal(
					configkit.EntityMessageNames{
						Consumed: message.NameRoles{
							cfixtures.MessageATypeName: message.EventRole,
							cfixtures.MessageBTypeName: message.EventRole,
							cfixtures.MessageETypeName: message.TimeoutRole,
							cfixtures.MessageFTypeName: message.TimeoutRole,
						},
						Produced: message.NameRoles{
							cfixtures.MessageCTypeName: message.CommandRole,
							cfixtures.MessageDTypeName: message.CommandRole,
							cfixtures.MessageETypeName: message.TimeoutRole,
							cfixtures.MessageFTypeName: message.TimeoutRole,
						},
					},
				))

				projection := apps[0].Handlers().Projections()[0]
				Expect(projection.Identity()).To(
					Equal(
						configkit.Identity{
							Name: "<projection>",
							Key:  "559dcb05-2b63-4567-bb25-3f69c569f8ec",
						},
					),
				)
				Expect(projection.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/conditional-branches.ProjectionHandler",
					),
				)
				Expect(projection.HandlerType()).To(Equal(configkit.ProjectionHandlerType))

				Expect(projection.MessageNames()).To(Equal(
					configkit.EntityMessageNames{
						Consumed: message.NameRoles{
							cfixtures.MessageATypeName: message.EventRole,
							cfixtures.MessageBTypeName: message.EventRole,
						},
						Produced: message.NameRoles{},
					},
				))

				integration := apps[0].Handlers().Integrations()[0]
				Expect(integration.Identity()).To(
					Equal(
						configkit.Identity{
							Name: "<integration>",
							Key:  "92cce461-8d30-409b-8d5a-406f656cef2d",
						},
					),
				)
				Expect(integration.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/conditional-branches.IntegrationHandler",
					),
				)
				Expect(integration.HandlerType()).To(Equal(configkit.IntegrationHandlerType))

				Expect(integration.MessageNames()).To(Equal(
					configkit.EntityMessageNames{
						Consumed: message.NameRoles{
							cfixtures.MessageATypeName: message.CommandRole,
							cfixtures.MessageBTypeName: message.CommandRole,
						},
						Produced: message.NameRoles{
							cfixtures.MessageCTypeName: message.EventRole,
							cfixtures.MessageDTypeName: message.EventRole,
						},
					},
				))
			})
		})

		When("nil is passed to a call of *Configurer.Routes() methods", func() {
			It("does not populate messages", func() {
				apps := FromDir("testdata/handlers/nil-routes")
				Expect(apps).To(HaveLen(1))
				Expect(apps[0].Handlers().Aggregates()).To(HaveLen(1))
				Expect(apps[0].Handlers().Processes()).To(HaveLen(1))
				Expect(apps[0].Handlers().Projections()).To(HaveLen(1))
				Expect(apps[0].Handlers().Integrations()).To(HaveLen(1))

				aggregate := apps[0].Handlers().Aggregates()[0]
				Expect(aggregate.Identity()).To(
					Equal(
						configkit.Identity{
							Name: "<aggregate>",
							Key:  "df648962-7d96-427e-8bc2-5a4efdb4cc4b",
						},
					),
				)
				Expect(aggregate.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/nil-routes.AggregateHandler",
					),
				)
				Expect(aggregate.HandlerType()).To(Equal(configkit.AggregateHandlerType))

				Expect(aggregate.MessageNames()).To(Equal(
					configkit.EntityMessageNames{
						Consumed: message.NameRoles{},
						Produced: message.NameRoles{},
					},
				))

				process := apps[0].Handlers().Processes()[0]
				Expect(process.Identity()).To(
					Equal(
						configkit.Identity{
							Name: "<process>",
							Key:  "e7bcb97c-627e-44db-ba05-d8a86e5bf88b",
						},
					),
				)
				Expect(process.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/nil-routes.ProcessHandler",
					),
				)
				Expect(process.HandlerType()).To(Equal(configkit.ProcessHandlerType))

				Expect(process.MessageNames()).To(Equal(
					configkit.EntityMessageNames{
						Consumed: message.NameRoles{},
						Produced: message.NameRoles{},
					},
				))

				projection := apps[0].Handlers().Projections()[0]
				Expect(projection.Identity()).To(
					Equal(
						configkit.Identity{
							Name: "<projection>",
							Key:  "9208f704-641c-44e4-91dc-5274598b30bd",
						},
					),
				)
				Expect(projection.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/nil-routes.ProjectionHandler",
					),
				)
				Expect(projection.HandlerType()).To(Equal(configkit.ProjectionHandlerType))

				Expect(projection.MessageNames()).To(Equal(
					configkit.EntityMessageNames{
						Consumed: message.NameRoles{},
						Produced: message.NameRoles{},
					},
				))

				integration := apps[0].Handlers().Integrations()[0]
				Expect(integration.Identity()).To(
					Equal(
						configkit.Identity{
							Name: "<integration>",
							Key:  "363039e5-2938-4b2c-9bec-dcb29dee2da1",
						},
					),
				)
				Expect(integration.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/nil-routes.IntegrationHandler",
					),
				)
				Expect(integration.HandlerType()).To(Equal(configkit.IntegrationHandlerType))

				Expect(integration.MessageNames()).To(Equal(
					configkit.EntityMessageNames{
						Consumed: message.NameRoles{},
						Produced: message.NameRoles{},
					},
				))
			})
		})
	})

	When("the application multiple handlers of each type", func() {
		It("returns all of the handler configurations", func() {
			apps := FromDir("testdata/handlers/multiple")
			Expect(apps).To(HaveLen(1))

			matchIdentities(
				apps[0].Handlers(),
				configkit.Identity{
					Name: "<first-aggregate>",
					Key:  "e6300d8d-6530-405e-9729-e9ca21df23d3",
				},
				configkit.Identity{
					Name: "<second-aggregate>",
					Key:  "feeb96d0-c56b-4e58-9cd0-d393683c2ec7",
				},
				configkit.Identity{
					Name: "<first-process>",
					Key:  "d33198e0-f1f7-4c2d-8ac2-98f68a44414e",
				},
				configkit.Identity{
					Name: "<second-process>",
					Key:  "0311717c-cf51-4292-8ed9-95125302a18e",
				},
				configkit.Identity{
					Name: "<first-projection>",
					Key:  "9174783f-4f12-4619-b5c6-c4ab70bd0937",
				},
				configkit.Identity{
					Name: "<second-projection>",
					Key:  "2e22850e-7c84-4b3f-b8b3-25ac743d90f2",
				},
				configkit.Identity{
					Name: "<first-integration>",
					Key:  "14cf2812-eead-43b3-9c9c-10db5b469e94",
				},
				configkit.Identity{
					Name: "<second-integration>",
					Key:  "6bed3fbc-30e2-44c7-9a5b-e440ffe370d9",
				},
			)
		})
	})

	When("a nil value is passed as a handler", func() {
		It("does not add a handler to the application configuration", func() {
			apps := FromDir("testdata/handlers/nil-handler")
			Expect(apps).To(HaveLen(1))
			Expect(apps[0].Handlers()).To(Equal(configkit.HandlerSet{}))
		})
	})

	When("a handler with a non-pointer methodset is registered as a pointer", func() {
		It("includes the handler in the application configuration", func() {
			apps := FromDir("testdata/handlers/pointer-handler-with-non-pointer-methodset")
			Expect(apps).To(HaveLen(1))
			Expect(apps[0].Handlers().Aggregates()).To(HaveLen(1))
			Expect(apps[0].Handlers().Processes()).To(HaveLen(1))
			Expect(apps[0].Handlers().Projections()).To(HaveLen(1))
			Expect(apps[0].Handlers().Integrations()).To(HaveLen(1))

			aggregate := apps[0].Handlers().Aggregates()[0]
			Expect(aggregate.Identity()).To(
				Equal(
					configkit.Identity{
						Name: "<aggregate>",
						Key:  "ee1814e3-194d-438a-916e-ee7766598646",
					},
				),
			)
			Expect(aggregate.TypeName()).To(
				Equal(
					"*github.com/dogmatiq/configkit/static/testdata/handlers/pointer-handler-with-non-pointer-methodset.AggregateHandler",
				),
			)
			Expect(aggregate.HandlerType()).To(Equal(configkit.AggregateHandlerType))

			Expect(aggregate.MessageNames()).To(Equal(
				configkit.EntityMessageNames{
					Consumed: message.NameRoles{
						cfixtures.MessageATypeName: message.CommandRole,
						cfixtures.MessageBTypeName: message.CommandRole,
					},
					Produced: message.NameRoles{
						cfixtures.MessageCTypeName: message.EventRole,
						cfixtures.MessageDTypeName: message.EventRole,
					},
				},
			))

			process := apps[0].Handlers().Processes()[0]
			Expect(process.Identity()).To(
				Equal(
					configkit.Identity{
						Name: "<process>",
						Key:  "39af6b34-5fa1-4f3a-b049-40a5e1d9b33b",
					},
				),
			)
			Expect(process.TypeName()).To(
				Equal(
					"*github.com/dogmatiq/configkit/static/testdata/handlers/pointer-handler-with-non-pointer-methodset.ProcessHandler",
				),
			)
			Expect(process.HandlerType()).To(Equal(configkit.ProcessHandlerType))

			Expect(process.MessageNames()).To(Equal(
				configkit.EntityMessageNames{
					Consumed: message.NameRoles{
						cfixtures.MessageATypeName: message.EventRole,
						cfixtures.MessageBTypeName: message.EventRole,
						cfixtures.MessageETypeName: message.TimeoutRole,
						cfixtures.MessageFTypeName: message.TimeoutRole,
					},
					Produced: message.NameRoles{
						cfixtures.MessageCTypeName: message.CommandRole,
						cfixtures.MessageDTypeName: message.CommandRole,
						cfixtures.MessageETypeName: message.TimeoutRole,
						cfixtures.MessageFTypeName: message.TimeoutRole,
					},
				},
			))

			projection := apps[0].Handlers().Projections()[0]
			Expect(projection.Identity()).To(
				Equal(
					configkit.Identity{
						Name: "<projection>",
						Key:  "3dfcd7cd-1f63-47a1-9be7-3242bd252423",
					},
				),
			)
			Expect(projection.TypeName()).To(
				Equal(
					"*github.com/dogmatiq/configkit/static/testdata/handlers/pointer-handler-with-non-pointer-methodset.ProjectionHandler",
				),
			)
			Expect(projection.HandlerType()).To(Equal(configkit.ProjectionHandlerType))

			Expect(projection.MessageNames()).To(Equal(
				configkit.EntityMessageNames{
					Consumed: message.NameRoles{
						cfixtures.MessageATypeName: message.EventRole,
						cfixtures.MessageBTypeName: message.EventRole,
					},
					Produced: message.NameRoles{},
				},
			))

			integration := apps[0].Handlers().Integrations()[0]
			Expect(integration.Identity()).To(
				Equal(
					configkit.Identity{
						Name: "<integration>",
						Key:  "1425ca64-0448-4bfd-b18d-9fe63a95995f",
					},
				),
			)
			Expect(integration.TypeName()).To(
				Equal(
					"*github.com/dogmatiq/configkit/static/testdata/handlers/pointer-handler-with-non-pointer-methodset.IntegrationHandler",
				),
			)
			Expect(integration.HandlerType()).To(Equal(configkit.IntegrationHandlerType))

			Expect(integration.MessageNames()).To(Equal(
				configkit.EntityMessageNames{
					Consumed: message.NameRoles{
						cfixtures.MessageATypeName: message.CommandRole,
						cfixtures.MessageBTypeName: message.CommandRole,
					},
					Produced: message.NameRoles{
						cfixtures.MessageCTypeName: message.EventRole,
						cfixtures.MessageDTypeName: message.EventRole,
					},
				},
			))
		})
	})
})
