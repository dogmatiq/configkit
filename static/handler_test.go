package static_test

import (
	"github.com/dogmatiq/configkit"
	cfixtures "github.com/dogmatiq/configkit/fixtures"
	"github.com/dogmatiq/configkit/message"
	. "github.com/dogmatiq/configkit/static"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/tools/go/packages"
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
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/handlers/single",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
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

		When("deprecated configurer methods are called", func() {
			It("returns a single configuration for each handler type", func() {
				cfg := packages.Config{
					Mode: packages.LoadAllSyntax,
					Dir:  "testdata/handlers/deprecated/single",
				}

				pkgs, err := packages.Load(&cfg, "./...")
				Expect(err).NotTo(HaveOccurred())

				apps := FromPackages(pkgs)
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
							Key:  "bf11e5eb-8cda-4498-a12e-35bf224aade7",
						},
					),
				)
				Expect(aggregate.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/deprecated/single.AggregateHandler",
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
							Key:  "b1ed1327-01fd-44ce-9cb8-c25a560e4c92",
						},
					),
				)
				Expect(process.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/deprecated/single.ProcessHandler",
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
							Key:  "91fc1201-5c47-40e7-ae10-ea3c96f7264d",
						},
					),
				)
				Expect(projection.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/deprecated/single.ProjectionHandler",
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
							Key:  "ad1ab39c-f497-4c85-89b7-72f62d2d6c28",
						},
					),
				)
				Expect(integration.TypeName()).To(
					Equal(
						"github.com/dogmatiq/configkit/static/testdata/handlers/deprecated/single.IntegrationHandler",
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

			When("a nil value is passed as a message", func() {
				It("does not add the message to the handler configuration", func() {
					cfg := packages.Config{
						Mode: packages.LoadAllSyntax,
						Dir:  "testdata/handlers/deprecated/nil-message",
					}

					pkgs, err := packages.Load(&cfg, "./...")
					Expect(err).NotTo(HaveOccurred())

					apps := FromPackages(pkgs)
					Expect(apps).To(HaveLen(1))
					Expect(apps[0].Handlers().Aggregates()).To(HaveLen(1))

					aggregate := apps[0].Handlers().Aggregates()[0]
					Expect(aggregate.Identity()).To(
						Equal(
							configkit.Identity{
								Name: "<nil-message-aggregate>",
								Key:  "da31f2cf-40c8-439e-aadc-042e30100908",
							},
						),
					)
					Expect(aggregate.TypeName()).To(
						Equal(
							"github.com/dogmatiq/configkit/static/testdata/handlers/deprecated/nil-message.AggregateHandler",
						),
					)
					Expect(aggregate.HandlerType()).To(Equal(configkit.AggregateHandlerType))
					Expect(aggregate.MessageNames()).To(Equal(
						configkit.EntityMessageNames{
							Consumed: message.NameRoles{
								cfixtures.MessageATypeName: message.CommandRole,
							},
							Produced: message.NameRoles{
								cfixtures.MessageBTypeName: message.EventRole,
							},
						},
					))

					progress := apps[0].Handlers().Processes()[0]
					Expect(progress.Identity()).To(
						Equal(
							configkit.Identity{
								Name: "<nil-message-process>",
								Key:  "16dea3ff-095f-4788-b632-3c6dd6903417",
							},
						),
					)
					Expect(progress.TypeName()).To(
						Equal(
							"github.com/dogmatiq/configkit/static/testdata/handlers/deprecated/nil-message.ProcessHandler",
						),
					)
					Expect(progress.HandlerType()).To(Equal(configkit.ProcessHandlerType))
					Expect(progress.MessageNames()).To(Equal(
						configkit.EntityMessageNames{
							Consumed: message.NameRoles{
								cfixtures.MessageATypeName: message.EventRole,
								cfixtures.MessageCTypeName: message.TimeoutRole,
							},
							Produced: message.NameRoles{
								cfixtures.MessageBTypeName: message.CommandRole,
								cfixtures.MessageCTypeName: message.TimeoutRole,
							},
						},
					))

					projection := apps[0].Handlers().Projections()[0]
					Expect(projection.Identity()).To(
						Equal(
							configkit.Identity{
								Name: "<nil-message-projection>",
								Key:  "ccaff8ea-f3c4-4d5c-8216-cb408b792998",
							},
						),
					)
					Expect(projection.TypeName()).To(
						Equal(
							"github.com/dogmatiq/configkit/static/testdata/handlers/deprecated/nil-message.ProjectionHandler",
						),
					)
					Expect(projection.HandlerType()).To(Equal(configkit.ProjectionHandlerType))
					Expect(projection.MessageNames()).To(Equal(
						configkit.EntityMessageNames{
							Consumed: message.NameRoles{
								cfixtures.MessageATypeName: message.EventRole,
							},
							Produced: message.NameRoles{},
						},
					))

					integration := apps[0].Handlers().Integrations()[0]
					Expect(integration.Identity()).To(
						Equal(
							configkit.Identity{
								Name: "<nil-message-integration>",
								Key:  "6042d127-d64c-4bfa-88ca-a6b1e0055759",
							},
						),
					)
					Expect(integration.TypeName()).To(
						Equal(
							"github.com/dogmatiq/configkit/static/testdata/handlers/deprecated/nil-message.IntegrationHandler",
						),
					)
					Expect(integration.HandlerType()).To(Equal(configkit.IntegrationHandlerType))
					Expect(integration.MessageNames()).To(Equal(
						configkit.EntityMessageNames{
							Consumed: message.NameRoles{
								cfixtures.MessageATypeName: message.CommandRole,
							},
							Produced: message.NameRoles{
								cfixtures.MessageBTypeName: message.EventRole,
							},
						},
					))
				})
			})
		})
	})
	When("the application multiple handlers of each type", func() {
		It("returns all of the handler configurations", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/handlers/multiple",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
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
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/handlers/nil-handler",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)
			Expect(apps).To(HaveLen(1))
			Expect(apps[0].Handlers()).To(Equal(configkit.HandlerSet{}))
		})
	})
})
