package main

import (
	"time"

	"gateway/services/pb"
	"gateway/weblib"
	"gateway/wrappers"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"
)

func main() {
	etcdReg := etcd.NewRegistry(
		registry.Addrs("127.0.0.1:2379"),
	)
	userMicroService := micro.NewService(
		micro.Name("userService.client"),
		micro.WrapClient(wrappers.NewUserWrapper),
	)

	userService := pb.NewUserService("rpcUserService", userMicroService.Client())
	server := web.NewService(
		web.Name("httpService"),
		web.Address(":4000"),
		web.Handler(weblib.NewRouter(userService)),
		web.Registry(etcdReg),
		web.RegisterTTL(time.Second*30),
		web.RegisterInterval(time.Second*15),
		web.Metadata(map[string]string{"protocol": "http"}),
	)

	// 接收命令行参数
	_ = server.Init()
	_ = server.Run()
}
