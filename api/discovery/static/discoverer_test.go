package static_test

import (
	"context"
	"time"

	"github.com/dogmatiq/configkit/api/discovery"
	. "github.com/dogmatiq/configkit/api/discovery/static"
	. "github.com/dogmatiq/configkit/api/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Discoverer", func() {
	var (
		ctx              context.Context
		cancel           context.CancelFunc
		obs              *TargetObserver
		target1, target2 *discovery.Target
		disc             *Discoverer
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)

		target1 = &discovery.Target{Name: "<target-1>"}
		target2 = &discovery.Target{Name: "<target-2>"}
		obs = &TargetObserver{}
		disc = &Discoverer{
			Targets:  []*discovery.Target{target1, target2},
			Observer: obs,
		}
	})

	AfterEach(func() {
		cancel()
	})

	Describe("func Run()", func() {
		It("notifies the observer immediately", func() {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			var targets []*discovery.Target
			obs.TargetAvailableFunc = func(t *discovery.Target) {
				targets = append(targets, t)

				if len(targets) == len(disc.Targets) {
					obs.TargetUnavailableFunc = nil
					cancel()
				}
			}

			obs.TargetUnavailableFunc = func(t *discovery.Target) {
				Fail("observer unexpectedly notified of target unavailability")
			}

			err := disc.Run(ctx)
			Expect(err).To(Equal(context.Canceled))
			Expect(targets).To(ConsistOf(target1, target2))
		})

		It("notifies the observer when the discoverer is stopped", func() {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			count := 0
			obs.TargetAvailableFunc = func(t *discovery.Target) {
				count++
				if count == 2 {
					cancel()
				}
			}

			var targets []*discovery.Target
			obs.TargetUnavailableFunc = func(t *discovery.Target) {
				targets = append(targets, t)
			}

			err := disc.Run(ctx)
			Expect(err).To(Equal(context.Canceled))
			Expect(targets).To(ConsistOf(target1, target2))
		})
	})
})
