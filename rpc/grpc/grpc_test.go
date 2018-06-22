package grpc

import (
	"fmt"
	"log"
	"net"
	"sync"
	"testing"
	"time"

	pb "github.com/shohi/gopattern/rpc/grpc/greeter"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port        = ":50051"
	defaultName = "world"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func TestGRPC(t *testing.T) {
	var address = "localhost" + port
	var wg sync.WaitGroup
	wg.Add(1)

	// Server
	go func() {
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterGreeterServer(s, &server{})
		s.Serve(lis)
	}()

	// client
	go func() {
		defer wg.Done()

		// Set up a connection to the server.
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		var c [20]pb.GreeterClient
		// Contact the server and print out its response.
		name := defaultName

		//warm up
		i := 0
		for ; i < 20; i++ {
			c[i] = pb.NewGreeterClient(conn)
			invoke(c[i], name)
		}
		//
		log.Print("sync ==> ")
		syncTest(c[0], name)

		log.Print("async ==> ")
		asyncTest(c, name)
	}()

	wg.Wait()

}

func invoke(c pb.GreeterClient, name string) {
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	_ = r
}

func syncTest(c pb.GreeterClient, name string) {
	i := 10000
	t := time.Now().UnixNano()
	for ; i > 0; i-- {
		invoke(c, name)
	}
	fmt.Println("took", (time.Now().UnixNano()-t)/1000000, "ms")
}

func asyncTest(c [20]pb.GreeterClient, name string) {
	var wg sync.WaitGroup
	wg.Add(10000)

	i := 10000
	t := time.Now().UnixNano()
	for ; i > 0; i-- {
		go func() { invoke(c[i%20], name); wg.Done() }()
	}
	wg.Wait()
	fmt.Println("took", (time.Now().UnixNano()-t)/1000000, "ms")
}
