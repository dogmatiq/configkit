package discovery_test

import (
	. "github.com/dogmatiq/configkit/api/discovery"
	"github.com/dogmatiq/configkit/api/discovery/fixtures" // can't dot-import due to conflict
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
