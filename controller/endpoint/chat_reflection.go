package controller

import (
	"bytes"
	"context"
	"emobackend/config"
	"emobackend/model"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllChatReflections(c *fiber.Ctx) error {
	var reflections []model.ChatReflection

	collection := config.DB.Collection("gemini_chat")
	cursor, err := collection.Find(context.Background(), map[string]interface{}{}, 
		// Sort by created_at descending
		&options.FindOptions{Sort: map[string]interface{}{"created_at": -1}},
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal mengambil data dari database",
		})
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &reflections); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal memproses data dari database",
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

	// Simpan ke DB via MongoDB
reflection := model.ChatReflection{
    Message:     input.Message,
    AIReply:     reply,
    IsAnonymous: true,
    CreatedAt:   time.Now(),
}

collection := config.DB.Collection("gemini_chat") // Replace with your actual collection name
if _, err := collection.InsertOne(context.Background(), reflection); err != nil {
    return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan ke database"})
}

return c.JSON(fiber.Map{
    "reply": reply,
})
}

func callGeminiAndSaveReflection(userID string, message string, isAnonymous bool) error {
	prompt := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": "Mood saya sedang begini: " + message + ". Tolong berikan motivasi atau refleksi yang suportif."},
				},
			},
		},
	}

	payload, _ := json.Marshal(prompt)
	apiKey := os.Getenv("GEMINI_API_KEY")
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=" + apiKey

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var geminiRes struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&geminiRes); err != nil {
		return err
	}

	reply := "Terima kasih sudah berbagi. Kamu tidak sendirian 💛"
	if len(geminiRes.Candidates) > 0 && len(geminiRes.Candidates[0].Content.Parts) > 0 {
		reply = geminiRes.Candidates[0].Content.Parts[0].Text
	}

	// Simpan ke chat_reflections
ref := model.ChatReflection{
    UserID:      &userID,
    Message:     message,
    AIReply:     reply,
    IsAnonymous: isAnonymous,
    CreatedAt:   time.Now(),
}

collection := config.DB.Collection("gemini_chat")
_, insertErr := collection.InsertOne(context.Background(), ref)
if insertErr != nil {
	return insertErr
}

return nil
}
