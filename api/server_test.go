package api

import (
	"context"
	"net"
	"time"

	"github.com/dogmatiq/configkit/api/internal/pb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("func RegisterServer()", func() {
	It("panics if one of the applications can not be marshaled", func() {
		Expect(func() {
			s := grpc.NewServer()
			RegisterServer(s, &application{})
		}).To(Panic())
	})
})

var _ = Describe("type server", func() {
	var (
		ctx      context.Context
		cancel   func()
		listener net.Listener
		gserver  *grpc.Server
		client   pb.ConfigClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)

		var err error
		listener, err = net.Listen("tcp", ":")
		Expect(err).ShouldNot(HaveOccurred())

		gserver = grpc.NewServer()
		RegisterServer(gserver)

		go gserver.Serve(listener)

		conn, err := grpc.Dial(
			listener.Addr().String(),
			grpc.WithInsecure(),
		)
		Expect(err).ShouldNot(HaveOccurred())

		client = pb.NewConfigClient(conn)
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

	Describe("func Watch()", func() {
		It("blocks until the server is stopped", func() {
			go func() {
				time.Sleep(500 * time.Millisecond)
				gserver.Stop()
			}()

			stream, err := client.Watch(ctx, &pb.WatchRequest{})
			Expect(err).ShouldNot(HaveOccurred())

			_, err = stream.Recv()
			Expect(err).ShouldNot(HaveOccurred())

			_, err = stream.Recv()
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.Unavailable))
		})
	})
})
