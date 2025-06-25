package middleware

import (
	"encoding/json"
	"strings"

	"emobackend/helper"

	"github.com/gofiber/fiber/v2"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Cek apakah request bertipe JSON dan mengandung "is_anonymous": true
		if strings.Contains(c.Get("Content-Type"), "application/json") {
			body := c.Body()
			var temp map[string]interface{}
			if err := json.Unmarshal(body, &temp); err == nil {
				if isAnon, ok := temp["is_anonymous"].(bool); ok && isAnon {
					// Skip validasi token kalau request anonim
					return c.Next()
				}
			}
		}

		//Ambil Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing Authorization header",
			})
		}

		// Ambil token JWT dari header
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Verifikasi token JWT
		payload, err := helper.VerifyJWTToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Simpan user_id ke context (bisa juga simpan name, role, dll)
		c.Locals("user_id", payload.ID)
		c.Locals("user_name", payload.User)
		c.Locals("user_role", payload.Role)

		// Lanjut ke handler berikutnya
		return c.Next()
	}
}

func JWTOptional() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			// Jika tidak ada token, anggap anonymous
			return c.Next()
		}

		// Tetap verifikasi kalau token ada
		claims, err := helper.ParseJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Simpan user_id ke context
		c.Locals("user_id", claims.ID)
		return c.Next()
	}
}

