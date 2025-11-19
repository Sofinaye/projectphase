package data

import (
	"sync"
	"task_manager/models"
)

type TaskService struct {
	mu     sync.Mutex
	tasks  map[int]models.Task
	nextID int
}

func NewTaskService() *TaskService {
	return &TaskService{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}
