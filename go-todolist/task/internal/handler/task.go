package handler

import (
	"context"
	"fmt"

	"task/internal/repository"
	"task/internal/service/pb"
	"task/pkg/e"
)

type TaskService struct{}

func NewTaskService() *TaskService {
	return &TaskService{}
}

func (*TaskService) TaskCreate(ctx context.Context, req *pb.TaskRequest) (resp *pb.CommonResponse, err error) {
	fmt.Println("in taskService")
	var task repository.Task
	resp = new(pb.CommonResponse)
	if err = task.TaskCreate(req); err != nil {
		resp.Code = e.Error
		resp.Msg = e.GetMsg(e.Error)
		resp.Data = err.Error()
		return resp, err
	}

	resp.Msg = e.GetMsg(e.Success)
	resp.Code = e.Success
	return resp, nil
}

func (*TaskService) TaskShow(ctx context.Context, req *pb.TaskRequest) (resp *pb.TasksDetailResponse, err error) {
	var task repository.Task
	resp = new(pb.TasksDetailResponse)
	resp.Code = e.Success
	tasks, err := task.TaskShow(req)
	if err != nil {
		resp.Code = e.Error
		return resp, err
	}
	resp.TaskDetail = repository.BuildTasks(tasks)
	return resp, nil
}

func (*TaskService) TaskUpdate(ctx context.Context, req *pb.TaskRequest) (resp *pb.CommonResponse, err error) {
	var task repository.Task
	resp = new(pb.CommonResponse)
	if err = task.TaskUpdate(req); err != nil {
		resp.Code = e.Error
		resp.Msg = e.GetMsg(e.Error)
		resp.Data = err.Error()
		return resp, err
	}

	resp.Msg = e.GetMsg(uint(resp.Code))
	return resp, nil
}

func (*TaskService) TaskDelete(ctx context.Context, req *pb.TaskRequest) (resp *pb.CommonResponse, err error) {
	var task repository.Task
	resp = new(pb.CommonResponse)
	if err = task.TaskDelete(req); err != nil {
		resp.Code = e.Error
		resp.Msg = e.GetMsg(e.Error)
		resp.Data = err.Error()
		return resp, err
	}

	resp.Msg = e.GetMsg(uint(resp.Code))
	return resp, nil
}
