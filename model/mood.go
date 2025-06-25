package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MoodReflection untuk koleksi mood_reflections
type MoodReflection struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     string             `bson:"user_id" json:"user_id"`
	UserName   string             `bson:"user_name" json:"user_name"`
	Mood       string             `bson:"mood" json:"mood"`
	Reflection string             `bson:"reflection" json:"reflection"`
	Timestamp  time.Time          `bson:"timestamp" json:"timestamp"`
}

// MoodInput untuk menerima input dari user (misalnya dari frontend)
type MoodInput struct {
	Mood        string `json:"mood" bson:"mood" validate:"required,oneof=happy neutral sad frustrated"`
	Message     string `json:"message" bson:"message"`
	IsAnonymous bool   `json:"isAnonymous" bson:"is_anonymous"`
	UserID      string `json:"userId,omitempty" bson:"user_id,omitempty"` 
}

