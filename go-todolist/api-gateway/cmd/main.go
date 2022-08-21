package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gateway/config"
	"gateway/discovery"
	"gateway/internal/service/pb"
	"gateway/routes"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

func main() {
	config.InitConfig()

	// 服务发现
	etcdAddress := []string{viper.GetString("etcd.addres")}
	etcdRegister := discovery.NewReslover(etcdAddress, logrus.New())
	resolver.Register(etcdRegister)
	go startListen()
	{
		osSignal := make(chan os.Signal, 1)
		signal.Notify(osSignal, os.Interrupt, os.Kill, syscall.SIGALRM, syscall.SIGINT, syscall.SIGKILL)
		s := <-osSignal
		fmt.Println("exit!", s)
	}
	fmt.Println("gateway listen on :3000")
}

func startListen() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	userConn, err := grpc.Dial("127.0.0.1:10001", opts...)
	if err != nil {
		panic(err)
	}
	userService := pb.NewUserServiceClient(userConn)

	taskConn, err := grpc.Dial("127.0.0.1:10002", opts...)
	if err != nil {
		panic(err)
	}
	taskService := pb.NewTaskServiceClient(taskConn)

	ginRouter := routes.NewRouter(userService, taskService)
	server := &http.Server{
		Addr:           ":" + viper.GetString("server.port"),
		Handler:        ginRouter,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("err", err)
	}

}
