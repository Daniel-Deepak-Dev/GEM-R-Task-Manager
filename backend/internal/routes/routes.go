package routes

import (
	"backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, taskHandler *handlers.TaskHandler) {
	api := app.Group("/api")

	api.Get("/tasks", taskHandler.GetTasks)
	api.Post("/tasks", taskHandler.CreateTask)
	api.Put("/tasks/:id", taskHandler.UpdateTask)
	api.Delete("/tasks/:id", taskHandler.DeleteTask)
	api.Get("/tasks/:id", taskHandler.GetTask)
}
