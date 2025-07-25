package helper

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"emobackend/model"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var JwtKey = []byte(os.Getenv("JWT_SECRET"))

// Encode JWT dengan data user dan role, berlaku sekian jam
func EncodeWithRoleHours(role, name, email string, userID primitive.ObjectID, hours int) (string, error) {
	claims := jwt.MapClaims{
		"id":    userID.Hex(),
		"name":  name,
		"email": email,
		"role":  role,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Duration(hours) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %v", err)
	}
	return tokenString, nil
}

// Validasi token dan ambil payload (tanpa PASETO)
func VerifyJWTToken(tokenString string) (model.Payload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Pastikan pakai HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		return model.Payload{}, fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return model.Payload{}, errors.New("failed to parse claims")
	}

	// Konversi ke Payload struct
	payload := model.Payload{
		User: claims["name"].(string),
		Role: claims["role"].(string),
		Iat:  int64(claims["iat"].(float64)),
		Nbf:  int64(claims["iat"].(float64)), // default ke iat
		Exp:  int64(claims["exp"].(float64)),
		ID:   claims["id"].(string), // tambahkan jika payload butuh ID
	}

	// Validasi waktu secara manual (opsional)
	now := time.Now().Unix()
	if payload.Nbf > now {
		return model.Payload{}, errors.New("token belum aktif (not before)")
	}
	if payload.Exp < now {
		return model.Payload{}, errors.New("token sudah expired")
	}

	return payload, nil
}



func ParseJWT(tokenString string) (*model.Claims, error) {
	if tokenString == "" {
		return nil, errors.New("token kosong")
	}

	// Hilangkan prefix "Bearer " jika ada
	if after, ok :=strings.CutPrefix(tokenString, "Bearer "); ok  {
		tokenString = after
	}

	// Parsing token dengan klaim kustom
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validasi metode signing
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode signing tidak didukung")
		}
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("token tidak valid atau gagal diparsing")
	}

	claimsMap, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("klaim token tidak valid")
	}

	// Konversi "id" ke ObjectID
	idHex, ok := claimsMap["id"].(string)
	if !ok {
		return nil, errors.New("ID tidak ditemukan dalam klaim")
	}

	objectID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, errors.New("gagal konversi ID ke ObjectID")
	}

	claims := &model.Claims{
		ID:        objectID,
		Name:      getStringFromClaims(claimsMap, "name"),
		Email:     getStringFromClaims(claimsMap, "email"),
		IssuedAt:  int64(getFloatFromClaims(claimsMap, "iat")),
		ExpiresAt: int64(getFloatFromClaims(claimsMap, "exp")),
	}

	// Validasi waktu (manual tambahan, meskipun jwt udah handle)
	now := time.Now().Unix()
	if claims.ExpiresAt < now {
		return nil, errors.New("token sudah expired")
	}

	return claims, nil
}

func getStringFromClaims(m jwt.MapClaims, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func getFloatFromClaims(m jwt.MapClaims, key string) float64 {
	if val, ok := m[key].(float64); ok {
		return val
	}
	return 0
}

