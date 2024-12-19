package handlers

import (
	"context"
	"net/http"
	"time"

	"TASK/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Reference to the MongoDB collection
var TaskCollection *mongo.Collection

// CreateTask handles creating a new task
func CreateTask(c echo.Context) error {

	task := new(models.Task)

	// Bind request body to the Task model
	if err := c.Bind(task); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Set task properties
	task.ID = primitive.NewObjectID()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	// Insert the task into the database
	_, err := TaskCollection.InsertOne(context.Background(), task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// Return the created task
	return c.JSON(http.StatusCreated, task)
}
func GetAllTasks(c echo.Context) error {
	pointer, err := TaskCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	var tasks []models.Task
	if err := pointer.All(context.Background(), &tasks); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, tasks)
}

//get task

func GetTask(c echo.Context) error {
	id := c.Param("id")
	objid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	var task models.Task
	err = TaskCollection.FindOne(context.Background(), bson.M{"_id": objid}).Decode(&task)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, task)
}

//Updating task

func UpdateTask(c echo.Context) error {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid ID"})
	}
	task := new(models.Task)
	if err := c.Bind(task); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	task.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"updated_at":  task.UpdatedAt,
		},
	}
	_, err = TaskCollection.UpdateOne(context.Background(), bson.M{"_id": objectId}, update)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Task Updated Successfully"})

}

//Deleting task

func DeleteTask(c echo.Context) error {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"error": "Invalid ID"})
	}
	_, err = TaskCollection.DeleteOne(context.Background(), bson.M{"_id": objectId})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Task deleted successfully"})
}