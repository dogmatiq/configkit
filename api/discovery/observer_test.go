package discovery_test

import (
	"sync"

	"github.com/dogmatiq/configkit/api"
	. "github.com/dogmatiq/configkit/api/discovery"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type ObserverSet", func() {
	var (
		set              *ObserverSet
		obs1, obs2       *observer
		client1, client2 *api.Client
	)

	BeforeEach(func() {
		set = &ObserverSet{}

		obs1 = &observer{}
		obs2 = &observer{}

		client1 = &api.Client{}
		client2 = &api.Client{}
	})

	Describe("func Connected()", func() {
		It("notifies the observers about the connection", func() {
			var observers []Observer

			obs1.ConnectedFunc = func(c *api.Client) {
				defer GinkgoRecover()
				Expect(c).To(BeIdenticalTo(client1))

				observers = append(observers, obs1)
			}

			obs2.ConnectedFunc = func(c *api.Client) {
				defer GinkgoRecover()
				Expect(c).To(BeIdenticalTo(client1))

				observers = append(observers, obs2)
			}

			set.Add(obs1)
			set.Add(obs2)
			set.Connected(client1)

			Expect(observers).To(ConsistOf(obs1, obs2))
		})

		It("does nothing if the client is already connected", func() {
			set.Add(obs1)
			set.Connected(client1)

			obs1.ConnectedFunc = func(c *api.Client) {
				defer GinkgoRecover()
				Fail("observer unexpectedly notified of the same connection")
			}

			set.Connected(client1)
		})
	})

	Describe("func Disconnected()", func() {
		It("notifies the observers about the disconnection", func() {
			var observers []Observer

			obs1.DisconnectedFunc = func(c *api.Client) {
				defer GinkgoRecover()
				Expect(c).To(BeIdenticalTo(client1))

				observers = append(observers, obs1)
			}

			obs2.DisconnectedFunc = func(c *api.Client) {
				defer GinkgoRecover()
				Expect(c).To(BeIdenticalTo(client1))

				observers = append(observers, obs2)
			}

			set.Add(obs1)
			set.Add(obs2)
			set.Connected(client1)
			set.Disconnected(client1)

			Expect(observers).To(ConsistOf(obs1, obs2))
		})

		It("does nothing if the target is not in the registry", func() {
			set.Add(obs1)

			obs1.DisconnectedFunc = func(c *api.Client) {
				defer GinkgoRecover()
				Fail("observer unexpectedly notified of unkonwn target")
			}

			set.Disconnected(client1)
		})
	})

	Describe("func Add()", func() {
		It("notifies the observer about existing connections", func() {
			set.Connected(client1)
			set.Connected(client2)

			var clients []*api.Client

			obs1.ConnectedFunc = func(c *api.Client) {
				clients = append(clients, c)
			}

			set.Add(obs1)

			Expect(clients).To(ConsistOf(client1, client2))
		})

		It("does nothing if the observer is already registered", func() {
			set.Connected(client1)
			set.Add(obs1)

			obs1.ConnectedFunc = func(c *api.Client) {
				defer GinkgoRecover()
				Fail("observer unexpectedly notified when registered twice")
			}

			set.Add(obs1)
		})
	})

	Describe("func Remove()", func() {
		It("synthesizes a disconnection notification for existing connections", func() {
			set.Connected(client1)
			set.Connected(client2)
			set.Add(obs1)

			var clients []*api.Client

			obs1.DisconnectedFunc = func(c *api.Client) {
				clients = append(clients, c)
			}

			set.Remove(obs1)

			Expect(clients).To(ConsistOf(client1, client2))
		})

		It("prevents the observer from receiving further notifications", func() {
			obs1.ConnectedFunc = func(c *api.Client) {
				defer GinkgoRecover()
				Fail("unregistered observer unexpectedly notified")
			}

			set.Add(obs1)
			set.Remove(obs1)
			set.Connected(client1)
		})

		It("does nothing if the observer is not already registered", func() {
			set.Connected(client1)

			obs1.DisconnectedFunc = func(c *api.Client) {
				defer GinkgoRecover()
				Fail("observer unexpectedly notified when not registered")
			}

			set.Remove(obs1)
		})
	})
})

type observer struct {
	m                sync.Mutex
	ConnectedFunc    func(*api.Client)
	DisconnectedFunc func(*api.Client)
}

func (o *observer) Connected(c *api.Client) {
	if o.ConnectedFunc != nil {
		o.m.Lock()
		defer o.m.Unlock()
		o.ConnectedFunc(c)
	}
}

func (o *observer) Disconnected(c *api.Client) {
	if o.DisconnectedFunc != nil {
		o.m.Lock()
		defer o.m.Unlock()
		o.DisconnectedFunc(c)
	}
}
