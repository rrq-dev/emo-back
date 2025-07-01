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

		// Routes yang bisa anonymous atau login
	api.Post("/moods", middleware.JWTOptional(), control.SubmitMoods)
	api.Get("/moods", control.GetAllMoodsData) // Public route
	app.Get("/prompts", control.GetAllSystemPrompts)
	app.Get("/prompts-sg", control.GetPromptSuggestions)


		// Gemini route
	api.Get("/reflection/latest", control.GetAllChatReflections)
	api.Post("/reflection", control.PostChatReflection)


		//chatroom
	api.Post("/chat-session", control.PostChatSession)
	api.Get("/chat-session", control.GetChatBySession)



		// feedback
	app.Post("/feedback", control.SubmitFeedback)
	


}


