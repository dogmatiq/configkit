package discovery_test

import (
	. "github.com/dogmatiq/configkit/api/discovery"
	"github.com/dogmatiq/configkit/api/discovery/fixtures" // can't dot-import due to conflict
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	_ TargetPublisher = (*TargetObserverSet)(nil)
	_ TargetObserver  = (*TargetObserverSet)(nil)
)

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
