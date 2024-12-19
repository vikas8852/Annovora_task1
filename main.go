package main

import (
	"TASK/db"
	"TASK/handlers"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db.Connect("mongodb://admin:admin%40academichub@34.55.9.182:27017/")

	if db.Client == nil {
		log.Fatal("MongoDB client is not initialized. Exiting...")
	}

	handlers.TaskCollection = db.Client.Database("vikas").Collection("tasks")
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/tasks", handlers.CreateTask)
	e.GET("/tasks", handlers.GetAllTasks)
	e.GET("/tasks/:id", handlers.GetTask)
	e.PUT("/tasks/:id", handlers.UpdateTask)
	e.DELETE("/tasks/:id", handlers.DeleteTask)

	e.Logger.Fatal(e.Start(":5000"))
}