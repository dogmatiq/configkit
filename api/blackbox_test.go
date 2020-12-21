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
				c.Identity("<app-1>", "<app-key-1>")

				c.RegisterAggregate(&AggregateMessageHandler{
					ConfigureFunc: func(c dogma.AggregateConfigurer) {
						c.Identity("<aggregate>", "<aggregate-key>")
						c.ConsumesCommandType(MessageC{})
						c.ProducesEventType(MessageE{})
					},
				})

				c.RegisterProcess(&ProcessMessageHandler{
					ConfigureFunc: func(c dogma.ProcessConfigurer) {
						c.Identity("<process>", "<process-key>")
						c.ConsumesEventType(MessageE{})
						c.ProducesCommandType(MessageC{})
						c.SchedulesTimeoutType(MessageT{})
					},
				})
			},
		}

		app2 = &Application{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("<app-2>", "<app-key-2>")

				c.RegisterIntegration(&IntegrationMessageHandler{
					ConfigureFunc: func(c dogma.IntegrationConfigurer) {
						c.Identity("<integration>", "<integration-key>")
						c.ConsumesCommandType(MessageI{})
						c.ProducesEventType(MessageJ{})
					},
				})

				c.RegisterProjection(&ProjectionMessageHandler{
					ConfigureFunc: func(c dogma.ProjectionConfigurer) {
						c.Identity("<projection>", "<projection-key>")
						c.ConsumesEventType(MessageE{})
						c.ConsumesEventType(MessageJ{})
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
