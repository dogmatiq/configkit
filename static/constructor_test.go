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

var _ = Describe("func FromPackages() (constructor function)", func() {
	When("the handler is created by a call to a 'constructor' function", func() {
		It("builds the configuration from the adapted type", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/handlers/constructor",
			}

			pkgs := loadPackages(cfg)

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
					"github.com/dogmatiq/configkit/static/testdata/handlers/constructor.AggregateHandler",
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
					"github.com/dogmatiq/configkit/static/testdata/handlers/constructor.ProcessHandler",
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
					"github.com/dogmatiq/configkit/static/testdata/handlers/constructor.ProjectionHandler",
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
					"github.com/dogmatiq/configkit/static/testdata/handlers/constructor.IntegrationHandler",
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
