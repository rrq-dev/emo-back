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
)

type ChatMessage struct {
	Role string `json:"role"`  // "user" or "model"
	Text string `json:"text"`
}

type ChatSessionRequest struct {
	SessionID string        `json:"session_id"`
	Messages  []ChatMessage `json:"messages"`
	IsAnonymous bool        `json:"is_anonymous"`
}

func PostChatSession(c *fiber.Ctx) error {
	var req ChatSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if len(req.Messages) == 0 || req.SessionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing session or message"})
	}

	// Ambil message terakhir untuk dikirim ke Gemini
	last := req.Messages[len(req.Messages)-1]

	// Format Gemini request
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": last.Text},
				},
			},
		},
	}
	jsonPayload, _ := json.Marshal(payload)

	apiKey := os.Getenv("GEMINI_API_KEY")
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=" + apiKey

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to connect to Gemini"})
	}
	defer resp.Body.Close()

	var geminiRes model.GeminiReply
	if err := json.NewDecoder(resp.Body).Decode(&geminiRes); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Invalid Gemini response"})
	}

	reply := "Saya tidak bisa menjawab saat ini."
	if len(geminiRes.Candidates) > 0 && len(geminiRes.Candidates[0].Content.Parts) > 0 {
		reply = geminiRes.Candidates[0].Content.Parts[0].Text
	}

	// Simpan user messages
	collection := config.DB.Collection("gemini_chat")
	for _, m := range req.Messages {
		collection.InsertOne(context.TODO(), model.ChatReflection{
			SessionID:   req.SessionID,
			Message:     m.Text,
			AIReply:     "",
			IsAnonymous: req.IsAnonymous,
			Role:        m.Role,
			CreatedAt:   time.Now(),
		})
	}

	// Simpan AI balasan
	collection.InsertOne(context.TODO(), model.ChatReflection{
		SessionID:   req.SessionID,
		Message:     "",
		AIReply:     reply,
		IsAnonymous: req.IsAnonymous,
		Role:        "model",
		CreatedAt:   time.Now(),
	})

	return c.JSON(fiber.Map{
		"reply": reply,
	})
}

