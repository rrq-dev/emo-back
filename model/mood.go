package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MoodReflection untuk menyimpan mood yang sudah dicatat user
type MoodReflection struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     string             `bson:"user_id" json:"user_id"`
	UserName   string             `bson:"user_name" json:"user_name"`
	Mood       string             `bson:"mood" json:"mood"`
	Reflection string             `bson:"reflection" json:"reflection"`
	Timestamp  time.Time          `bson:"timestamp" json:"timestamp"`
	Processed  bool               `bson:"processed,omitempty" json:"processed,omitempty"` // optional, for Gemini processing status
}

// MoodInput untuk menerima input dari frontend
type MoodInput struct {
	Mood        string `json:"mood" bson:"mood" validate:"required,oneof=happy neutral sad frustrated"`
	Message     string `json:"message" bson:"message"`
	IsAnonymous bool   `json:"is_anonymous" bson:"is_anonymous"` // konsisten dengan backend & frontend
	UserID      string `json:"user_id,omitempty" bson:"user_id,omitempty"` // snake_case biar sama dengan backend fields
}

type SystemPrompt struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Text string             `bson:"text" json:"text"`
}
