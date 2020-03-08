package discovery_test

import (
	"context"
	"net"
	"time"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/api"
	. "github.com/dogmatiq/configkit/api/discovery"
	dfixtures "github.com/dogmatiq/configkit/api/discovery/fixtures" // can't dot-import due to conflict
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures" // can't dot-import due to conflict
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("type Inspector", func() {
	var (
		ctx        context.Context
		cancel     func()
		listener   net.Listener
		gserver    *grpc.Server
		gconn      *grpc.ClientConn
		cfg1, cfg2 configkit.Application
		obs        *dfixtures.ApplicationObserver
		inspector  *Inspector
		client     *Client
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 250*time.Millisecond)

		var err error
		listener, err = net.Listen("tcp", ":")
		Expect(err).ShouldNot(HaveOccurred())

		gserver = grpc.NewServer()

		cfg1 = configkit.FromApplication(&fixtures.Application{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("<app-1>", "<app-key-1>")
			},
		})

		cfg2 = configkit.FromApplication(&fixtures.Application{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("<app-2>", "<app-key-2>")
			},
		})

		api.RegisterServer(gserver, cfg1, cfg2)

		go gserver.Serve(listener)

		obs = &dfixtures.ApplicationObserver{
			ApplicationAvailableFunc: func(a *Application) {
				defer GinkgoRecover()
				Fail("unexpected application available notification")
			},
			ApplicationUnavailableFunc: func(a *Application) {
				defer GinkgoRecover()
				Fail("unexpected application unavailable notification")
			},
		}

		gconn, err = grpc.DialContext(ctx, listener.Addr().String(), grpc.WithInsecure())
		Expect(err).ShouldNot(HaveOccurred())
		inspector = &Inspector{
			Observer: obs,
		}

		client = &Client{
			Client: api.NewClient(gconn),
			Target: &Target{
				Name: listener.Addr().String(),
				Options: []grpc.DialOption{
					grpc.WithInsecure(),
				},
			},
		}
	})

	AfterEach(func() {
		if listener != nil {
			listener.Close()
		}

		if gserver != nil {
			gserver.Stop()
		}

		if gconn != nil {
			gconn.Close()
		}

		cancel()
	})

	Describe("Inspect()", func() {
		It("notifies the observer", func() {
			inspectCtx, cancelInspect := context.WithCancel(ctx)
			defer cancelInspect()

			var available, unavailable []*Application

			obs.ApplicationAvailableFunc = func(a *Application) {
				available = append(available, a)

				if len(available) == 2 {
					cancelInspect()
				}
			}

			obs.ApplicationUnavailableFunc = func(a *Application) {
				unavailable = append(unavailable, a)
			}

			err := inspector.Inspect(inspectCtx, client)
			Expect(err).To(Equal(context.Canceled))
			Expect(available).To(HaveLen(2))
			Expect(unavailable).To(ConsistOf(available))

			Expect(
				configkit.IsApplicationEqual(available[0], cfg1),
			).To(BeTrue())

			Expect(
				configkit.IsApplicationEqual(available[1], cfg2),
			).To(BeTrue())
		})

		It("does not notify the observer if the application is ignored", func() {
			inspector.Ignore = func(a configkit.Application) bool {
				return a.Identity().Key == "<app-key-1>"
			}

			inspectCtx, cancelInspect := context.WithCancel(ctx)
			defer cancelInspect()

			var available []*Application

			obs.ApplicationAvailableFunc = func(a *Application) {
				available = append(available, a)
			}

			obs.ApplicationUnavailableFunc = nil

			err := inspector.Inspect(inspectCtx, client)
			Expect(err).To(Equal(context.DeadlineExceeded))
			Expect(available).To(HaveLen(1))

			Expect(
				configkit.IsApplicationEqual(available[0], cfg2),
			).To(BeTrue())
		})

		It("returns immediately of all applications are ignored", func() {
			inspector.Ignore = func(a configkit.Application) bool {
				return true
			}

			err := inspector.Inspect(ctx, client)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("returns an error if the query fails", func() {
			gserver.Stop()

			err := inspector.Inspect(ctx, client)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.Unavailable))
		})
	})
})
