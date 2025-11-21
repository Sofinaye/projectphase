package main

import (
	"log"
	"task_manager/data"
	"task_manager/router"
)

func main() {
	taskService := data.NewTaskService()
	r := router.SetupRouter(taskService)

	log.Println("Task Management API is running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
