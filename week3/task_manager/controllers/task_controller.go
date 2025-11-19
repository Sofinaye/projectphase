package controllers

import (
	"net/http"
	"task_manager/data"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	service *data.TaskService
}

func NewTaskController(service *data.TaskService) *TaskController {
	return &TaskController{service: service}
}
func (tc *TaskController) GetTasks(c *gin.Context) {
	tasks := tc.service.GetAllTasks()
	c.JSON(http.StatusOK, gin.H{
		"data": tasks,
	})
}
