package discovery_test

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc/status"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/api"
	. "github.com/dogmatiq/configkit/api/discovery"
	dfixtures "github.com/dogmatiq/configkit/api/discovery/fixtures" // can't dot-import due to conflict
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var _ = Describe("type Inspector", func() {
	var (
		ctx       context.Context
		cancel    func()
		listener  net.Listener
		gserver   *grpc.Server
		gconn     *grpc.ClientConn
		cfg       configkit.Application
		obs       *dfixtures.ApplicationObserver
		inspector *Inspector
		client    *Client
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 250*time.Millisecond)

		var err error
		listener, err = net.Listen("tcp", ":")
		Expect(err).ShouldNot(HaveOccurred())

		gserver = grpc.NewServer()

		app := &fixtures.Application{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("<app>", "<app-key>")
			},
		}

		cfg = configkit.FromApplication(app)
		api.RegisterServer(gserver, cfg)

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
			available := make(chan struct{})
			unavailable := make(chan struct{})

			obs.ApplicationAvailableFunc = func(a *Application) {
				defer GinkgoRecover()

				Expect(a.Client).To(Equal(client))
				close(available)
			}

			obs.ApplicationUnavailableFunc = func(a *Application) {
				defer GinkgoRecover()

				Expect(a.Client).To(Equal(client))
				close(unavailable)
			}

			inspectCtx, cancelInspect := context.WithCancel(ctx)
			defer cancelInspect()

			go inspector.Inspect(inspectCtx, client)

			select {
			case <-available:
			case <-ctx.Done():
				Expect(ctx.Err()).ShouldNot(HaveOccurred())
			}

			cancelInspect()

			select {
			case <-unavailable:
			case <-ctx.Done():
				Expect(ctx.Err()).ShouldNot(HaveOccurred())
			}
		})

		It("does not notify the observer if the application is ignored", func() {
			inspector.Ignore = []string{"<app-key>"}

			err := inspector.Inspect(ctx, client)
			Expect(err).To(Equal(context.DeadlineExceeded))
		})

		It("inspects the application", func() {
			inspectCtx, cancelInspect := context.WithCancel(ctx)
			defer cancelInspect()

			obs.ApplicationAvailableFunc = func(a *Application) {
				defer GinkgoRecover()
				defer cancelInspect()

				Expect(configkit.IsApplicationEqual(
					a,
					cfg,
				))
			}

			obs.ApplicationUnavailableFunc = nil

			err := inspector.Inspect(inspectCtx, client)
			Expect(err).To(Equal(context.Canceled))
		})

		It("returns an error if the query fails", func() {
			gserver.Stop()

			err := inspector.Inspect(ctx, client)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.Unavailable))
		})
	})
})
