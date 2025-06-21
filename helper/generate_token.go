package helper

import (
	"crypto/rand"
	"emobackend/model"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(userID uint, email, name string) (string, error) {
	claims := model.Claims{
		ID:    userID,
		Email: email,
		Name:  name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // 24-hour expiration
			Issuer:    "emobackend",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func GenerateUserID() string {
	bytes := make([]byte, 6) // menghasilkan 12 karakter hex
	if _, err := rand.Read(bytes); err != nil {
		// fallback jika gagal
		return "user-fallback"
	}
	return fmt.Sprintf("user-%s", hex.EncodeToString(bytes))
}