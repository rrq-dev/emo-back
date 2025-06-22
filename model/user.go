package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Email     string         `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password  string         `gorm:"type:text;not null" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

}

type Payload struct {
	User string `json:"user"`  // Username
	Role string `json:"role"`  // Role (e.g. "user", "admin")
	Iat  int64  `json:"iat"`   // Issued At
	Nbf  int64  `json:"nbf"`   // Not Before
	Exp  int64  `json:"exp"`   // Expiration Time
}