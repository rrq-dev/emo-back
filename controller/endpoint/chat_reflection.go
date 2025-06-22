package controller

import (
	"bytes"
	"emobackend/config"
	"emobackend/model"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetAllChatReflections(c *fiber.Ctx) error {
	var reflections []model.ChatReflection

	if err := config.DB.Order("created_at DESC").Find(&reflections).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal mengambil data dari database",
		})
	}

	return c.JSON(reflections)
}

func PostChatReflection(c *fiber.Ctx) error {
	var input model.ChatRequest
	if err := c.BodyParser(&input); err != nil || input.Message == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Message is required"})
	}

	// Prompt empatik untuk Gemini
	prompt := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": "Kamu adalah sahabat reflektif dan suportif. Tanggapi dengan empati dan motivasi berdasarkan curhatan ini: " + input.Message},
				},
			},
		},
	}
	payload, _ := json.Marshal(prompt)

	// Request ke Gemini API
	apiKey := os.Getenv("GEMINI_API_KEY")
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=" + apiKey
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to connect to Gemini"})
	}
	defer resp.Body.Close()

	var geminiRes model.GeminiReply
	json.NewDecoder(resp.Body).Decode(&geminiRes)

	// Ambil respon Gemini
	reply := "Maaf, saya belum bisa membalas curhatan kamu."
	if len(geminiRes.Candidates) > 0 && len(geminiRes.Candidates[0].Content.Parts) > 0 {
		reply = geminiRes.Candidates[0].Content.Parts[0].Text
	}

	// Simpan ke DB via GORM
	reflection := model.ChatReflection{
		Message:     input.Message,
		AIReply:     reply,
		IsAnonymous: true,
		CreatedAt:   time.Now(),
	}

	if err := config.DB.Create(&reflection).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan ke database"})
	}

	return c.JSON(fiber.Map{
		"reply": reply,
	})
}