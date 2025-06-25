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

		// ✅ Routes yang bisa anonymous atau login
	api.Post("/moods", middleware.JWTOptional(), control.SubmitMoods)
	api.Get("/moods", control.GetAllMoodsData) // Public route

	 // ✅ Routes yang HARUS login
	// api.Get("/moods/:userId", middleware.JWTProtected(), control.GetUserMoods)
	// api.Get("/chat-reflections", middleware.JWTProtected(), control.GetLatestReflection)
	// api.Post("/chat-refleksi", middleware.JWTProtected(), control.SubmitChatReflection)
}


