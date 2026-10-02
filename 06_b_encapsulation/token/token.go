package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// Manager defines the public behavior for generating and validating signed tokens.
//
// Key Encapsulation Lesson:
// Here, encapsulation reaches its highest level:
// 1. The interface 'Manager' is EXPORTED.
// 2. The concrete struct implementing it ('hmacManager') is UNEXPORTED!
// 3. The constructor returns the interface type 'Manager'.
// External callers cannot even name the underlying struct, let alone inspect
// its secret fields or internal hash helpers.
type Manager interface {
	Generate(subject string, ttl time.Duration) string
	Validate(tokenStr string) (string, bool)
}

// hmacManager is an UNEXPORTED concrete struct.
type hmacManager struct {
	secretKey []byte
}

// NewManager creates a token Manager while hiding the concrete implementation.
func NewManager(secret string) Manager {
	return &hmacManager{
		secretKey: []byte(secret),
	}
}

// Generate creates a signed token format: "<subject>.<expiry_unix_milli>.<signature>"
func (m *hmacManager) Generate(subject string, ttl time.Duration) string {
	expiry := time.Now().Add(ttl).UnixMilli()
	payload := fmt.Sprintf("%s.%d", subject, expiry)
	signature := m.sign(payload)
	return fmt.Sprintf("%s.%s", payload, signature)
}

// Validate checks the signature and expiration, returning the subject if valid.
func (m *hmacManager) Validate(tokenStr string) (string, bool) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return "", false
	}

	subject := parts[0]
	expiryStr := parts[1]
	signature := parts[2]

	payload := fmt.Sprintf("%s.%s", subject, expiryStr)
	expectedSig := m.sign(payload)

	// Constant time compare to avoid timing attacks
	if !hmac.Equal([]byte(signature), []byte(expectedSig)) {
		return "", false
	}

	var expiry int64
	_, err := fmt.Sscanf(expiryStr, "%d", &expiry)
	if err != nil || time.Now().UnixMilli() > expiry {
		return "", false // Expired or malformed timestamp
	}

	return subject, true
}

// sign is an UNEXPORTED helper method.
func (m *hmacManager) sign(data string) string {
	mac := hmac.New(sha256.New, m.secretKey)
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}
