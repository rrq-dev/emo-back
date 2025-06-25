package middleware

import (
	"encoding/json"
	"strings"

	"emobackend/helper"

	"github.com/gofiber/fiber/v2"
)

// JWTProtected - Middleware untuk route yang HARUS login
func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ✅ PERBAIKAN 1: Cek request body untuk anonymous user
		if strings.Contains(c.Get("Content-Type"), "application/json") {
			body := c.Body()
			if len(body) > 0 {
				var temp map[string]interface{}
				if err := json.Unmarshal(body, &temp); err == nil {
					// Cek berbagai kemungkinan field name
					if isAnon, ok := temp["is_anonymous"].(bool); ok && isAnon {
						// Skip validasi token untuk request anonymous
						c.Locals("is_anonymous", true)
						return c.Next()
					}
					if isAnon, ok := temp["isAnonymous"].(bool); ok && isAnon {
						// Untuk camelCase dari frontend
						c.Locals("is_anonymous", true)
						return c.Next()
					}
				}
			}
		}

		// Ambil Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing Authorization header",
			})
		}

		// Ambil token JWT dari header
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			// Bearer prefix tidak ditemukan
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Authorization format. Use Bearer <token>",
			})
		}

		// Verifikasi token JWT
		payload, err := helper.VerifyJWTToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// ✅ PERBAIKAN 2: Simpan data user ke context dengan field yang konsisten
		c.Locals("user_id", payload.ID)
		c.Locals("user_name", payload.Name)  // Ubah dari payload.User ke payload.Name
		c.Locals("user_email", payload.Email) // Tambah email jika ada
		c.Locals("is_anonymous", false)

		return c.Next()
	}
}

// JWTOptional - Middleware untuk route yang bisa anonymous atau login
func JWTOptional() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ✅ PERBAIKAN 3: Cek request body untuk menentukan anonymous
		var isRequestAnonymous bool
		if strings.Contains(c.Get("Content-Type"), "application/json") {
			body := c.Body()
			if len(body) > 0 {
				var temp map[string]interface{}
				if err := json.Unmarshal(body, &temp); err == nil {
					if isAnon, ok := temp["is_anonymous"].(bool); ok && isAnon {
						isRequestAnonymous = true
					}
					if isAnon, ok := temp["isAnonymous"].(bool); ok && isAnon {
						isRequestAnonymous = true
					}
				}
			}
		}

		// Jika request explicitly anonymous, skip token validation
		if isRequestAnonymous {
			c.Locals("is_anonymous", true)
			return c.Next()
		}

		// Cek Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// Tidak ada token, anggap anonymous
			c.Locals("is_anonymous", true)
			return c.Next()
		}

		// Ada token, coba verifikasi
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			// Format Bearer salah, tapi karena optional, anggap anonymous
			c.Locals("is_anonymous", true)
			return c.Next()
		}

		// Verifikasi token
		payload, err := helper.VerifyJWTToken(tokenStr)
		if err != nil {
			// Token invalid, anggap anonymous (tidak error karena optional)
			c.Locals("is_anonymous", true)
			return c.Next()
		}

		// Token valid, simpan data user
		c.Locals("user_id", payload.ID)
		c.Locals("user_name", payload.Name)
		c.Locals("user_email", payload.Email)
		c.Locals("is_anonymous", false)

		return c.Next()
	}
}

// ✅ BONUS: Helper function untuk mengambil user info dari context
func GetUserFromContext(c *fiber.Ctx) (userID interface{}, userName string, isAnonymous bool) {
	userID = c.Locals("user_id")
	userName, _ = c.Locals("user_name").(string)
	isAnonymous, _ = c.Locals("is_anonymous").(bool)
	return
}