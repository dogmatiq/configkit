package api_test

import (
	"context"
	"net"
	"time"

	"github.com/dogmatiq/configkit"
	. "github.com/dogmatiq/configkit/api"
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/dogma/fixtures"
	"github.com/dogmatiq/interopspec/configspec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
)

var _ = Context("end-to-end tests", func() {
	var (
		ctx        context.Context
		cancel     func()
		app1, app2 dogma.Application
		cfg1, cfg2 configkit.Application
		listener   net.Listener
		gserver    *grpc.Server
		client     *Client
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)

		app1 = &Application{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("<app-1>", "b1101bbf-8a62-436d-9044-e6fd3d0e5385")

				c.RegisterAggregate(&AggregateMessageHandler{
					ConfigureFunc: func(c dogma.AggregateConfigurer) {
						c.Identity("<aggregate>", "938b829d-e4d7-4780-bf06-ea349453ba8f")
						c.Routes(
							dogma.HandlesCommand[MessageC](),
							dogma.RecordsEvent[MessageE](),
						)
					},
				})

				c.RegisterProcess(&ProcessMessageHandler{
					ConfigureFunc: func(c dogma.ProcessConfigurer) {
						c.Identity("<process>", "2a87972b-547d-416b-b6e5-4dddb1187658")
						c.Routes(
							dogma.HandlesEvent[MessageE](),
							dogma.ExecutesCommand[MessageC](),
							dogma.SchedulesTimeout[MessageT](),
						)
					},
				})
			},
		}

		app2 = &Application{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("<app-2>", "7d3927ce-d879-40a4-bd67-0fafc79d3c36")

				c.RegisterIntegration(&IntegrationMessageHandler{
					ConfigureFunc: func(c dogma.IntegrationConfigurer) {
						c.Identity("<integration>", "e6f0ad02-d301-4f46-a03d-4f9d0d20f5cf")
						c.Routes(
							dogma.HandlesCommand[MessageI](),
							dogma.RecordsEvent[MessageJ](),
						)
					},
				})

				c.RegisterProjection(&ProjectionMessageHandler{
					ConfigureFunc: func(c dogma.ProjectionConfigurer) {
						c.Identity("<projection>", "280a58bd-f154-46d7-863b-23ce70e49d2a")
						c.Routes(
							dogma.HandlesEvent[MessageE](),
							dogma.HandlesEvent[MessageJ](),
						)
						c.Disable()
					},
				})
			},
		}

		cfg1 = configkit.FromApplication(app1)
		cfg2 = configkit.FromApplication(app2)

		var err error
		listener, err = net.Listen("tcp", ":")
		Expect(err).ShouldNot(HaveOccurred())

		gserver = grpc.NewServer()
		configspec.RegisterConfigAPIServer(
			gserver,
			NewServer(cfg1, cfg2),
		)

		go gserver.Serve(listener)

		conn, err := grpc.Dial(
			listener.Addr().String(),
			grpc.WithInsecure(),
		)
		Expect(err).ShouldNot(HaveOccurred())

		client = NewClient(conn)
	})

	AfterEach(func() {
		if listener != nil {
			listener.Close()
		}

		if gserver != nil {
			gserver.Stop()
		}

		cancel()
	})

	Describe("func ListApplications()", func() {
		It("returns the application configurations", func() {
			configs, err := client.ListApplications(ctx)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(configs).To(HaveLen(2))

			var res1, res2 configkit.Application

			for _, cfg := range configs {
				switch cfg.Identity() {
				case cfg1.Identity():
					res1 = cfg
				case cfg2.Identity():
					res2 = cfg
				default:
					Fail("unexpected config in response")
				}
			}

			if !configkit.IsApplicationEqual(res1, cfg1) {
				Fail(
					"expected:\n\n" +
						configkit.ToString(res1) +
						"\nto equal:\n\n" +
						configkit.ToString(cfg1),
				)
			}

			if !configkit.IsApplicationEqual(res2, cfg2) {
				Fail(
					"expected:\n\n" +
						configkit.ToString(res2) +
						"\nto equal:\n\n" +
						configkit.ToString(cfg2),
				)
			}
		})

		It("returns an error if the gRPC call fails", func() {
			gserver.Stop()
			_, err := client.ListApplications(ctx)
			Expect(err).Should(HaveOccurred())
		})
	})
})
