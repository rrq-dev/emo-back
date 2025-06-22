package controller

import (
	"emobackend/config"
	"emobackend/helper"
	"emobackend/model"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
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
		userData, ok := userDataRaw.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token tidak valid (tidak bisa dibaca)",
			})
		}

		idVal, okID := userData["id"].(string)
		nameVal, okName := userData["name"].(string)
		if !okID || !okName {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Token tidak memiliki data user yang lengkap",
			})
		}

		userID = idVal
		userName = nameVal

	}

	// Simpan mood refleksi seperti biasa
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

	// Auto-trigger Gemini untuk merespon curhat
	go func() {
		if input.Message != "" {
			err := callGeminiAndSaveReflection(userID, input.Message, input.IsAnonymous)
			if err != nil {
				log.Println("Gagal menyimpan refleksi Gemini:", err)
			}
		}
	}()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil menyimpan refleksi mood",
		"data":    reflection,
	})
}







