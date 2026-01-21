package routes

import (
	"backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, taskHandler *handlers.TaskHandler, skillHandler *handlers.SkillHandler, progressHandler *handlers.ProgressHandler) {
	api := app.Group("/api")

	api.Get("/tasks", taskHandler.GetTasks)
	api.Post("/tasks", taskHandler.CreateTask)
	api.Get("/tasks/:id", taskHandler.GetTask)
	api.Put("/tasks/:id", taskHandler.UpdateTask)
	api.Delete("/tasks/:id", taskHandler.DeleteTask)

	// Skill Routes
	api.Get("/skills", skillHandler.GetSkills)
	api.Post("/skills", skillHandler.CreateSkill)
	api.Get("/skills/:id", skillHandler.GetSkill)
	api.Put("/skills/:id", skillHandler.UpdateSkill)
	api.Delete("/skills/:id", skillHandler.DeleteSkill)

	// Progress Routes
	api.Post("/progress", progressHandler.CreateProgressItem)
	api.Get("/progress/skill/:skillId", progressHandler.GetItemsBySkill)
	api.Put("/progress/:id", progressHandler.UpdateProgressItem)
	api.Delete("/progress/:id", progressHandler.DeleteProgressItem)
}
