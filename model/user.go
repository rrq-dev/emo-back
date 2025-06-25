package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User untuk koleksi users
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// Payload untuk data JWT (nggak perlu diubah banyak)
type Payload struct {
	Name string `json:"name"`
	Email string `json:"email"`
	User string `json:"user"`
	Role string `json:"role"`
	Iat  int64  `json:"iat"`
	Nbf  int64  `json:"nbf"`
	Exp  int64  `json:"exp"`
	ID   string `json:"id"`
}

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Feedback struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Message   string    `json:"message" bson:"message" validate:"required"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}