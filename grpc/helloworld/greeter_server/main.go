package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "helloworld/helloworld"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
)

var port = flag.Int("port", 50051, "The server port")

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelleReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelleReply{Message: "Hello " + in.GetName()}, nil
}

var cli *clientv3.Client

func init() {
	var err error
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 10,
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())

	resp, err := cli.Grant(context.TODO(), 5)
	if err != nil {
		panic(err)
	}
	// if _, err = cli.Put(
	// 	context.TODO(), "foo/bar/my-service",
	// 	fmt.Sprintf("localhost:%d", *port),
	// 	clientv3.WithLease(resp.ID)); err != nil {
	// 	panic(err)
	// }

	em, _ := endpoints.NewManager(cli, "foo/bar/my-service")
	err = em.AddEndpoint(
		context.TODO(),
		"foo/bar/my-service/v1",
		endpoints.Endpoint{Addr: fmt.Sprintf("localhost:%d", *port)},
		clientv3.WithLease(resp.ID))
	if err != nil {
		panic(err)
	}

	if _, err = cli.KeepAlive(context.TODO(), resp.ID); err != nil {
		panic(err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
