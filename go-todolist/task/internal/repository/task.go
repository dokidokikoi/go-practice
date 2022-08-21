package repository

import "task/internal/service/pb"

type Task struct {
	TaskId    uint `gorm:"primarykey"`
	UserId    uint `gorm:"index"`
	Status    int  `gorm:"default:0"`
	Title     string
	Content   string `gorm:"type:longtext"`
	StartTime int64
	EndTime   int64
}

func (*Task) TaskCreate(req *pb.TaskRequest) error {
	task := Task{
		UserId:    uint(req.UserID),
		Title:     req.Title,
		Content:   req.Content,
		StartTime: int64(req.StartTime),
		EndTime:   int64(req.EndTime),
	}

	return DB.Create(&task).Error
}

func (*Task) TaskShow(req *pb.TaskRequest) (taskList []Task, err error) {
	err = DB.Model(&Task{}).Where("user_id=?", req.UserID).Find(&taskList).Error
	if err != nil {
		return nil, err
	}

	return taskList, nil
}

func BuildTask(task *Task) (taskModel *pb.TaskModel) {
	return &pb.TaskModel{
		TaskID:    uint32(task.TaskId),
		UserID:    uint32(task.UserId),
		Status:    uint32(task.Status),
		Title:     task.Title,
		Content:   task.Content,
		StartTime: uint32(task.StartTime),
		EndTime:   uint32(task.EndTime),
	}
}

func BuildTasks(taskList []Task) (tasks []*pb.TaskModel) {
	for _, v := range taskList {
		t := BuildTask(&v)
		tasks = append(tasks, t)
	}
	return tasks
}

func (*Task) TaskUpdate(req *pb.TaskRequest) error {
	t := Task{}
	if err := DB.Where("task_id=?", req.TaskID).First(&t).Error; err != nil {
		return err
	}

	t.Status = int(req.Status)
	t.Title = req.Title
	t.Content = req.Content
	t.StartTime = int64(req.StartTime)
	t.EndTime = int64(req.EndTime)
	return DB.Save(&t).Error
}

func (*Task) TaskDelete(req *pb.TaskRequest) error {
	return DB.Model(&Task{}).Where("task_id=?", req.TaskID).Delete(&Task{}).Error
}
