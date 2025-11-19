package controllers

import (
	"task_manager/data"
)

type TaskController struct {
	service *data.TaskService
}

func NewTaskController(service *data.TaskService) *TaskController {
	return &TaskController{service: service}
}
