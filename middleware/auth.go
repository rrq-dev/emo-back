package middleware

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

// JWTProtected mengamankan route menggunakan JWT
func JWTProtected() fiber.Handler {
	secret := os.Getenv("JWT_SECRET") // Pastikan JWT_SECRET di .env kamu

	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(secret),
		ContextKey:   "user", // data token nanti disimpan di c.Locals("user")
		ErrorHandler: jwtError,
	})
}

// jwtError akan dijalankan kalau token tidak valid / tidak ada
func jwtError(c *fiber.Ctx, err error) error {
	log.Println("[JWT ERROR]:", err.Error()) // Tambahkan ini untuk debugging
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "User tidak terautentikasi",
		"error":   err.Error(),
	})
}
