package common

import (
	"context"
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

// GenerateRandomInstanceID 生成随机示例ID
func GenerateRandomInstanceID() string {
	bytes := make([]byte, 31) // 31字节 * 2 = 62个十六进制字符，取前61位
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)[:61] // 取前61位
}

// GetInstanceID 获取实例ID，如果不存在则生成61位随机字符串并存储到Redis
// 仅在分组示例ID不存在时调用, 即全局共享组专属ID
func GetInstanceID() string {
	const instanceKey = "system:instance_id"
	ctx := context.Background()

	// 从Redis获取
	if RDB != nil {
		id, err := RDB.Get(ctx, instanceKey).Result()
		if err == nil && id != "" {
			return id
		}
	}

	newID := GenerateRandomInstanceID()
	if newID == "" {
		SysError("Failed to generate instance ID")
		return ""
	}

	// 存储到Redis（永久存储）
	if RDB != nil {
		err := RDB.Set(ctx, instanceKey, newID, 0).Err()
		if err != nil {
			SysError("Failed to store instance ID to Redis: " + err.Error())
		}
	}

	SysLog("Generated new instance ID: " + newID)
	return newID
}
