package middleware

import (
	"emobackend/config"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

// JWTProtected mengamankan route menggunakan JWT
func JWTProtected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   config.JwtKey, // ✅ Gunakan config.JwtKey, bukan os.Getenv lagi
		ContextKey:   "user",
		ErrorHandler: jwtError,
	})
}

// jwtError akan dijalankan kalau token tidak valid / tidak ada
func jwtError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "User tidak terautentikasi",
		"error":   err.Error(),
	})
}
