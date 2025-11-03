package main

import (
	"library_management/controllers"
	"library_management/services"
)

func main() {
	lib := services.NewLibrary()
	controller := controllers.NewLibraryController(lib)
	controller.Run()
}
