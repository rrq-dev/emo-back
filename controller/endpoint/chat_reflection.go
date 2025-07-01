package controller

import (
	"bytes"
	"context"
	"emobackend/config"
	"emobackend/model"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"emobackend/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
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
	var input struct {
		Messages    []string `json:"messages"`
		SessionID   string   `json:"session_id"`
		IsAnonymous bool     `json:"is_anonymous"`
	}
	if err := c.BodyParser(&input); err != nil || len(input.Messages) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Messages is required"})
	}

	// Ambil user id dari JWT jika ada
	userID := ""
	token := c.Get("Authorization")
	if token != "" {
		if payload, err := helper.VerifyJWTToken(token); err == nil {
			userID = payload.ID
		}
	}

	// Buat session id jika kosong
	sessionID := input.SessionID
	if sessionID == "" {
		sessionID = uuid.New().String()
	}

	collection := config.DB.Collection("gemini_chat")
	var lastReply string
	for _, msg := range input.Messages {
		// Prompt ke Gemini
		prompt := map[string]interface{}{
			"contents": []map[string]interface{}{
				{
					"role": "user",
					"parts": []map[string]string{
						{"text": msg},
					},
				},
			},
		}
		payload, _ := json.Marshal(prompt)
		apiKey := os.Getenv("GEMINI_API_KEY")
		url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=" + apiKey
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to connect to Gemini"})
		}
		defer resp.Body.Close()
		var geminiRes model.GeminiReply
		json.NewDecoder(resp.Body).Decode(&geminiRes)
		reply := "Maaf, saya belum bisa membalas curhatan kamu."
		if len(geminiRes.Candidates) > 0 && len(geminiRes.Candidates[0].Content.Parts) > 0 {
			reply = geminiRes.Candidates[0].Content.Parts[0].Text
		}
		lastReply = reply
		// Simpan chat user
		reflectionUser := model.ChatReflection{
			UserID:      &userID,
			SessionID:   sessionID,
			Message:     msg,
			AIReply:     "",
			IsAnonymous: input.IsAnonymous,
			Role:        "user",
			CreatedAt:   time.Now(),
		}
		fmt.Println("Insert Chat User:", reflectionUser)
		collection.InsertOne(context.Background(), reflectionUser)
		// Simpan balasan Gemini
		reflectionAI := model.ChatReflection{
			UserID:      &userID,
			SessionID:   sessionID,
			Message:     "",
			AIReply:     reply,
			IsAnonymous: input.IsAnonymous,
			Role:        "model",
			CreatedAt:   time.Now(),
		}
		fmt.Println("Insert Chat AI:", reflectionAI)
		collection.InsertOne(context.Background(), reflectionAI)
	}
	return c.JSON(fiber.Map{
		"reply": lastReply,
		"session_id": sessionID,
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

func ProcessLatestMoodToGemini() error {
	// Ambil mood terbaru dari submit_mood
	var latest model.MoodReflection

	opts := options.FindOne().SetSort(map[string]interface{}{"timestamp": -1})
	err := config.DB.Collection("submit_mood").FindOne(context.TODO(), bson.M{}, opts).Decode(&latest)
	if err != nil {
		return fmt.Errorf("gagal ambil data submit_mood: %v", err)
	}

	// Panggil Gemini untuk proses refleksi
	prompt := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": "Mood saya: " + latest.Mood + ". Curhatan saya: " + latest.Reflection + ". Tolong berikan refleksi atau dukungan yang suportif."},
				},
			},
		},
	}

	payload, _ := json.Marshal(prompt)
	apiKey := os.Getenv("GEMINI_API_KEY")
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=" + apiKey

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("gagal request ke Gemini: %v", err)
	}
	defer resp.Body.Close()

	var geminiRes model.GeminiReply
	if err := json.NewDecoder(resp.Body).Decode(&geminiRes); err != nil {
		return fmt.Errorf("gagal decode response Gemini: %v", err)
	}

	reply := "Refleksi tidak tersedia saat ini."
	if len(geminiRes.Candidates) > 0 && len(geminiRes.Candidates[0].Content.Parts) > 0 {
		reply = geminiRes.Candidates[0].Content.Parts[0].Text
	}

	// Simpan ke gemini_chat
	reflection := model.ChatReflection{
		UserID:      &latest.UserID,
		Message:     latest.Reflection,
		AIReply:     reply,
		IsAnonymous: latest.UserName == "Anonymous",
		CreatedAt:   time.Now(),
	}

	_, err = config.DB.Collection("gemini_chat").InsertOne(context.TODO(), reflection)
	if err != nil {
		return fmt.Errorf("gagal simpan ke gemini_chat: %v", err)
	}

	return nil
}

