package common

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"time"

	"github.com/google/uuid"
)

func GetSessionSecret() string {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "claude-scheduler-default-secret"
	}
	return secret
}

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}