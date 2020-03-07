package discovery_test

import (
	"context"
	"net"
	"time"

	"github.com/dogmatiq/dodeca/logging"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/api"
	. "github.com/dogmatiq/configkit/api/discovery"
	"github.com/dogmatiq/configkit/api/discovery/fixtures" // can't dot-import due to conflict
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/dogma/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
)

var _ TargetObserver = (*Connector)(nil)

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

	Context("when dialing fails", func() {
		BeforeEach(func() {
			listener.Close()
			target.Options = append(target.Options, grpc.WithBlock())
		})

		It("does not notify the observer", func() {
			connector.TargetAvailable(target)
			defer connector.TargetUnavailable(target)
			<-ctx.Done() // wait out the timeout
		})
	})

	Context("when the server is down", func() {
		BeforeEach(func() {
			listener.Close()
		})

		It("does not notify the observer", func() {
			connector.TargetAvailable(target)
			defer connector.TargetUnavailable(target)
			<-ctx.Done() // wait out the timeout
		})
	})

	Context("when the server is does not implement the config API", func() {
		It("does not notify the observer", func() {
			connector.TargetAvailable(target)
			defer connector.TargetUnavailable(target)
			<-ctx.Done() // wait out the timeout
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
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			barrier := make(chan bool)

			obs.ClientConnectedFunc = func(c *Client) {
				defer GinkgoRecover()

				Expect(c.Target).To(Equal(target))
				barrier <- true
			}

			obs.ClientDisconnectedFunc = func(c *Client) {
				defer GinkgoRecover()

				Expect(c.Target).To(Equal(target))
				barrier <- false
			}

			connector.TargetAvailable(target)
			defer connector.TargetUnavailable(target)

			select {
			case connect := <-barrier:
				Expect(connect).To(BeTrue())
			case <-ctx.Done():
				Expect(ctx.Err()).ShouldNot(HaveOccurred())
			}

			connector.TargetUnavailable(target)

			select {
			case connect := <-barrier:
				Expect(connect).To(BeFalse())
			case <-ctx.Done():
				Expect(ctx.Err()).ShouldNot(HaveOccurred())
			}
		})

		It("connects to the server", func() {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			obs.ClientConnectedFunc = func(c *Client) {
				defer GinkgoRecover()
				defer cancel()

				idents, err := c.ListApplicationIdentities(ctx)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(idents).To(ConsistOf(
					configkit.MustNewIdentity("<app>", "<app-key>"),
				))
			}

			obs.ClientDisconnectedFunc = nil

			connector.TargetAvailable(target)
			defer connector.TargetUnavailable(target)
			<-ctx.Done()
		})

		It("notifies of a disconnection if the server goes offline", func() {
			barrier := make(chan struct{})

			obs.ClientConnectedFunc = func(c *Client) {
				gserver.Stop()
			}

			obs.ClientDisconnectedFunc = func(c *Client) {
				barrier <- struct{}{}
			}

			connector.TargetAvailable(target)
			defer connector.TargetUnavailable(target)

			select {
			case <-barrier:
			case <-ctx.Done():
				Expect(ctx.Err()).ShouldNot(HaveOccurred())
			}
		})
	})

	Describe("func TargetAvailable()", func() {
		It("does not panic if the target is already available", func() {
			connector.TargetAvailable(target)
			defer connector.TargetUnavailable(target)

			connector.TargetAvailable(target)
		})
	})

	Describe("func TargetUnavailable()", func() {
		It("does not panic if the target is already unavailable", func() {
			connector.TargetUnavailable(target)
		})
	})
})
