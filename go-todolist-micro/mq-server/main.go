package main

import (
	"mq/conf"
	"mq/service"
)

func main() {
	conf.Init()

	forever := make(chan bool)
	service.CreateTask()
	<-forever
}
