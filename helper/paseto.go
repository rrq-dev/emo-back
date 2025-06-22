package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"emobackend/model"

	"aidanwoods.dev/go-paseto"
)

func EncodeWithRoleHours(role, username string, hours int64) (string, error) {
	privatekey := os.Getenv("PRIVATE_KEY")
	token := paseto.NewToken()
	// Set metadata: waktu pembuatan, masa berlaku, dll
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(time.Duration(hours) * time.Hour))
	token.SetString("user", username)
	token.SetString("role", role)
	key, err := paseto.NewV4AsymmetricSecretKeyFromHex(privatekey)
	return token.V4Sign(key, nil), err
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