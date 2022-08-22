package main

import (
	"user/conf"
	"user/core"
	"user/service/pb"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
)

func main() {
	conf.Init()
	// etcd 注册件
	etcdReg := etcd.NewRegistry(
		registry.Addrs("127.0.0.1:2379"),
	)

	// 得到一个微服务实例
	microService := micro.NewService(
		micro.Name("rpcUserService"),
		micro.Address("127.0.0.1:8082"),
		micro.Registry(etcdReg),
	)

	// 结构命令行参数，初始化
	microService.Init()
	// 服务注册
	_ = pb.RegisterUserServiceHandler(microService.Server(), new(core.UserService))
	// 启动微服务
	_ = microService.Run()
}
