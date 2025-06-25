package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"emobackend/model"

	"aidanwoods.dev/go-paseto"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func EncodeWithRoleHours(role, name, email string, userID primitive.ObjectID, hours int) (string, error) {
	claims := &model.Claims{
		ID:    userID,
		Name:  name,
		Email: email,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Duration(hours) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    userID.Hex(),
		"name":  name,
		"email": email,
		"role":  role,
		"iat":   claims.IssuedAt,
		"exp":   claims.ExpiresAt,
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %v", err)
	}
	return tokenString, nil
}

func VerifyJWTToken(tokenString string) (model.Payload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validasi algoritma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return model.Payload{}, fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return model.Payload{}, errors.New("failed to parse claims")
	}

	// Ekstrak payload dari JWT
	payload := model.Payload{
		User: claims["name"].(string),
		Role: claims["role"].(string),
		Iat:  int64(claims["iat"].(float64)),
		Nbf:  int64(claims["iat"].(float64)), // asumsi aktif sejak iat
		Exp:  int64(claims["exp"].(float64)),
	}

	// Optional: validasi waktu manual (jwt lib juga udah handle sih)
	now := time.Now().Unix()
	if payload.Nbf > now {
		return model.Payload{}, errors.New("token belum aktif (not before)")
	}
	if payload.Exp < now {
		return model.Payload{}, errors.New("token sudah expired")
	}

	return payload, nil
}

func Decoder(tokenStr string) (model.Payload, error) {
	publicKey := os.Getenv("PUBLIC_KEY")

	pubKey, err := paseto.NewV4AsymmetricPublicKeyFromHex(publicKey)
	if err != nil {
		return model.Payload{}, fmt.Errorf("failed to parse public key: %w", err)
	}

	parser := paseto.NewParser()
	token, err := parser.ParseV4Public(pubKey, tokenStr, nil)
	if err != nil {
		return model.Payload{}, fmt.Errorf("failed to parse paseto token: %w", err)
	}

	var payload model.Payload
	if err := json.Unmarshal(token.ClaimsJSON(), &payload); err != nil {
		return model.Payload{}, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// ✅ Validasi waktu: belum aktif atau expired
	now := time.Now().Unix()

	if payload.Nbf > now {
		return model.Payload{}, errors.New("token belum aktif (not before)")
	}
	if payload.Exp < now {
		return model.Payload{}, errors.New("token sudah expired")
	}

	return payload, nil
}