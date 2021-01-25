package static_test

import (
	"github.com/dogmatiq/configkit"
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

			Expect(apps).To(
				ConsistOf(
					PointTo(
						MatchAllFields(Fields{
							"IdentityValue": Equal(
								configkit.Identity{
									Name: "<single-aggregate-app>",
									Key:  "3bc3849b-abe0-4c4e-9db4-e48dc28c9a26",
								},
							),
							"TypeNameValue": Equal(
								"github.com/dogmatiq/configkit/static/testdata/aggregates/single-aggregate-app.App",
							),
							"MessageNamesValue": Equal(configkit.EntityMessageNames{}),
							"HandlersValue": MatchAllKeys(Keys{
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
												message.NameOf(fixtures.MessageC{}): Equal(message.CommandRole),
											}),
											"Produced": MatchAllKeys(Keys{
												message.NameOf(fixtures.MessageD{}): Equal(message.EventRole),
												message.NameOf(fixtures.MessageE{}): Equal(message.EventRole),
												message.NameOf(fixtures.MessageF{}): Equal(message.EventRole),
											}),
										}),
										"HandlerTypeValue": Equal(configkit.AggregateHandlerType),
									}),
								),
							}),
						}),
					),
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

			Expect(apps).To(
				ConsistOf(
					PointTo(
						MatchAllFields(Fields{
							"IdentityValue": Equal(
								configkit.Identity{
									Name: "<multiple-aggregate-app>",
									Key:  "66abab58-e0cd-4d50-90c2-b0bd093ea2ba",
								},
							),
							"TypeNameValue": Equal(
								"github.com/dogmatiq/configkit/static/testdata/aggregates/multiple-aggregate-app.App",
							),
							"MessageNamesValue": Equal(configkit.EntityMessageNames{}),
							"HandlersValue": MatchAllKeys(Keys{
								configkit.Identity{
									Name: "<first-aggregate>",
									Key:  "e6300d8d-6530-405e-9729-e9ca21df23d3",
								}: PointTo(
									MatchAllFields(Fields{
										"IdentityValue": Equal(
											configkit.Identity{
												Name: "<first-aggregate>",
												Key:  "e6300d8d-6530-405e-9729-e9ca21df23d3",
											},
										),
										"TypeNameValue": Equal(
											"github.com/dogmatiq/configkit/static/testdata/aggregates/multiple-aggregate-app.FirstAggregateHandler",
										),
										"MessageNamesValue": MatchAllFields(Fields{
											"Consumed": MatchAllKeys(Keys{
												message.NameOf(fixtures.MessageA{}): Equal(message.CommandRole),
												message.NameOf(fixtures.MessageB{}): Equal(message.CommandRole),
												message.NameOf(fixtures.MessageC{}): Equal(message.CommandRole),
											}),
											"Produced": MatchAllKeys(Keys{
												message.NameOf(fixtures.MessageD{}): Equal(message.EventRole),
												message.NameOf(fixtures.MessageE{}): Equal(message.EventRole),
												message.NameOf(fixtures.MessageF{}): Equal(message.EventRole),
											}),
										}),
										"HandlerTypeValue": Equal(configkit.AggregateHandlerType),
									}),
								),
								configkit.Identity{
									Name: "<second-aggregate>",
									Key:  "feeb96d0-c56b-4e58-9cd0-d393683c2ec7",
								}: PointTo(
									MatchAllFields(Fields{
										"IdentityValue": Equal(
											configkit.Identity{
												Name: "<second-aggregate>",
												Key:  "feeb96d0-c56b-4e58-9cd0-d393683c2ec7",
											},
										),
										"TypeNameValue": Equal(
											"github.com/dogmatiq/configkit/static/testdata/aggregates/multiple-aggregate-app.SecondAggregateHandler",
										),
										"MessageNamesValue": MatchAllFields(Fields{
											"Consumed": MatchAllKeys(Keys{
												message.NameOf(fixtures.MessageG{}): Equal(message.CommandRole),
												message.NameOf(fixtures.MessageH{}): Equal(message.CommandRole),
												message.NameOf(fixtures.MessageI{}): Equal(message.CommandRole),
											}),
											"Produced": MatchAllKeys(Keys{
												message.NameOf(fixtures.MessageJ{}): Equal(message.EventRole),
												message.NameOf(fixtures.MessageK{}): Equal(message.EventRole),
												message.NameOf(fixtures.MessageL{}): Equal(message.EventRole),
											}),
										}),
										"HandlerTypeValue": Equal(configkit.AggregateHandlerType),
									}),
								),
							}),
						}),
					),
				),
			)
		})
	})
})
