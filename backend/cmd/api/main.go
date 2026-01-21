package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	// Import your internal packages
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/routes"
)

func main() {
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Connect to Database
	client := database.ConnectDB(cfg.MongoURI)
	// Disconnect when main exits
	defer client.Disconnect(context.Background())

	// 3. Initialize Handler
	// We inject the specific collection into the handler
	coll := client.Database("taskmanager").Collection("tasks")
	taskHandler := handlers.NewTaskHandler(coll)

	// 4. Setup Fiber App
	app := fiber.New()
	app.Use(cors.New())

	// 5. Setup Routes
	routes.SetupRoutes(app, taskHandler)

	// 6. Start Server
	log.Println("Server running on port", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
