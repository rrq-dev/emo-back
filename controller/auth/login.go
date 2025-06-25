package controller

import (
	"context"
	"emobackend/config"
	"emobackend/helper"
	"emobackend/model"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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
	collection := config.DB.Collection("users")
	filter := bson.M{"email": input.Email}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
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

	// Generate JWT token
	token, err := helper.EncodeWithRoleHours("user", user.Name, user.Email, user.ID, 2)
	if err != nil {
		log.Printf("JWT token generation failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
			"details": fmt.Sprintf("JWT error: %v", err),
		})
	}

	// Simpan token JWT ke local storage
	c.Set("token", token)
	c.Set("expires", fmt.Sprintf("%d", time.Now().Add(2*time.Hour).Unix()))

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