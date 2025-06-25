package main

import (
	"emobackend/config"
	"emobackend/routes"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func init(){
	_ = godotenv.Load() // Load environment variables from .env file
}
func main() {
	app := fiber.New()
	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
    AllowOrigins: strings.Join(config.GetAllowedOrigins(), ","),
    AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
}))

	config.MongoConnect()

	routes.SetupURL(app)

	app.Listen(":1506")

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Endpoint not found",
		})
	})
}