package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "helloworld/helloworld"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultName = "world"

var (
	addr = flag.String("addr", "localhost:50051", "the address to connact to")
	name = flag.String("name", defaultName, "Name to greet")
)

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

	// resp, _ := cli.Get(context.Background(), "foo/bar/my-service")
	// conn, err := grpc.Dial(string(resp.Kvs[0].Value), grpc.WithTransportCredentials(insecure.NewCredentials()))

	etcdResolver, err := resolver.NewBuilder(cli)
	conn, err := grpc.Dial(
		"etcd:///foo/bar/my-service",
		grpc.WithResolvers(etcdResolver),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect; %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
