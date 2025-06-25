package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Response standar dari API
type APIResponse struct {
	Status  string      `json:"status"`           // "success", "error", dll
	Message string      `json:"message"`          // Penjelasan singkat
	Data    interface{} `json:"data,omitempty"`   // Bisa mood, user, dll
}

// JWT Claims — ganti ID jadi ObjectID agar match dengan MongoDB
type Claims struct {
	ID    primitive.ObjectID `json:"id" bson:"_id"`
	Name  string             `json:"name" bson:"name"`
	Email string             `json:"email" bson:"email"`
	IssuedAt int64           `json:"iat" bson:"-"`
	ExpiresAt int64          `json:"exp" bson:"-"`
}


