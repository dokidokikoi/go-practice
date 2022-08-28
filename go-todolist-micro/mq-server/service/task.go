package service

import (
	"encoding/json"
	"log"

	"mq/model"
)

func CreateTask() {
	ch, err := model.MQ.Channel()
	if err != nil {
		panic(err)
	}
	q, _ := ch.QueueDeclare("task_queue", true, false, false, false, nil)
	ch.Qos(1, 0, false)
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	// 处于一个监听状态，一致监听我们的生产端的生产，所以这里我们要阻塞主进程
	go func() {
		for d := range msgs {
			var t model.Task
			err := json.Unmarshal(d.Body, &t)
			if err != nil {
				panic(err)
			}
			model.DB.Create(&t)
			log.Println("Done")
			_ = d.Ack(false)
		}
	}()
}
