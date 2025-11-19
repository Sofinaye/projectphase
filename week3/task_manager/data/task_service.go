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
