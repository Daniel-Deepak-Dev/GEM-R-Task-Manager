package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Task represents our to-do item structure
type Task struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Completed   bool               `json:"completed" bson:"completed"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

var taskCollection *mongo.Collection

func main() {
	// 1. Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.Background())

	// Get the tasks collection (MongoDB will create it automatically if it doesn't exist)
	taskCollection = client.Database("taskmanager").Collection("tasks")
	log.Println("Connected to MongoDB!")

	// 2. Create Fiber app
	app := fiber.New()

	// 3. Configure CORS - This allows your React Native app to communicate with the backend
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // In production, specify your frontend URL instead of "*"
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	// 4. Define API routes
	app.Get("/api/tasks", getTasks)          // Get all tasks
	app.Post("/api/tasks", createTask)       // Create a new task
	app.Put("/api/tasks/:id", updateTask)    // Update a task
	app.Delete("/api/tasks/:id", deleteTask) // Delete a task

	// 5. Start the server
	log.Println("Server running on http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}

// GET /api/tasks - Fetch all tasks from MongoDB
func getTasks(c *fiber.Ctx) error {
	var tasks []Task

	// Find all documents in the tasks collection
	cursor, err := taskCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch tasks"})
	}
	defer cursor.Close(context.Background())

	// Decode the results into our tasks slice
	if err := cursor.All(context.Background(), &tasks); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to decode tasks"})
	}

	// If no tasks found, return empty array instead of null
	if tasks == nil {
		tasks = []Task{}
	}

	return c.JSON(tasks)
}

// POST /api/tasks - Create a new task
func createTask(c *fiber.Ctx) error {
	task := new(Task)

	// Parse the JSON body from the request
	if err := c.BodyParser(task); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body :("})
	}

	// Set default values
	task.ID = primitive.NewObjectID()
	task.CreatedAt = time.Now()
	if task.Completed == false {
		task.Completed = false // Explicitly set default
	}

	// Insert the task into MongoDB
	_, err := taskCollection.InsertOne(context.Background(), task)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create task :D"})
	}

	return c.Status(201).JSON(task)
}

// PUT /api/tasks/:id - Update an existing task
func updateTask(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid task ID LOL!"})
	}

	task := new(Task)
	if err := c.BodyParser(task); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body haha!"})
	}

	// Update the task in MongoDB
	update := bson.M{
		"$set": bson.M{
			"title":       task.Title,
			"description": task.Description,
			"completed":   task.Completed,
		},
	}

	result, err := taskCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update task"})
	}

	if result.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found xD"})
	}

	task.ID = objectID
	return c.JSON(task)
}

// DELETE /api/tasks/:id - Delete a task
func deleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid task ID :("})
	}

	result, err := taskCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete task :D"})
	}

	if result.DeletedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Task not found :D"})
	}

	return c.JSON(fiber.Map{"message": "Task deleted successfully :D"})
}
