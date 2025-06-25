package controller

import (
	"context"
	"emobackend/config"
	"emobackend/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func SubmitFeedback(c *fiber.Ctx) error {
	var input struct {
		Message string `json:"message" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Gagal parsing input",
		})
	}

	// Validasi manual (optional kalau pakai validator lib)
	if input.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Pesan feedback tidak boleh kosong",
		})
	}

	feedback := model.Feedback{
		ID:        uuid.New().String(),
		Message:   input.Message,
		CreatedAt: time.Now(),
	}

	collection := config.DB.Collection("feedback")
	_, err := collection.InsertOne(context.TODO(), feedback)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan feedback",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Feedback berhasil dikirim!",
	})
}

