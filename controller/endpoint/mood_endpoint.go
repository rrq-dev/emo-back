package controller

import (
	"emobackend/config"
	"emobackend/helper"
	"emobackend/model"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)


func GetAllMoodReflections(c *fiber.Ctx) error {
	var reflections []model.MoodReflection

	if err := config.DB.Find(&reflections).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data refleksi mood",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data refleksi mood",
		"data":    reflections,
	})
}

func GetReflections(c *fiber.Ctx) error {
	userID := c.Params("userID")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User ID tidak boleh kosong",
		})
	}

	var reflections []model.MoodReflection
	err := config.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&reflections).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data refleksi mood",
			"error":   err.Error(),
		})
	}

	if len(reflections) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Belum ada data refleksi untuk user ini",
			"data":    []model.MoodReflection{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data refleksi mood",
		"data":    reflections,
	})
}


func SubmitMoodReflections(c *fiber.Ctx) error {
	var input model.MoodInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Input tidak valid",
			"error":   err.Error(),
		})
	}

	var userID, userName string

	if input.IsAnonymous {
		userID = "anon-" + helper.GenerateUserID()
		userName = "Anonymous"
	} else {
		userDataRaw := c.Locals("user")
		if userDataRaw == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "User tidak terautentikasi",
			})
		}

		payload, ok := userDataRaw.(model.Payload)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Gagal membaca token Paseto",
			})
		}

		userID = "user-" + payload.User
		userName = payload.User
	}

	reflection := model.MoodReflection{
		UserID:     userID,
		UserName:   userName,
		Mood:       input.Mood,
		Reflection: input.Message,
		Timestamp:  time.Now(),
	}

	if err := config.DB.Create(&reflection).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menyimpan refleksi mood",
			"error":   err.Error(),
		})
	}

	// Trigger Gemini async
	go func() {
		if input.Message != "" {
			err := callGeminiAndSaveReflection(userID, input.Message, input.IsAnonymous)
			if err != nil {
				fmt.Println("Gagal menyimpan refleksi Gemini:", err)
			}
		}
	}()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil menyimpan refleksi mood",
		"data":    reflection,
	})
}









