package middleware

import (
	"emobackend/config"
	"log"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func JWTProtected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   config.JwtKey,
		ContextKey:   "user",
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	log.Println("[JWT ERROR]:", err.Error())
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "User tidak terautentikasi",
		"error":   err.Error(),
	})
}
