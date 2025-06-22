package controller

import (
	"emobackend/config"
	"emobackend/helper"
	"emobackend/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(c *fiber.Ctx) error {
	var input model.User
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user model.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Wrong password"})
	}

	// Generate PASETO token
	token, err := helper.EncodeWithRoleHours("user", user.Name, 2)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "paseto_token",
		Value:    token,
		Expires:  time.Now().Add(2 * time.Hour),
		HTTPOnly: true,                // aman dari JS access
		Secure:   true,                // aktifkan kalau pakai HTTPS
		SameSite: "Lax",               // atau "Strict" kalau mau lebih ketat
		Path:     "/",
	})

	// Kirim juga di body biar frontend bisa simpan juga kalau perlu
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   token, // opsional kalau hanya mau pakai cookie
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

