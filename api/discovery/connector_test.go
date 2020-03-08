package discovery_test

import (
	"context"
	"net"
	"time"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/api"
	. "github.com/dogmatiq/configkit/api/discovery"
	"github.com/dogmatiq/configkit/api/fixtures" // can't dot-import due to conflict
	"github.com/dogmatiq/dodeca/logging"
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/dogma/fixtures"
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
		obs       *fixtures.ClientObserver
		connector *Connector
		target    *Target
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 250*time.Millisecond)

		var err error
		listener, err = net.Listen("tcp", ":")
		Expect(err).ShouldNot(HaveOccurred())

		gserver = grpc.NewServer()

		obs = &fixtures.ClientObserver{
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

	Describe("Watch()", func() {
		Context("when dialing fails", func() {
			BeforeEach(func() {
				listener.Close()
				target.Options = append(target.Options, grpc.WithBlock())
			})

			It("does not notify the observer", func() {
				err := connector.Watch(ctx, target)
				Expect(err).To(Equal(context.DeadlineExceeded))
			})
		})

		Context("when the server is down", func() {
			BeforeEach(func() {
				listener.Close()
			})

			It("does not notify the observer", func() {
				err := connector.Watch(ctx, target)
				Expect(err).To(Equal(context.DeadlineExceeded))
			})
		})

		Context("when the server does not implement the config API", func() {
			It("does not notify the observer", func() {
				err := connector.Watch(ctx, target)
				Expect(err).To(Equal(context.DeadlineExceeded))
			})
		})

		Context("when the server implements the config API", func() {
			BeforeEach(func() {
				app := &Application{
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

				watchCtx, cancelWatch := context.WithCancel(ctx)
				defer cancelWatch()

				go connector.Watch(watchCtx, target)

				select {
				case <-connected:
				case <-ctx.Done():
					Expect(ctx.Err()).ShouldNot(HaveOccurred())
				}

				cancelWatch()

				select {
				case <-disconnected:
				case <-ctx.Done():
					Expect(ctx.Err()).ShouldNot(HaveOccurred())
				}
			})

			It("connects to the server", func() {
				watchCtx, cancelWatch := context.WithCancel(ctx)
				defer cancelWatch()

				obs.ClientConnectedFunc = func(c *Client) {
					defer GinkgoRecover()
					defer cancelWatch()

					idents, err := c.ListApplicationIdentities(ctx)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(idents).To(ConsistOf(
						configkit.MustNewIdentity("<app>", "<app-key>"),
					))
				}

				obs.ClientDisconnectedFunc = nil

				err := connector.Watch(watchCtx, target)
				Expect(err).To(Equal(context.Canceled))
			})

			It("provides the underlying connection", func() {
				watchCtx, cancelWatch := context.WithCancel(ctx)
				defer cancelWatch()

				obs.ClientConnectedFunc = func(c *Client) {
					defer GinkgoRecover()
					defer cancelWatch()

					client := api.NewClient(c.Connection)
					idents, err := client.ListApplicationIdentities(ctx)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(idents).To(ConsistOf(
						configkit.MustNewIdentity("<app>", "<app-key>"),
					))
				}

				obs.ClientDisconnectedFunc = nil

				err := connector.Watch(watchCtx, target)
				Expect(err).To(Equal(context.Canceled))
			})

			It("notifies of a disconnection if the server goes offline", func() {
				watchCtx, cancelWatch := context.WithCancel(ctx)
				defer cancelWatch()

				obs.ClientConnectedFunc = func(c *Client) {
					gserver.Stop()
				}

				obs.ClientDisconnectedFunc = func(c *Client) {
					cancel()
				}

				err := connector.Watch(watchCtx, target)
				Expect(err).To(Equal(context.Canceled))
			})
		})
	})
})
