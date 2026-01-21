package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	// Replace 'backend' with your actual module name from go.mod
	"backend/internal/models"
)

// TaskHandler struct holds the database collection
type TaskHandler struct {
	collection *mongo.Collection
}

// NewTaskHandler creates a new instance of TaskHandler
func NewTaskHandler(col *mongo.Collection) *TaskHandler {
	return &TaskHandler{collection: col}
}

func (h *TaskHandler) GetTasks(c *fiber.Ctx) error {
	var tasks []models.Task
	cursor, err := h.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &tasks); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Return empty list if nil
	if tasks == nil {
		tasks = []models.Task{}
	}
	return c.JSON(tasks)
}

func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	task := new(models.Task)
	if err := c.BodyParser(task); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	task.ID = primitive.NewObjectID()
	task.CreatedAt = time.Now()

	_, err := h.collection.InsertOne(context.Background(), task)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create task"})
	}

	return c.Status(201).JSON(task)
}

// UpdateTask handles updating an existing task
func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	// Use a pointer to the Task model so we handle optional fields correctly
	updateData := new(models.Task)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	// Create the update query
	update := bson.M{
		"$set": bson.M{
			"title":       updateData.Title,
			"description": updateData.Description,
			"completed":   updateData.Completed,
		},
	}

	// Use h.collection here
	result, err := h.collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update task"})
	}

	if result.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}

	// Return the updated data (or fetch the fresh document if strict accuracy is needed)
	updateData.ID = objectID
	return c.JSON(updateData)
}

// DeleteTask handles removing a task
func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	// Use h.collection here
	result, err := h.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete task"})
	}

	if result.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
	}

	return c.JSON(fiber.Map{"message": "Task deleted successfully"})
}

func (h *TaskHandler) GetTask(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var task models.Task

	// FindOne returns a SingleResult. We call Decode(&task) to get the error.
	err = h.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&task)

	if err != nil {
		// Check specifically if the error is "document not found"
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"error": "Task not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get task"})
	}

	return c.JSON(task)
}
