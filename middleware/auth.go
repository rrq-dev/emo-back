package middleware

import (
	"emobackend/helper"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func PasetoMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token tidak ditemukan",
			})
		}

		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Format token tidak valid",
			})
		}

		token := parts[1]

		payload, err := helper.Decoder(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token tidak valid",
				"error":   err.Error(),
			})
		}

		// Simpan ke context
		c.Locals("user", payload)
		return c.Next()
	}
}
