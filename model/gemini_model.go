package model

import "time"

type ChatReflection struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID      *string   `json:"user_id,omitempty"`
	Message     string    `json:"message"`
	AIReply     string    `json:"ai_reply"`
	IsAnonymous bool      `json:"is_anonymous"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type ChatRequest struct {
	Message string `json:"message"`
}

type GeminiReply struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}