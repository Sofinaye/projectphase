package router

import (
	"task_manager/controllers"
	"task_manager/data"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin router and routes.
func SetupRouter(taskService *data.TaskService) *gin.Engine {
	r := gin.Default()

	taskController := controllers.NewTaskController(taskService)

	// Routes
	r.GET("/tasks", taskController.GetTasks)
	r.GET("/tasks/:id", taskController.GetTaskByID)
	r.POST("/tasks", taskController.CreateTask)
	r.PUT("/tasks/:id", taskController.UpdateTask)
	r.DELETE("/tasks/:id", taskController.DeleteTask)

	return r
}
