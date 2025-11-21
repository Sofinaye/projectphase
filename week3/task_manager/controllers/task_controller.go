package controllers

import (
	"net/http"
	"strconv"
	"task_manager/data"
	"task_manager/models"

	"github.com/gin-gonic/gin"
)

// TaskController handles incoming HTTP requests.
type TaskController struct {
	service *data.TaskService
}

// NewTaskController creates a new TaskController.
func NewTaskController(service *data.TaskService) *TaskController {
	return &TaskController{service: service}
}

// GetTasks handles GET /tasks
func (tc *TaskController) GetTasks(c *gin.Context) {
	tasks, err := tc.service.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "failed to fetch tasks",
			"detail": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

// GetTaskByID handles GET /tasks/:id
func (tc *TaskController) GetTaskByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	task, found, err := tc.service.GetTaskByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "failed to fetch task",
			"detail": err.Error(),
		})
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

// CreateTask handles POST /tasks
func (tc *TaskController) CreateTask(c *gin.Context) {
	var input models.TaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid request payload",
			"detail": err.Error(),
		})
		return
	}

	task, err := tc.service.CreateTask(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "failed to create task",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": task})
}

// UpdateTask handles PUT /tasks/:id
func (tc *TaskController) UpdateTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	var input models.TaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid request payload",
			"detail": err.Error(),
		})
		return
	}

	updated, err := tc.service.UpdateTask(id, input)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  "failed to update task",
				"detail": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updated})
}

// DeleteTask handles DELETE /tasks/:id
func (tc *TaskController) DeleteTask(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	if err := tc.service.DeleteTask(id); err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  "failed to delete task",
				"detail": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}
