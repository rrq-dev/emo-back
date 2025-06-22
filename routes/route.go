package routes

import (
	ctrl "emobackend/controller/auth"
	control "emobackend/controller/endpoint"
	"emobackend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupURL(app *fiber.App) {
	api := app.Group("/api")

	// Auth Links
	api.Post("/login", ctrl.Login)
	api.Post("/register", ctrl.Register)

	// Endpoint
	api.Get("/moods", control.GetAllMoodReflections)
	api.Post("/moods", middleware.PasetoMiddleware(), control.SubmitMoodReflections)
	api.Get("/moods/:userID", control.GetReflections)
	app.Get("/chat-reflections", control.GetAllChatReflections)
	app.Post("/chat-refleksi", control.PostChatReflection)

}
