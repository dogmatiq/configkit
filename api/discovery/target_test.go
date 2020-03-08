package discovery_test

import (
	"context"
	"time"

	. "github.com/dogmatiq/configkit/api/discovery"
	"github.com/dogmatiq/configkit/api/discovery/fixtures" // can't dot-import due to conflict
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ TargetObserver = (*TargetObserverSet)(nil)

var _ = Describe("type TargetObserverSet", func() {
	var (
		set              *TargetObserverSet
		obs1, obs2       *fixtures.TargetObserver
		target1, target2 *Target
	)

	BeforeEach(func() {
		set = &TargetObserverSet{}

		obs1 = &fixtures.TargetObserver{}
		obs2 = &fixtures.TargetObserver{}

		target1 = &Target{}
		target2 = &Target{}
	})

	Describe("func NewClientObserverSet()", func() {
		It("returns a set containing the given observers", func() {
			set.RegisterTargetObserver(obs1)
			set.RegisterTargetObserver(obs2)

			Expect(
				NewTargetObserverSet(obs1, obs2),
			).To(Equal(set))
		})
	})

	Describe("func TargetAvailable()", func() {
		It("notifies the observers about the target availability", func() {
			var observers []TargetObserver

			obs1.TargetAvailableFunc = func(t *Target) {
				defer GinkgoRecover()
				Expect(t).To(BeIdenticalTo(target1))

				observers = append(observers, obs1)
			}

			obs2.TargetAvailableFunc = func(t *Target) {
				defer GinkgoRecover()
				Expect(t).To(BeIdenticalTo(target1))

				observers = append(observers, obs2)
			}

			set.RegisterTargetObserver(obs1)
			set.RegisterTargetObserver(obs2)
			set.TargetAvailable(target1)

			Expect(observers).To(ConsistOf(obs1, obs2))
		})

		It("does nothing if the target is already available", func() {
			set.RegisterTargetObserver(obs1)
			set.TargetAvailable(target1)

			obs1.TargetAvailableFunc = func(t *Target) {
				defer GinkgoRecover()
				Fail("observer unexpectedly notified of the same target")
			}

			set.TargetAvailable(target1)
		})
	})

	Describe("func TargetUnavailable()", func() {
		It("notifies the observers about the target unavailability", func() {
			var observers []TargetObserver

			obs1.TargetUnavailableFunc = func(t *Target) {
				defer GinkgoRecover()
				Expect(t).To(BeIdenticalTo(target1))

				observers = append(observers, obs1)
			}

			obs2.TargetUnavailableFunc = func(t *Target) {
				defer GinkgoRecover()
				Expect(t).To(BeIdenticalTo(target1))

				observers = append(observers, obs2)
			}

			set.RegisterTargetObserver(obs1)
			set.RegisterTargetObserver(obs2)
			set.TargetAvailable(target1)
			set.TargetUnavailable(target1)

			Expect(observers).To(ConsistOf(obs1, obs2))
		})

		It("does nothing if the target is not in the registry", func() {
			set.RegisterTargetObserver(obs1)

			obs1.TargetUnavailableFunc = func(t *Target) {
				defer GinkgoRecover()
				Fail("observer unexpectedly notified of unknown target")
			}

			set.TargetUnavailable(target1)
		})
	})

	Describe("func RegisterTargetObserver()", func() {
		It("notifies the observer about existing targets", func() {
			set.TargetAvailable(target1)
			set.TargetAvailable(target2)

			var targets []*Target

			obs1.TargetAvailableFunc = func(t *Target) {
				targets = append(targets, t)
			}

			set.RegisterTargetObserver(obs1)

			Expect(targets).To(ConsistOf(target1, target2))
		})

		It("does nothing if the observer is already registered", func() {
			set.TargetAvailable(target1)
			set.RegisterTargetObserver(obs1)

			obs1.TargetAvailableFunc = func(t *Target) {
				defer GinkgoRecover()
				Fail("observer unexpectedly notified when registered twice")
			}

			set.RegisterTargetObserver(obs1)
		})
	})

	Describe("func UnregisterTargetObserver()", func() {
		It("synthesizes an unavailable notification for existing targets", func() {
			set.TargetAvailable(target1)
			set.TargetAvailable(target2)
			set.RegisterTargetObserver(obs1)

			var targets []*Target

			obs1.TargetUnavailableFunc = func(t *Target) {
				targets = append(targets, t)
			}

			set.UnregisterTargetObserver(obs1)

			Expect(targets).To(ConsistOf(target1, target2))
		})

		It("prevents the observer from receiving further notifications", func() {
			obs1.TargetAvailableFunc = func(t *Target) {
				defer GinkgoRecover()
				Fail("unregistered observer unexpectedly notified")
			}

			set.RegisterTargetObserver(obs1)
			set.UnregisterTargetObserver(obs1)
			set.TargetAvailable(target1)
		})

		It("does nothing if the observer is not already registered", func() {
			set.TargetAvailable(target1)

			obs1.TargetUnavailableFunc = func(t *Target) {
				defer GinkgoRecover()
				Fail("observer unexpectedly notified when not registered")
			}

			set.UnregisterTargetObserver(obs1)
		})
	})
})

var _ TargetObserver = (*TargetExecutor)(nil)

var _ = Describe("type TargetExecutor", func() {
	var (
		ctx    context.Context
		cancel func()
		exec   *TargetExecutor
		target *Target
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 250*time.Millisecond)

		exec = &TargetExecutor{
			Task: func(context.Context, *Target) {},
		}

		target = &Target{}
	})

	AfterEach(func() {
		cancel()
	})

	Describe("func TargetAvailable()", func() {
		It("starts a goroutine for the given target", func() {
			barrier := make(chan struct{})

			exec.Task = func(_ context.Context, t *Target) {
				defer GinkgoRecover()
				defer close(barrier)

				Expect(t).To(Equal(target))
			}

			exec.TargetAvailable(target)
			defer exec.TargetUnavailable(target)

			select {
			case <-barrier:
			case <-ctx.Done():
				Expect(ctx.Err()).ShouldNot(HaveOccurred())
			}
		})

		It("does not panic if the target is already available", func() {
			exec.TargetAvailable(target)
			defer exec.TargetUnavailable(target)

			exec.TargetAvailable(target)
		})
	})

	Describe("func TargetUnavailable()", func() {
		It("cancels the context associated with the goroutine and waits for the function to finish", func() {
			barrier := make(chan struct{})

			exec.Task = func(funcCtx context.Context, t *Target) {
				defer GinkgoRecover()
				defer close(barrier)

				select {
				case <-funcCtx.Done():
					// ok
				case <-ctx.Done():
					Expect(ctx.Err()).ShouldNot(HaveOccurred())
				}
			}

			exec.TargetAvailable(target)
			exec.TargetUnavailable(target)

			select {
			case <-barrier:
			case <-ctx.Done():
				Expect(ctx.Err()).ShouldNot(HaveOccurred())
			}
		})

		It("does not panic if the target is already unavailable", func() {
			exec.TargetUnavailable(target)
		})
	})
})
