package middleware

import (
	"emobackend/helper"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Cek apakah body mengandung "is_anonymous": true
		body := c.Body()
		if strings.Contains(string(body), `"is_anonymous":true`) {
			return c.Next() // Lewatin middleware kalau anonim
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing Authorization header",
			})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return helper.JwtKey, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		// Simpan user_id ke context
		c.Locals("user_id", claims["user_id"])
		return c.Next()
	}
}
