package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mehtahetul/proximate/internal/db"
)

func RateLimit(endpoint string, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			// For public endpoints (login, register), fall back to IP address
			userID = c.ClientIP()
		}

		key := fmt.Sprintf("ratelimit:%s:%s", endpoint, userID)
		ctx := context.Background()

		// Atomically increment the counter
		count, err := db.RedisClient.Incr(ctx, key).Result()
		if err != nil {
			// Redis is down — fail open (let the request through)
			c.Next()
			return
		}

		// First request in this window — set the expiry
		if count == 1 {
			db.RedisClient.Expire(ctx, key, window)
		}

		if count > int64(limit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests, slow down",
			})
			return
		}

		c.Next()
	}
}