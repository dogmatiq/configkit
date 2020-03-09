package discovery_test

import (
	"context"
	"time"

	. "github.com/dogmatiq/configkit/api/discovery"
	"github.com/dogmatiq/configkit/api/fixtures" // can't dot-import due to conflict
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ ClientObserver = (*ClientObserverSet)(nil)

var _ = Describe("type ClientObserverSet", func() {
	var (
		set              *ClientObserverSet
		obs1, obs2       *fixtures.ClientObserver
		client1, client2 *Client
	)

	BeforeEach(func() {
		set = &ClientObserverSet{}

		obs1 = &fixtures.ClientObserver{}
		obs2 = &fixtures.ClientObserver{}

		client1 = &Client{}
		client2 = &Client{}
	})

	Describe("func NewClientObserverSet()", func() {
		It("returns a set containing the given observers", func() {
			set.RegisterClientObserver(obs1)
			set.RegisterClientObserver(obs2)

			Expect(
				NewClientObserverSet(obs1, obs2),
			).To(Equal(set))
		})
	})

	Describe("func ClientConnected()", func() {
		It("notifies the observers about the connection", func() {
			var observers []ClientObserver

			obs1.ClientConnectedFunc = func(c *Client) {
				defer GinkgoRecover()
				Expect(c).To(BeIdenticalTo(client1))

				observers = append(observers, obs1)
			}

			obs2.ClientConnectedFunc = func(c *Client) {
				defer GinkgoRecover()
				Expect(c).To(BeIdenticalTo(client1))

				observers = append(observers, obs2)
			}

			set.RegisterClientObserver(obs1)
			set.RegisterClientObserver(obs2)
			set.ClientConnected(client1)

			Expect(observers).To(ConsistOf(obs1, obs2))
		})

		It("does nothing if the client is already connected", func() {
			set.RegisterClientObserver(obs1)
			set.ClientConnected(client1)

			obs1.ClientConnectedFunc = func(c *Client) {
				defer GinkgoRecover()
				Fail("observer unexpectedly notified of the same connection")
			}

			set.ClientConnected(client1)
		})
	})

	Describe("func ClientDisconnected()", func() {
		It("notifies the observers about the disconnection", func() {
			var observers []ClientObserver

			obs1.ClientDisconnectedFunc = func(c *Client) {
				defer GinkgoRecover()
				Expect(c).To(BeIdenticalTo(client1))

				observers = append(observers, obs1)
			}

			obs2.ClientDisconnectedFunc = func(c *Client) {
				defer GinkgoRecover()
				Expect(c).To(BeIdenticalTo(client1))

				observers = append(observers, obs2)
			}

			set.RegisterClientObserver(obs1)
			set.RegisterClientObserver(obs2)
			set.ClientConnected(client1)
			set.ClientDisconnected(client1)

			Expect(observers).To(ConsistOf(obs1, obs2))
		})

		It("does nothing if the client is not in the registry", func() {
			set.RegisterClientObserver(obs1)

			obs1.ClientDisconnectedFunc = func(c *Client) {
				defer GinkgoRecover()
				Fail("observer unexpectedly notified of unknown client")
			}

			set.ClientDisconnected(client1)
		})
	})

	Describe("func RegisterClientObserver()", func() {
		It("notifies the observer about existing connections", func() {
			set.ClientConnected(client1)
			set.ClientConnected(client2)

			var clients []*Client

			obs1.ClientConnectedFunc = func(c *Client) {
				clients = append(clients, c)
			}

			set.RegisterClientObserver(obs1)

			Expect(clients).To(ConsistOf(client1, client2))
		})

		It("does nothing if the observer is already registered", func() {
			set.ClientConnected(client1)
			set.RegisterClientObserver(obs1)

			obs1.ClientConnectedFunc = func(c *Client) {
				defer GinkgoRecover()
				Fail("observer unexpectedly notified when registered twice")
			}

			set.RegisterClientObserver(obs1)
		})
	})

	Describe("func UnregisterClientObserver()", func() {
		It("synthesizes a disconnection notification for existing connections", func() {
			set.ClientConnected(client1)
			set.ClientConnected(client2)
			set.RegisterClientObserver(obs1)

			var clients []*Client

			obs1.ClientDisconnectedFunc = func(c *Client) {
				clients = append(clients, c)
			}

			set.UnregisterClientObserver(obs1)

			Expect(clients).To(ConsistOf(client1, client2))
		})

		It("prevents the observer from receiving further notifications", func() {
			obs1.ClientConnectedFunc = func(c *Client) {
				defer GinkgoRecover()
				Fail("unregistered observer unexpectedly notified")
			}

			set.RegisterClientObserver(obs1)
			set.UnregisterClientObserver(obs1)
			set.ClientConnected(client1)
		})

		It("does nothing if the observer is not already registered", func() {
			set.ClientConnected(client1)

			obs1.ClientDisconnectedFunc = func(c *Client) {
				defer GinkgoRecover()
				Fail("observer unexpectedly notified when not registered")
			}

			set.UnregisterClientObserver(obs1)
		})
	})
})

var _ ClientObserver = (*ClientExecutor)(nil)

var _ = Describe("type ClientExecutor", func() {
	var (
		ctx    context.Context
		cancel func()
		exec   *ClientExecutor
		client *Client
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 250*time.Millisecond)

		exec = &ClientExecutor{
			Task: func(context.Context, *Client) {},
		}

		client = &Client{}
	})

	AfterEach(func() {
		cancel()
	})

	Describe("func ClientConnected()", func() {
		It("starts a goroutine for the given client", func() {
			barrier := make(chan struct{})

			exec.Task = func(_ context.Context, c *Client) {
				defer GinkgoRecover()
				defer close(barrier)

				Expect(c).To(Equal(client))
			}

			exec.ClientConnected(client)
			defer exec.ClientDisconnected(client)

			select {
			case <-barrier:
			case <-ctx.Done():
				Expect(ctx.Err()).ShouldNot(HaveOccurred())
			}
		})

		It("does not panic if the client is already connected", func() {
			exec.ClientConnected(client)
			defer exec.ClientDisconnected(client)

			exec.ClientConnected(client)
		})
	})

	Describe("func ClientDisconnected()", func() {
		It("cancels the context associated with the goroutine and waits for the function to finish", func() {
			barrier := make(chan struct{})

			exec.Task = func(funcCtx context.Context, c *Client) {
				defer GinkgoRecover()
				defer close(barrier)

				select {
				case <-funcCtx.Done():
					// ok
				case <-ctx.Done():
					Expect(ctx.Err()).ShouldNot(HaveOccurred())
				}
			}

			exec.ClientConnected(client)
			exec.ClientDisconnected(client)

			select {
			case <-barrier:
			case <-ctx.Done():
				Expect(ctx.Err()).ShouldNot(HaveOccurred())
			}
		})

		It("does not panic if the client is already disconnected", func() {
			exec.ClientDisconnected(client)
		})
	})
})
