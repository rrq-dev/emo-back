package helper

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"emobackend/model"

	"aidanwoods.dev/go-paseto"
)

func EncodeWithRoleHours(role, name string, hours int) (string, error) {
	// Get private key from environment
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex == "" {
		return "", fmt.Errorf("PRIVATE_KEY environment variable not set")
	}

	log.Printf("Private key hex length: %d", len(privateKeyHex))

	// Decode hex private key
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		log.Printf("Failed to decode private key hex: %v", err)
		return "", fmt.Errorf("invalid private key format: %v", err)
	}

	log.Printf("Private key bytes length: %d", len(privateKeyBytes))

	// Create PASETO private key
	privateKey, err := paseto.NewV4AsymmetricSecretKeyFromBytes(privateKeyBytes)
	if err != nil {
		log.Printf("Failed to create PASETO private key: %v", err)
		return "", fmt.Errorf("failed to create PASETO key: %v", err)
	}

	// Create token with claims
	token := paseto.NewToken()
	
	// Set standard claims
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(time.Duration(hours) * time.Hour))
	
	// Set custom claims
	token.SetString("role", role)
	token.SetString("name", name)
	
	// Log claims for debugging
	log.Printf("Creating token with claims: role=%s, name=%s, exp=%v", 
		role, name, time.Now().Add(time.Duration(hours) * time.Hour))

	// Sign token
	signedToken := token.V4Sign(privateKey, nil)
	
	if signedToken == "" {
		return "", fmt.Errorf("failed to sign token - empty result")
	}

	log.Printf("Token signed successfully, length: %d", len(signedToken))
	
	return signedToken, nil
}

// ✅ Function untuk verify token (berguna untuk middleware)
func VerifyPasetoToken(tokenString string) (*paseto.Token, error) {
	publicKeyHex := os.Getenv("PUBLIC_KEY")
	if publicKeyHex == "" {
		return nil, fmt.Errorf("PUBLIC_KEY environment variable not set")
	}

	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid public key format: %v", err)
	}

	publicKey, err := paseto.NewV4AsymmetricPublicKeyFromBytes(publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create PASETO public key: %v", err)
	}

	parser := paseto.NewParser()
	token, err := parser.ParseV4Public(publicKey, tokenString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	return token, nil
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