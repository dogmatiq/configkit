package discovery_test

import (
	"context"
	"net"
	"time"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/api"
	. "github.com/dogmatiq/configkit/api/discovery"
	apifixtures "github.com/dogmatiq/configkit/api/fixtures" // can't dot-import due to conflict
	"github.com/dogmatiq/dodeca/logging"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures" // can't dot-import due to conflict
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
)

var _ = Describe("type Connector", func() {
	var (
		ctx       context.Context
		cancel    func()
		listener  net.Listener
		gserver   *grpc.Server
		obs       *apifixtures.ClientObserver
		connector *Connector
		target    *Target
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 250*time.Millisecond)

		var err error
		listener, err = net.Listen("tcp", ":")
		Expect(err).ShouldNot(HaveOccurred())

		gserver = grpc.NewServer()

		obs = &apifixtures.ClientObserver{
			ClientConnectedFunc: func(c *Client) {
				defer GinkgoRecover()
				Fail("unexpected client connected notification")
			},
			ClientDisconnectedFunc: func(c *Client) {
				defer GinkgoRecover()
				Fail("unexpected client disconnected notification")
			},
		}

		connector = &Connector{
			Observer: obs,
			Logger:   logging.DiscardLogger{},
		}

		target = &Target{
			Name: listener.Addr().String(),
			Options: []grpc.DialOption{
				grpc.WithInsecure(),
			},
		}
	})

	JustBeforeEach(func() {
		go gserver.Serve(listener)
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

	Describe("Run()", func() {
		Context("when dialing fails", func() {
			BeforeEach(func() {
				listener.Close()
				target.Options = append(target.Options, grpc.WithBlock())
			})

			It("does not notify the observer", func() {
				err := connector.Run(ctx, target)
				Expect(err).To(Equal(context.DeadlineExceeded))
			})
		})

		Context("when the server is down", func() {
			BeforeEach(func() {
				listener.Close()
			})

			It("does not notify the observer", func() {
				err := connector.Run(ctx, target)
				Expect(err).To(Equal(context.DeadlineExceeded))
			})
		})

		Context("when the server does not implement the config API", func() {
			It("does not notify the observer", func() {
				err := connector.Run(ctx, target)
				Expect(err).To(Equal(context.DeadlineExceeded))
			})
		})

		Context("when the target is ignored", func() {
			BeforeEach(func() {
				connector.Ignore = func(t *Target) bool {
					return t == target
				}
			})

			It("returns immediately", func() {
				err := connector.Run(ctx, target)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when the server implements the config API", func() {
			BeforeEach(func() {
				app := &fixtures.Application{
					ConfigureFunc: func(c dogma.ApplicationConfigurer) {
						c.Identity("<app>", "<app-key>")
					},
				}

				cfg := configkit.FromApplication(app)
				api.RegisterServer(gserver, cfg)
			})

			It("notifies the observer", func() {
				connected := make(chan struct{})
				disconnected := make(chan struct{})

				obs.ClientConnectedFunc = func(c *Client) {
					defer GinkgoRecover()

					Expect(c.Target).To(Equal(target))
					close(connected)
				}

				obs.ClientDisconnectedFunc = func(c *Client) {
					defer GinkgoRecover()

					Expect(c.Target).To(Equal(target))
					close(disconnected)
				}

				runCtx, cancelRun := context.WithCancel(ctx)
				defer cancelRun()

				go connector.Run(runCtx, target)

				select {
				case <-connected:
				case <-ctx.Done():
					Expect(ctx.Err()).ShouldNot(HaveOccurred())
				}

				cancelRun()

				select {
				case <-disconnected:
				case <-ctx.Done():
					Expect(ctx.Err()).ShouldNot(HaveOccurred())
				}
			})

			It("connects to the server", func() {
				runCtx, cancelRun := context.WithCancel(ctx)
				defer cancelRun()

				obs.ClientConnectedFunc = func(c *Client) {
					defer GinkgoRecover()
					defer cancelRun()

					idents, err := c.ListApplicationIdentities(ctx)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(idents).To(ConsistOf(
						configkit.MustNewIdentity("<app>", "<app-key>"),
					))
				}

				obs.ClientDisconnectedFunc = nil

				err := connector.Run(runCtx, target)
				Expect(err).To(Equal(context.Canceled))
			})

			It("provides the underlying connection", func() {
				runCtx, cancelRun := context.WithCancel(ctx)
				defer cancelRun()

				obs.ClientConnectedFunc = func(c *Client) {
					defer GinkgoRecover()
					defer cancelRun()

					client := api.NewClient(c.Connection)
					idents, err := client.ListApplicationIdentities(ctx)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(idents).To(ConsistOf(
						configkit.MustNewIdentity("<app>", "<app-key>"),
					))
				}

				obs.ClientDisconnectedFunc = nil

				err := connector.Run(runCtx, target)
				Expect(err).To(Equal(context.Canceled))
			})

			It("notifies of a disconnection if the server goes offline", func() {
				runCtx, cancelRun := context.WithCancel(ctx)
				defer cancelRun()

				obs.ClientConnectedFunc = func(c *Client) {
					gserver.Stop()
				}

				obs.ClientDisconnectedFunc = func(c *Client) {
					cancel()
				}

				err := connector.Run(runCtx, target)
				Expect(err).To(Equal(context.Canceled))
			})
		})
	})
})
