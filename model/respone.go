package model

import "github.com/golang-jwt/jwt"

type APIResponse struct {
	Status  string      `json:"status"`         // "success", "error", dll
	Message string      `json:"message"`        // Penjelasan singkat
	Data    interface{} `json:"data,omitempty"` // Data (bisa mood list, user info, dll)
}

type Claims struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
}