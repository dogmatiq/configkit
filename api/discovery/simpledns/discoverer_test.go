package simpledns_test

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/dogmatiq/configkit/api/discovery"
	. "github.com/dogmatiq/configkit/api/discovery/simpledns"
	. "github.com/dogmatiq/configkit/api/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Discoverer", func() {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		obs    *TargetObserver
		res    *resolver
		disc   *Discoverer
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)

		obs = &TargetObserver{}
		res = &resolver{}
		disc = &Discoverer{
			QueryHost: "<query-host>",
			Observer:  obs,
			Resolver:  res,
		}
	})

	AfterEach(func() {
		cancel()
	})

	Describe("func Run()", func() {
		It("notifies the observer when targets become available", func() {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			res.LookupHostFunc = func(_ context.Context, host string) ([]string, error) {
				// canceling here relies on the fact that the discoverer does
				// not check the context while notifying the observer.
				cancel()

				Expect(host).To(Equal("<query-host>"))
				return []string{"<host1>", "<host2>"}, nil
			}

			var targets []*discovery.Target
			obs.TargetAvailableFunc = func(t *discovery.Target) {
				targets = append(targets, t)
			}

			err := disc.Run(ctx)
			Expect(err).To(Equal(context.Canceled))
			Expect(targets).To(ConsistOf(
				&discovery.Target{
					Name: "<host1>:https",
					MetaData: discovery.MetaData{
						QueryHostKey: "<query-host>",
					},
				},
				&discovery.Target{
					Name: "<host2>:https",
					MetaData: discovery.MetaData{
						QueryHostKey: "<query-host>",
					},
				},
			))
		})

		It("notifies the observer when targets become unavailable", func() {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			disc.Interval = 10 * time.Millisecond
			first := true

			res.LookupHostFunc = func(_ context.Context, host string) ([]string, error) {
				if first {
					first = false
					return []string{"<host1>", "<host2>"}, nil
				}

				return []string{"<host2>"}, nil
			}

			obs.TargetUnavailableFunc = func(t *discovery.Target) {
				defer GinkgoRecover()

				Expect(t).To(Equal(
					&discovery.Target{
						Name: "<host1>:https",
						MetaData: discovery.MetaData{
							QueryHostKey: "<query-host>",
						},
					},
				))

				// Prevent a failure when <host2> becomes unavailable simply
				// because the discover is stopped.
				obs.TargetUnavailableFunc = nil

				// canceling here relies on the fact that the discoverer does
				// not check the context while notifying the observer.
				cancel()
			}

			err := disc.Run(ctx)
			Expect(err).To(Equal(context.Canceled))
		})

		It("notifies the observer when the discoverer is stopped", func() {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			disc.Interval = 10 * time.Millisecond

			res.LookupHostFunc = func(_ context.Context, host string) ([]string, error) {
				// canceling here relies on the fact that the discoverer
				// does not check the context while notifying the observer.
				cancel()

				return []string{"<host1>", "<host2>"}, nil
			}

			var targets []*discovery.Target
			obs.TargetUnavailableFunc = func(t *discovery.Target) {
				targets = append(targets, t)
			}

			err := disc.Run(ctx)
			Expect(err).To(Equal(context.Canceled))
			Expect(targets).To(ConsistOf(
				&discovery.Target{
					Name: "<host1>:https",
					MetaData: discovery.MetaData{
						QueryHostKey: "<query-host>",
					},
				},
				&discovery.Target{
					Name: "<host2>:https",
					MetaData: discovery.MetaData{
						QueryHostKey: "<query-host>",
					},
				},
			))
		})

		It("uses net.DefaultResolver by default", func() {
			disc.QueryHost = "localhost"
			disc.Resolver = nil

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			obs.TargetAvailableFunc = func(t *discovery.Target) {
				// canceling here relies on the fact that the discoverer does
				// not check the context while notifying the observer.
				cancel()

				Expect(t.Name).To(
					Or(
						Equal("127.0.0.1:https"),
						Equal("[::1]:https"),
					),
				)
			}

			err := disc.Run(ctx)
			Expect(err).To(Equal(context.Canceled))
		})

		When("the resolver fails", func() {
			It("does not propagate not-found errors", func() {
				ctx, cancel := context.WithCancel(ctx)
				defer cancel()

				res.LookupHostFunc = func(context.Context, string) ([]string, error) {
					cancel()
					return nil, &net.DNSError{
						IsNotFound: true,
					}
				}

				err := disc.Run(ctx)
				Expect(err).To(Equal(context.Canceled)) // note: not the net.DNSError
			})

			It("propagates other errors", func() {
				res.LookupHostFunc = func(context.Context, string) ([]string, error) {
					return nil, errors.New("<error>")
				}

				err := disc.Run(ctx)
				Expect(err).To(MatchError("<error>"))
			})
		})
	})
})

type resolver struct {
	LookupHostFunc func(ctx context.Context, host string) ([]string, error)
}

func (r *resolver) LookupHost(ctx context.Context, host string) ([]string, error) {
	if r.LookupHostFunc == nil {
		return nil, &net.DNSError{
			IsNotFound: true,
		}
	}

	return r.LookupHostFunc(ctx, host)
}
