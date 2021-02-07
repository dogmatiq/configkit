package static_test

import (
	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/internal/entity"
	"github.com/dogmatiq/configkit/message"
	. "github.com/dogmatiq/configkit/static"
	"github.com/dogmatiq/dogma/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"golang.org/x/tools/go/packages"
)

var _ = Describe("func FromPackages() (aggregate analysis)", func() {
	When("the application contains a single aggregate handler", func() {
		It("returns the aggregate handler configuration", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/aggregates/single-aggregate-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)

			Expect(apps).To(HaveLen(1))
			Expect(apps[0].Handlers()).To(
				MatchAllKeys(
					Keys{
						configkit.Identity{
							Name: "<aggregate>",
							Key:  "ef16c9d1-d7b6-4c99-a0e7-a59218e544fc",
						}: PointTo(
							MatchAllFields(Fields{
								"IdentityValue": Equal(
									configkit.Identity{
										Name: "<aggregate>",
										Key:  "ef16c9d1-d7b6-4c99-a0e7-a59218e544fc",
									},
								),
								"TypeNameValue": Equal(
									"github.com/dogmatiq/configkit/static/testdata/aggregates/single-aggregate-app.AggregateHandler",
								),
								"MessageNamesValue": MatchAllFields(Fields{
									"Consumed": MatchAllKeys(Keys{
										message.NameOf(fixtures.MessageA{}): Equal(message.CommandRole),
										message.NameOf(fixtures.MessageB{}): Equal(message.CommandRole),
									}),
									"Produced": MatchAllKeys(Keys{
										message.NameOf(fixtures.MessageC{}): Equal(message.EventRole),
										message.NameOf(fixtures.MessageD{}): Equal(message.EventRole),
									}),
								}),
								"HandlerTypeValue": Equal(configkit.AggregateHandlerType),
							}),
						),
					},
				),
			)
		})
	})

	When("the application contains multiple aggregate handlers", func() {
		It("returns all of the aggregate handler configurations", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/aggregates/multiple-aggregate-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)

			Expect(apps).To(HaveLen(1))
			Expect(apps[0].Handlers()).To(
				Equal(
					configkit.HandlerSet{
						configkit.Identity{
							Name: "<first-aggregate>",
							Key:  "e6300d8d-6530-405e-9729-e9ca21df23d3",
						}: &entity.Handler{
							IdentityValue: configkit.Identity{
								Name: "<first-aggregate>",
								Key:  "e6300d8d-6530-405e-9729-e9ca21df23d3",
							},
							TypeNameValue: "github.com/dogmatiq/configkit/static/testdata/aggregates/multiple-aggregate-app.FirstAggregateHandler",
							MessageNamesValue: configkit.EntityMessageNames{
								Consumed: message.NameRoles{
									message.NameOf(fixtures.MessageA{}): message.CommandRole,
								},
								Produced: message.NameRoles{
									message.NameOf(fixtures.MessageB{}): message.EventRole,
								},
							},
							HandlerTypeValue: configkit.AggregateHandlerType,
						},
						configkit.Identity{
							Name: "<second-aggregate>",
							Key:  "feeb96d0-c56b-4e58-9cd0-d393683c2ec7",
						}: &entity.Handler{
							IdentityValue: configkit.Identity{
								Name: "<second-aggregate>",
								Key:  "feeb96d0-c56b-4e58-9cd0-d393683c2ec7",
							},
							TypeNameValue: "github.com/dogmatiq/configkit/static/testdata/aggregates/multiple-aggregate-app.SecondAggregateHandler",
							MessageNamesValue: configkit.EntityMessageNames{
								Consumed: message.NameRoles{
									message.NameOf(fixtures.MessageC{}): message.CommandRole,
								},
								Produced: message.NameRoles{
									message.NameOf(fixtures.MessageD{}): message.EventRole,
								},
							},
							HandlerTypeValue: configkit.AggregateHandlerType,
						},
					},
				),
			)
		})
	})

	When("a nil value passed as an aggregate handler", func() {
		It("does not returns the aggregate handler configuration", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/aggregates/nil-aggregate-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)

			Expect(apps).To(HaveLen(1))
			Expect(apps[0].Handlers()).To(Equal(configkit.HandlerSet{}))
		})
	})

	When("a nil value passed as a message while configuring an aggregate handler", func() {
		It("skips the nil value message", func() {
			cfg := packages.Config{
				Mode: packages.LoadAllSyntax,
				Dir:  "testdata/aggregates/nil-message-aggregate-app",
			}

			pkgs, err := packages.Load(&cfg, "./...")
			Expect(err).NotTo(HaveOccurred())

			apps := FromPackages(pkgs)

			Expect(apps).To(HaveLen(1))
			Expect(apps[0].Handlers()).To(
				Equal(
					configkit.HandlerSet{
						configkit.Identity{
							Name: "<nil-message-aggregate>",
							Key:  "99271492-1ec3-475f-b154-3e69cda11155",
						}: &entity.Handler{
							IdentityValue: configkit.Identity{
								Name: "<nil-message-aggregate>",
								Key:  "99271492-1ec3-475f-b154-3e69cda11155",
							},
							TypeNameValue: "github.com/dogmatiq/configkit/static/testdata/aggregates/nil-message-aggregate-app.AggregateHandler",
							MessageNamesValue: configkit.EntityMessageNames{
								Consumed: message.NameRoles{
									message.NameOf(fixtures.MessageA{}): message.CommandRole,
								},
								Produced: message.NameRoles{
									message.NameOf(fixtures.MessageB{}): message.EventRole,
								},
							},
							HandlerTypeValue: configkit.AggregateHandlerType,
						},
					},
				),
			)
		})
	})
})
