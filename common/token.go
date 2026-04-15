package common

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"
)

const (
	TokenStatusEnabled  = 1
	TokenStatusDisabled = 2
	TokenStatusExpired  = 3
	TokenStatusExhausted = 4

	KeyPrefix = "sk-"
	KeyLength = 48
)

func GenerateKey() string {
	rawKey := make([]byte, (KeyLength-len(KeyPrefix))/2)
	_, err := rand.Read(rawKey)
	if err != nil {
		return ""
	}
	return KeyPrefix + hex.EncodeToString(rawKey)
}

func GetTimestamp() int64 {
	return time.Now().Unix()
}

func IsValidKey(key string) bool {
	return strings.HasPrefix(key, KeyPrefix) && len(key) == KeyLength
}

func MaskKey(key string) string {
	if len(key) <= len(KeyPrefix)+8 {
		return key
	}
	return key[:len(KeyPrefix)+4] + strings.Repeat("*", len(key)-len(KeyPrefix)-8) + key[len(key)-4:]
}
