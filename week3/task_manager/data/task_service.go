package data

import (
	"errors"
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
func (s *TaskService) GetAllTasks() []models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]models.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		result = append(result, t)
	}
	return result
}

func (s *TaskService) GetTaskByID(id int) (models.Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	return task, ok
}

func (s *TaskService) CreateTask(input models.TaskInput) models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := models.Task{
		ID:          s.nextID,
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate,
		Status:      input.Status,
	}
	s.tasks[task.ID] = task
	s.nextID++
	return task
}

func (s *TaskService) UpdateTask(id int, input models.TaskInput) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return models.Task{}, errors.New("task not found")
	}

	task.Title = input.Title
	task.Description = input.Description
	task.DueDate = input.DueDate
	task.Status = input.Status

	s.tasks[id] = task
	return task, nil
}
