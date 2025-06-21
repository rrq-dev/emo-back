package model

import (
	"time"
)

type MoodReflection struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     string    `gorm:"type:varchar(255);not null" json:"user_id"`
	UserName   string    `gorm:"type:varchar(255)" json:"user_name"` // ✅ tambahkan ini
	Mood       string    `gorm:"type:varchar(50);not null" json:"mood"`
	Reflection string    `gorm:"type:text" json:"reflection"`
	Timestamp  time.Time `gorm:"autoCreateTime" json:"timestamp"`
}

type MoodInput struct {
	Mood        string `json:"mood" validate:"required,oneof=happy neutral sad frustrated"`
	Message     string `json:"message"`         // opsional
	IsAnonymous bool   `json:"is_anonymous"`    // true jika user tidak login
	UserID      string `json:"user_id"`         // opsional, akan di-override dari token jika login
}
