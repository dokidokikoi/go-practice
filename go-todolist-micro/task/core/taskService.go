package core

import (
	"context"
	"encoding/json"
	"errors"

	"task/model"
	"task/service/pb"

	"github.com/streadway/amqp"
)

func (*TaskService) CreateTask(ctx context.Context, req *pb.TaskRequest, resp *pb.TaskDetailResponse) error {
	ch, err := model.MQ.Channel()
	if err != nil {
		err = errors.New("rabbitMQ channel err" + err.Error())
		return err
	}

	q, _ := ch.QueueDeclare("task_queue", true, false, false, false, nil)
	body, _ := json.Marshal(req)
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})

	if err != nil {
		err = errors.New("rabbitMQ public err" + err.Error())
	}
	return err
}

func (*TaskService) GetTasksList(ctx context.Context, req *pb.TaskRequest, resp *pb.TaskListResponse) error {
	if req.Limit == 0 {
		req.Limit = 10
	}
	var taskData []model.Task
	var count uint32
	// 查找备忘录
	err := model.DB.Offset(req.Start).Limit(req.Limit).Where("uid=?", req.Uid).Find(&taskData).Error
	if err != nil {
		return errors.New("mysql find : " + err.Error())
	}
	// 统计数量
	model.DB.Model(&model.Task{}).Where("uid=?", req.Uid).Count(&count)
	// 返回proto里面定义的类型
	var taskRes []*pb.TaskModel
	for _, item := range taskData {
		taskRes = append(taskRes, BuildTask(item))
	}
	resp.TaskList = taskRes
	resp.Count = count
	return nil
}

// 获取详细的备忘录
func (*TaskService) GetTask(ctx context.Context, req *pb.TaskRequest, resp *pb.TaskDetailResponse) error {
	taskData := model.Task{}
	model.DB.First(&taskData, req.Id)
	taskRes := BuildTask(taskData)
	resp.TaskDetail = taskRes
	return nil
}

func (*TaskService) UpdateTask(ctx context.Context, req *pb.TaskRequest, resp *pb.TaskDetailResponse) error {
	taskData := model.Task{}
	// 查找该用户的这条信息
	model.DB.Model(&model.Task{}).Where("id= ? AND uid=?", req.Id, req.Uid).First(&taskData)
	taskData.Title = req.Title
	taskData.Status = int(req.Status)
	taskData.Content = req.Content
	model.DB.Save(&taskData)
	resp.TaskDetail = BuildTask(taskData)
	return nil
}

func (*TaskService) DeleteTask(ctx context.Context, req *pb.TaskRequest, resp *pb.TaskDetailResponse) error {
	err := model.DB.Model(&model.Task{}).Where("id =? AND uid=?", req.Id, req.Uid).Delete(&model.Task{}).Error
	if err != nil {
		return errors.New("删除失败：" + err.Error())
	}
	return nil
}
