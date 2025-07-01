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
	"github.com/google/uuid"
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
	PromptID    string        `json:"prompt_id"` 
}

func PostChatSession(c *fiber.Ctx) error {
	var req ChatSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if len(req.Messages) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Missing message"})
	}
	// Buat session id jika kosong
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = uuid.New().String()
	}
	// Ambil user id dari JWT jika ada
	userID := ""
	token := c.Get("Authorization")
	if token != "" {
		if payload, err := helper.VerifyJWTToken(token); err == nil {
			userID = payload.ID
		}
	}
	promptText := "Kamu adalah AI pendamping refleksi emosi. Balas dengan empati."
	if req.PromptID != "" {
		p, err := helper.GetPromptByID(req.PromptID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil prompt"})
		}
		promptText = p
	}
	contents := []map[string]interface{}{
		{
			"role": "system",
			"parts": []map[string]string{
				{"text": promptText},
			},
		},
	}
	for _, m := range req.Messages {
		contents = append(contents, map[string]interface{}{
			"role": m.Role,
			"parts": []map[string]string{
				{"text": m.Text},
			},
		})
	}
	payload := map[string]interface{}{"contents": contents}
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
	collection := config.DB.Collection("gemini_chat")
	for _, m := range req.Messages {
		collection.InsertOne(context.TODO(), model.ChatReflection{
			UserID:      &userID,
			SessionID:   sessionID,
			Message:     m.Text,
			AIReply:     "",
			IsAnonymous: req.IsAnonymous,
			Role:        m.Role,
			CreatedAt:   time.Now(),
		})
	}
	collection.InsertOne(context.TODO(), model.ChatReflection{
		UserID:      &userID,
		SessionID:   sessionID,
		Message:     "",
		AIReply:     reply,
		IsAnonymous: req.IsAnonymous,
		Role:        "model",
		CreatedAt:   time.Now(),
	})
	return c.JSON(fiber.Map{"reply": reply, "session_id": sessionID})
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
	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &chats); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal decode chat"})
	}

	// Definisikan struct respons yang sesuai dengan format JSON yang diinginkan
	type MongoObjectID struct {
		Oid string `json:"$oid"`
	}

	type MongoDate struct {
		Date string `json:"$date"`
	}

	type ChatReflectionResponse struct {
		ID          MongoObjectID `json:"_id"`
		UserID      *string       `json:"user_id,omitempty"`
		Message     string        `json:"message"`
		AIReply     string        `json:"ai_reply"`	
		IsAnonymous bool          `json:"is_anonymous"`
		Role        string        `json:"role"`
		SessionID   string        `json:"session_id"`
		CreatedAt   MongoDate     `json:"created_at"`
	}

	var responseChats []ChatReflectionResponse
	for _, chat := range chats {
		var userID *string
		if chat.UserID != nil {
			id := *chat.UserID // Dereference pointer untuk menyalin nilai
			userID = &id       // Re-reference salinan nilai agar bisa di-assign ke pointer di response struct
		}

		responseChats = append(responseChats, ChatReflectionResponse{
			ID: MongoObjectID{
				Oid: chat.ID.Hex(), // Konversi ObjectID ke string heksadesimalnya
			},
			UserID:      userID,
			Message:     chat.Message,
			AIReply:     chat.AIReply,
			IsAnonymous: chat.IsAnonymous,
			Role:        chat.Role,
			SessionID:   chat.SessionID,
			CreatedAt: MongoDate{
				Date: chat.CreatedAt.Format(time.RFC3339Nano), // Format waktu ke ISO 8601 dengan milidetik
			},
		})
	}

	return c.JSON(responseChats)
}

