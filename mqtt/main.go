package main

import (
	"fmt"
	"os"
	"os/signal"

	"practice/mqtt/service"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// 定义MQTT客户端连接参数
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("mqtt-client")
	opts.SetPassword("password")
	opts.SetUsername("test")

	// 创建MQTT客户端实例
	client := mqtt.NewClient(opts)

	// 客户端连接到MQTT服务器
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// 订阅主题
	topic := "mqttopic/process"
	if token := client.Subscribe(topic, 1, service.ProcessMessageHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// service.PrintData()

	// 等待信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	// 取消订阅主题和断开MQTT连接
	client.Unsubscribe(topic)
	client.Disconnect(250)
}
