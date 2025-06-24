package controller

import (
	"emobackend/config"
	"emobackend/helper"
	"emobackend/model"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(c *fiber.Ctx) error {
	var input model.User
	if err := c.BodyParser(&input); err != nil {
		log.Printf("BodyParser error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Validasi input
	if input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and password required"})
	}

	log.Printf("Login attempt for email: %s", input.Email)

	var user model.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("User not found: %s", input.Email)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		log.Printf("Database error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		log.Printf("Wrong password for user: %s", input.Email)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Wrong password"})
	}

	// ✅ Enhanced PASETO token generation dengan validation
	log.Printf("Attempting to generate PASETO token for user: %s", user.Email)
	
	// Cek environment variables
	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		log.Printf("PRIVATE_KEY environment variable not set")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server configuration error", 
			"details": "Private key not configured",
		})
	}

	// Generate PASETO token dengan error handling yang detailed
	token, err := helper.EncodeWithRoleHours("user", user.Name, 2)
	if err != nil {
		log.Printf("PASETO token generation failed: %v", err)
		log.Printf("Private key length: %d", len(privateKey))
		
		// Return detailed error untuk debugging
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
			"details": fmt.Sprintf("PASETO error: %v", err),
			"hint": "Check PRIVATE_KEY format and PASETO library version",
		})
	}

	log.Printf("PASETO token generated successfully for user: %s", user.Email)

	// Environment detection
	isProduction := os.Getenv("ENV") == "production" || os.Getenv("NODE_ENV") == "production"
	
	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "paseto_token",
		Value:    token,
		Expires:  time.Now().Add(2 * time.Hour),
		HTTPOnly: true,
		Secure:   isProduction,
		SameSite: "Lax",
		Path:     "/",
	})

	// Response yang konsisten
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

