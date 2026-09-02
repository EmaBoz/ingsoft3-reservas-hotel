package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
)

const saltSize = 16 // bytes

// HashPassword genera salt aleatorio y retorna saltHex$hashHex
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("empty password")
	}
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	h := sha256.New()
	h.Write(salt)
	h.Write([]byte(password))
	sum := h.Sum(nil)

	saltHex := hex.EncodeToString(salt)
	hashHex := hex.EncodeToString(sum)

	return saltHex + "$" + hashHex, nil
}

// CheckPassword compara el hashedPassword (salt$hash) con el password en tiempo constante
func CheckPassword(hashedPassword, password string) error {
	if hashedPassword == "" || password == "" {
		return errors.New("invalid password or hash")
	}
	parts := make([]string, 2)
	n, _ := fmt.Sscanf(hashedPassword, "%[^$]$%s", &parts[0], &parts[1])
	if n != 2 {
		// fallback: intentar split simple

		for i := 0; i < len(hashedPassword); i++ {
			// noop - we will do a simple split:
		}
		// usar split real:
		s := []byte(hashedPassword)
		idx := -1
		for i := range s {
			if s[i] == '$' {
				idx = i
				break
			}
		}
		if idx == -1 {
			return errors.New("invalid stored hash format")
		}
		parts[0] = string(s[:idx])
		parts[1] = string(s[idx+1:])
	}

	saltHex := parts[0]
	hashHex := parts[1]

	salt, err := hex.DecodeString(saltHex)
	if err != nil {
		return errors.New("invalid salt encoding")
	}

	expectedHash, err := hex.DecodeString(hashHex)
	if err != nil {
		return errors.New("invalid hash encoding")
	}

	h := sha256.New()
	h.Write(salt)
	h.Write([]byte(password))
	computed := h.Sum(nil)

	if subtle.ConstantTimeCompare(computed, expectedHash) != 1 {
		return errors.New("invalid credentials")
	}
	return nil
}
