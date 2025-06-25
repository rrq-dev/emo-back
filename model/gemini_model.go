package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChatReflection — untuk menyimpan chat history antara user dan AI
type ChatReflection struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      *string            `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Message     string             `bson:"message" json:"message"`
	AIReply     string             `bson:"ai_reply" json:"ai_reply"`
	IsAnonymous bool               `bson:"is_anonymous" json:"is_anonymous"`
	Role        string             `bson:"role" json:"role"`           // NEW
	SessionID   string             `bson:"session_id" json:"session_id"` // NEW
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}



// ChatRequest — request body yang dikirim dari user ke backend
type ChatRequest struct {
	Message string `json:"message" bson:"message"`
}

// GeminiReply — format response dari Gemini API
type GeminiReply struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}