package middleware

import (
	"claude-code-relay/common"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if common.RDB == nil {
			c.Next()
			return
		}

		// 使用IP作为限流键
		key := fmt.Sprintf("rate_limit:%s", c.ClientIP())
		ctx := context.Background()

		// 获取当前计数
		count, err := common.RDB.Get(ctx, key).Int()
		if err != nil && err.Error() != "redis: nil" {
			c.Next()
			return
		}

		if count >= maxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
				"code":  42901,
			})
			c.Abort()
			return
		}

		// 增加计数
		pipe := common.RDB.Pipeline()
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, window)
		_, err = pipe.Exec(ctx)
		if err != nil {
			common.SysError("Rate limit pipeline error: " + err.Error())
		}

		// 设置响应头
		c.Header("X-RateLimit-Limit", strconv.Itoa(maxRequests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(maxRequests-count-1))

		c.Next()
	}
}
