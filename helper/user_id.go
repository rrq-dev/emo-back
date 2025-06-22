package helper

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateUserID() string {
	bytes := make([]byte, 6) // menghasilkan 12 karakter hex
	if _, err := rand.Read(bytes); err != nil {
		// fallback jika gagal
		return "user-fallback"
	}
	return fmt.Sprintf("user-%s", hex.EncodeToString(bytes))
}