package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

const nonceSize = 12

// Encrypt encrypts plaintext with AES-256-GCM using the provided 32-byte key.
// Returns a base64-encoded string of nonce+ciphertext.
func Encrypt(plaintext []byte, key []byte) (string, error) {
	if len(key) != 32 {
		return "", fmt.Errorf("Encrypt: key must be exactly 32 bytes, got %d", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("Encrypt: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("Encrypt: %w", err)
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("Encrypt: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	combined := make([]byte, nonceSize+len(ciphertext))
	copy(combined, nonce)
	copy(combined[nonceSize:], ciphertext)

	return base64.StdEncoding.EncodeToString(combined), nil
}

// Decrypt decodes a base64-encoded AES-256-GCM ciphertext (nonce prepended)
// and returns the original plaintext.
func Decrypt(encoded string, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("Decrypt: key must be exactly 32 bytes, got %d", len(key))
	}

	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("Decrypt: base64 decode: %w", err)
	}

	if len(data) < nonceSize {
		return nil, fmt.Errorf("Decrypt: ciphertext too short")
	}

	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("Decrypt: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("Decrypt: %w", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("Decrypt: %w", err)
	}

	return plaintext, nil
}
