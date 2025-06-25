package controller

import (
	"bytes"
	"context"
	"emobackend/config"
	"emobackend/helper"
	"emobackend/model"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/mgo.v2/bson"
)

type ChatMessage struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type ChatSessionRequest struct {
	SessionID   string        `json:"session_id"`
	Messages    []ChatMessage `json:"messages"`
	IsAnonymous bool          `json:"is_anonymous"`
}


func PostChatSession(c *fiber.Ctx) error {
	var req ChatSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if len(req.Messages) == 0 || req.SessionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing session or message"})
	}

	// 🔥 Ambil prompt dari database
	promptText, err := helper.GetSystemPrompt()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil prompt dari database"})
	}

	// Bangun isi pesan dengan prompt di awal
	contents := []map[string]interface{}{
		{
			"role": "system",
			"parts": []map[string]string{
				{"text": promptText},
			},
		},
	}

	// Tambahkan semua pesan user sebelumnya
	for _, m := range req.Messages {
		contents = append(contents, map[string]interface{}{
			"role": m.Role,
			"parts": []map[string]string{
				{"text": m.Text},
			},
		})
	}

	// Payload ke Gemini
	payload := map[string]interface{}{
		"contents": contents,
	}

	jsonPayload, _ := json.Marshal(payload)

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=" + os.Getenv("GEMINI_API_KEY")
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

	// Simpan ke database
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

	// Simpan balasan Gemini
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



func GetChatBySession(c *fiber.Ctx) error {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "session_id dibutuhkan"})
	}

	var chats []model.ChatReflection

	cursor, err := config.DB.Collection("gemini_chat").
		Find(context.TODO(), bson.M{"session_id": sessionID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil chat"})
	}
	if err := cursor.All(context.TODO(), &chats); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal decode chat"})
	}

	return c.JSON(chats)
}

