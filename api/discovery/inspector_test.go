package discovery_test

import (
	"context"
	"errors"
	"time"

	"github.com/dogmatiq/linger/backoff"

	"github.com/dogmatiq/configkit"
	. "github.com/dogmatiq/configkit/api/discovery"
	apifixtures "github.com/dogmatiq/configkit/api/fixtures" // can't dot-import due to conflict
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures" // can't dot-import due to conflict
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Inspector", func() {
	var (
		ctx        context.Context
		cancel     func()
		cfg1, cfg2 configkit.Application
		obs        *apifixtures.ApplicationObserver
		inspector  *Inspector
		apiClient  *apifixtures.Client
		client     *Client
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 250*time.Millisecond)

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

		obs = &apifixtures.ApplicationObserver{
			ApplicationAvailableFunc: func(a *Application) {
				defer GinkgoRecover()
				Fail("unexpected application available notification")
			},
			ApplicationUnavailableFunc: func(a *Application) {
				defer GinkgoRecover()
				Fail("unexpected application unavailable notification")
			},
		}

		inspector = &Inspector{
			Observer:        obs,
			BackoffStrategy: backoff.Constant(100 * time.Millisecond),
		}

		apiClient = &apifixtures.Client{
			ListApplicationsFunc: func(context.Context) ([]configkit.Application, error) {
				return []configkit.Application{cfg1, cfg2}, nil
			},
		}

		client = &Client{
			Client: apiClient,
			Target: &Target{Name: "<target>"},
		}
	})

	AfterEach(func() {
		cancel()
	})

	Describe("Run()", func() {
		It("notifies the observer", func() {
			runCtx, cancelRun := context.WithCancel(ctx)
			defer cancelRun()

			var available, unavailable []*Application

			obs.ApplicationAvailableFunc = func(a *Application) {
				available = append(available, a)

				if len(available) == 2 {
					cancelRun()
				}
			}

			obs.ApplicationUnavailableFunc = func(a *Application) {
				unavailable = append(unavailable, a)
			}

			err := inspector.Run(runCtx, client)
			Expect(err).To(Equal(context.Canceled))
			Expect(available).To(HaveLen(2))
			Expect(available).To(ConsistOf(unavailable))
			Expect(configkit.IsApplicationEqual(available[0], cfg1)).To(BeTrue())
			Expect(configkit.IsApplicationEqual(available[1], cfg2)).To(BeTrue())
		})

		It("does not notify the observer if the application list is never obtained", func() {
			apiClient.ListApplicationsFunc = func(context.Context) ([]configkit.Application, error) {
				return nil, errors.New("<error>")
			}

			err := inspector.Run(ctx, client)
			Expect(err).To(Equal(context.DeadlineExceeded))
		})

		It("does not notify the observer if the application is ignored", func() {
			inspector.Ignore = func(a configkit.Application) bool {
				return a.Identity().Key == "<app-key-1>"
			}

			runCtx, cancelRun := context.WithCancel(ctx)
			defer cancelRun()

			var available []*Application

			obs.ApplicationAvailableFunc = func(a *Application) {
				available = append(available, a)
			}

			obs.ApplicationUnavailableFunc = nil

			err := inspector.Run(runCtx, client)
			Expect(err).To(Equal(context.DeadlineExceeded))
			Expect(available).To(HaveLen(1))
			Expect(configkit.IsApplicationEqual(available[0], cfg2)).To(BeTrue())
		})

		It("returns immediately if all applications are ignored", func() {
			inspector.Ignore = func(a configkit.Application) bool {
				return true
			}

			err := inspector.Run(ctx, client)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("retries if the query fails", func() {
			first := true

			apiClient.ListApplicationsFunc = func(context.Context) ([]configkit.Application, error) {
				if first {
					first = false
					return nil, errors.New("<error>")
				}

				return []configkit.Application{cfg1, cfg2}, nil
			}

			runCtx, cancelRun := context.WithCancel(ctx)
			defer cancelRun()

			obs.ApplicationAvailableFunc = func(a *Application) {
				cancelRun()
			}

			obs.ApplicationUnavailableFunc = nil

			err := inspector.Run(runCtx, client)
			Expect(err).To(Equal(context.Canceled))
		})
	})
})
