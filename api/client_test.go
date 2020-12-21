package api_test

import (
	"context"
	"net"
	"time"

	. "github.com/dogmatiq/configkit/api"
	"github.com/dogmatiq/interopspec/configspec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
)

type invalidServer struct{}

func (s *invalidServer) ListApplications(
	ctx context.Context,
	req *configspec.ListApplicationsRequest,
) (*configspec.ListApplicationsResponse, error) {
	return &configspec.ListApplicationsResponse{
		Applications: []*configspec.Application{
			{}, // invalid
		},
	}, nil
}

var _ = Describe("type Client", func() {
	var (
		ctx      context.Context
		cancel   func()
		listener net.Listener
		gserver  *grpc.Server
		client   *Client
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)

		var err error
		listener, err = net.Listen("tcp", ":")
		Expect(err).ShouldNot(HaveOccurred())

		gserver = grpc.NewServer()
		configspec.RegisterConfigAPIServer(gserver, &invalidServer{})

		go gserver.Serve(listener)

		conn, err := grpc.Dial(
			listener.Addr().String(),
			grpc.WithInsecure(),
		)
		Expect(err).ShouldNot(HaveOccurred())

		client = NewClient(conn)
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

	Describe("func ListApplications()", func() {
		It("returns an error if the server returns an invalid application", func() {
			_, err := client.ListApplications(ctx)
			Expect(err).Should(HaveOccurred())
		})
	})
})
