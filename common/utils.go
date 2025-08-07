package common

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
)

func GetSessionSecret() string {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "claude-code-relay-default-secret"
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

func GetSalt() string {
	salt := os.Getenv("SALT")
	if salt == "" {
		salt = "claude-code-relay-default-salt"
	}
	return salt
}

func HashPassword(password string) string {
	salt := GetSalt()
	hash := md5.Sum([]byte(password + salt))
	return hex.EncodeToString(hash[:])
}

func VerifyPassword(password, hashedPassword string) bool {
	return HashPassword(password) == hashedPassword
}

func CloseIO(c io.Closer) {
	err := c.Close()
	if nil != err {
		log.Println(err)
	}
}
