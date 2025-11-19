package controllers

import (
	"net/http"
	"strconv"
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

func (tc *TaskController) GetTaskByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	task, ok := tc.service.GetTaskByID(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}
