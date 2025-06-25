package controller

import (
	"context"
	"emobackend/config"
	"log"

	"emobackend/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var input model.RegisterInput
	if err := c.BodyParser(&input); err != nil {
		log.Printf("BodyParser error: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Validasi input
	if input.Name == "" || input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name, email, and password required"})
	}

	log.Printf("Register attempt for email: %s", input.Email)

	// Cek apakah email sudah digunakan
	var user model.User
	collection := config.DB.Collection("users")
	filter := bson.M{"email": input.Email}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err == nil {
		log.Printf("Email sudah digunakan: %s", input.Email)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already used"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		log.Printf("Hash password error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Hash password error"})
	}

	// Simpan user ke database
	user = model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}
	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Printf("Database error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// Response yang konsisten
	return c.JSON(fiber.Map{
		"message": "Register successful",
		"user": fiber.Map{
			"name":  user.Name,
			"email": user.Email,
		},
	})
}