package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const nearbyTTL = 30 * time.Second

func NearbyKey(userID string, lat, lng, radius float64) string {
	return fmt.Sprintf("nearby:%s:%.4f:%.4f:%.0f", userID, lat, lng, radius) //lat and lng rounded to 4 deciamal points i.e. approx 11 meters of precision
}

func GetNearby(ctx context.Context, client *redis.Client, key string) (string, bool) {
	val, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		// Key doesn't exist — cache miss
		return "", false
	}
	if err != nil {
		// Redis error — treat as cache miss, don't crash
		return "", false
	}
	return val, true
}

func SetNearby(ctx context.Context, client *redis.Client, key, value string) {
	// Fire and forget — if this fails, it's not fatal
	client.Set(ctx, key, value, nearbyTTL)
}
