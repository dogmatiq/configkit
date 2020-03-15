package discovery_test

import (
	"context"
	"sync"
	"time"

	. "github.com/dogmatiq/configkit/api/discovery"
	"github.com/dogmatiq/configkit/api/fixtures" // can't dot-import due to conflict
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ ApplicationObserver = (*ApplicationObserverSet)(nil)

var _ = Describe("type ApplicationObserverSet", func() {
	var (
		set        *ApplicationObserverSet
		obs1, obs2 *fixtures.ApplicationObserver
		app1, app2 *Application
	)

	BeforeEach(func() {
		set = &ApplicationObserverSet{}

		obs1 = &fixtures.ApplicationObserver{}
		obs2 = &fixtures.ApplicationObserver{}

		app1 = &Application{Client: &Client{Target: &Target{Name: "<target-1>"}}}
		app2 = &Application{Client: &Client{Target: &Target{Name: "<target-2>"}}}
	})

	Describe("func NewApplicationObserverSet()", func() {
		It("returns a set containing the given observers", func() {
			set.RegisterApplicationObserver(obs1)
			set.RegisterApplicationObserver(obs2)

			Expect(
				NewApplicationObserverSet(obs1, obs2),
			).To(Equal(set))
		})
	})

	Describe("func ApplicationAvailable()", func() {
		It("notifies the observers about the application availability", func() {
			var m sync.Mutex
			var observers []ApplicationObserver

			obs1.ApplicationAvailableFunc = func(a *Application) {
				defer GinkgoRecover()
				Expect(a).To(BeIdenticalTo(app1))

				m.Lock()
				defer m.Unlock()
				observers = append(observers, obs1)
			}

			obs2.ApplicationAvailableFunc = func(a *Application) {
				defer GinkgoRecover()
				Expect(a).To(BeIdenticalTo(app1))

				m.Lock()
				defer m.Unlock()
				observers = append(observers, obs2)
			}

			set.RegisterApplicationObserver(obs1)
			set.RegisterApplicationObserver(obs2)
			set.ApplicationAvailable(app1)

			Expect(observers).To(ConsistOf(obs1, obs2))
		})

		It("does nothing if the application is already available", func() {
			set.RegisterApplicationObserver(obs1)
			set.ApplicationAvailable(app1)

			obs1.ApplicationAvailableFunc = func(a *Application) {
				defer GinkgoRecover()
				Fail("observer unexpectedly notified of the same application")
			}

			set.ApplicationAvailable(app1)
		})
	})

	Describe("func ApplicationUnavailable()", func() {
		It("notifies the observers about the application unavailability", func() {
			var m sync.Mutex
			var observers []ApplicationObserver

			obs1.ApplicationUnavailableFunc = func(a *Application) {
				defer GinkgoRecover()
				Expect(a).To(BeIdenticalTo(app1))

				m.Lock()
				defer m.Unlock()
				observers = append(observers, obs1)
			}

			obs2.ApplicationUnavailableFunc = func(a *Application) {
				defer GinkgoRecover()
				Expect(a).To(BeIdenticalTo(app1))

				m.Lock()
				defer m.Unlock()
				observers = append(observers, obs2)
			}

			set.RegisterApplicationObserver(obs1)
			set.RegisterApplicationObserver(obs2)
			set.ApplicationAvailable(app1)
			set.ApplicationUnavailable(app1)

			Expect(observers).To(ConsistOf(obs1, obs2))
		})

		It("does nothing if the application is not in the registry", func() {
			set.RegisterApplicationObserver(obs1)

			obs1.ApplicationUnavailableFunc = func(a *Application) {
				defer GinkgoRecover()
				Fail("observer unexpectedly notified of unknown application")
			}

			set.ApplicationUnavailable(app1)
		})
	})

	Describe("func RegisterApplicationObserver()", func() {
		It("notifies the observer about existing applications", func() {
			set.ApplicationAvailable(app1)
			set.ApplicationAvailable(app2)

			var apps []*Application

			obs1.ApplicationAvailableFunc = func(a *Application) {
				apps = append(apps, a)
			}

			set.RegisterApplicationObserver(obs1)

			Expect(apps).To(ConsistOf(app1, app2))
		})

		It("does nothing if the observer is already registered", func() {
			set.ApplicationAvailable(app1)
			set.RegisterApplicationObserver(obs1)

			obs1.ApplicationAvailableFunc = func(a *Application) {
				defer GinkgoRecover()
				Fail("observer unexpectedly notified when registered twice")
			}

			set.RegisterApplicationObserver(obs1)
		})
	})

	Describe("func UnregisterApplicationObserver()", func() {
		It("synthesizes an unavailable notification for existing applications", func() {
			set.ApplicationAvailable(app1)
			set.ApplicationAvailable(app2)
			set.RegisterApplicationObserver(obs1)

			var apps []*Application

			obs1.ApplicationUnavailableFunc = func(a *Application) {
				apps = append(apps, a)
			}

			set.UnregisterApplicationObserver(obs1)

			Expect(apps).To(ConsistOf(app1, app2))
		})

		It("prevents the observer from receiving further notifications", func() {
			obs1.ApplicationAvailableFunc = func(a *Application) {
				defer GinkgoRecover()
				Fail("unregistered observer unexpectedly notified")
			}

			set.RegisterApplicationObserver(obs1)
			set.UnregisterApplicationObserver(obs1)
			set.ApplicationAvailable(app1)
		})

		It("does nothing if the observer is not already registered", func() {
			set.ApplicationAvailable(app1)

			obs1.ApplicationUnavailableFunc = func(a *Application) {
				defer GinkgoRecover()
				Fail("observer unexpectedly notified when not registered")
			}

			set.UnregisterApplicationObserver(obs1)
		})
	})
})

var _ ApplicationObserver = (*ApplicationExecutor)(nil)

var _ = Describe("type ApplicationExecutor", func() {
	var (
		ctx    context.Context
		cancel func()
		exec   *ApplicationExecutor
		app    *Application
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 250*time.Millisecond)

		exec = &ApplicationExecutor{
			Task: func(context.Context, *Application) {},
		}

		app = &Application{Client: &Client{Target: &Target{Name: "<target>"}}}
	})

	AfterEach(func() {
		cancel()
	})

	Describe("func ApplicationAvailable()", func() {
		It("starts a goroutine for the given application", func() {
			barrier := make(chan struct{})

			exec.Task = func(_ context.Context, a *Application) {
				defer GinkgoRecover()
				defer close(barrier)

				Expect(a).To(Equal(app))
			}

			exec.ApplicationAvailable(app)
			defer exec.ApplicationUnavailable(app)

			select {
			case <-barrier:
			case <-ctx.Done():
				Expect(ctx.Err()).ShouldNot(HaveOccurred())
			}
		})

		It("does not panic if the application is already available", func() {
			exec.ApplicationAvailable(app)
			defer exec.ApplicationUnavailable(app)

			exec.ApplicationAvailable(app)
		})
	})

	Describe("func ApplicationUnavailable()", func() {
		It("cancels the context associated with the goroutine and waits for the function to finish", func() {
			barrier := make(chan struct{})

			exec.Task = func(funcCtx context.Context, a *Application) {
				defer GinkgoRecover()
				defer close(barrier)

				select {
				case <-funcCtx.Done():
					// ok
				case <-ctx.Done():
					Expect(ctx.Err()).ShouldNot(HaveOccurred())
				}
			}

			exec.ApplicationAvailable(app)
			exec.ApplicationUnavailable(app)

			select {
			case <-barrier:
			case <-ctx.Done():
				Expect(ctx.Err()).ShouldNot(HaveOccurred())
			}
		})

		It("does not panic if the application is already unavailable", func() {
			exec.ApplicationUnavailable(app)
		})
	})
})
