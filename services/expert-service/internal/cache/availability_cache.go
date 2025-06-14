package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type AvailabilityCache interface {
	SetAvailability(key string, value []byte) error
	GetAvailability(key string) ([]byte, error)
	InvalidateExpert(expertID string) error
}

type availabilityCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewAvailabilityCache(client *redis.Client, ttl time.Duration) AvailabilityCache {
	return &availabilityCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *availabilityCache) SetAvailability(key string, value []byte) error {
	return c.client.Set(context.Background(), key, value, c.ttl).Err()
}

func (c *availabilityCache) GetAvailability(key string) ([]byte, error) {
	val, err := c.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, nil // Key not found
	}
	if err != nil {
		return nil, err
	}
	return []byte(val), nil
}

func (c *availabilityCache) InvalidateExpert(expertID string) error {
	// Invalidate availability check cache keys (e.g., "expertID:date")
	pattern1 := fmt.Sprintf("%s:*", expertID)
	keys1, err := c.client.Keys(context.Background(), pattern1).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys for pattern %s: %w", pattern1, err)
	}
	if len(keys1) > 0 {
		if err := c.client.Del(context.Background(), keys1...).Err(); err != nil {
			return fmt.Errorf("failed to delete keys for pattern %s: %w", pattern1, err)
		}
	}

	// Invalidate availability slot keys (e.g., "availability:expertID:date")
	pattern2 := fmt.Sprintf("availability:%s:*", expertID)
	keys2, err := c.client.Keys(context.Background(), pattern2).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys for pattern %s: %w", pattern2, err)
	}
	if len(keys2) > 0 {
		if err := c.client.Del(context.Background(), keys2...).Err(); err != nil {
			return fmt.Errorf("failed to delete keys for pattern %s: %w", pattern2, err)
		}
	}

	return nil
}

func (c *availabilityCache) SetAvailabilityRedis(expertID int, date string, isAvailable bool) error {
	key := c.getKey(expertID, date)
	value, _ := json.Marshal(isAvailable)

	return c.client.Set(context.Background(), key, value, c.ttl).Err()
}

func (c *availabilityCache) GetAvailabilityRedis(expertID int, date string) (bool, bool, error) {
	key := c.getKey(expertID, date)

	val, err := c.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return false, false, nil // Not found
	}
	if err != nil {
		return false, false, err
	}

	var isAvailable bool
	err = json.Unmarshal([]byte(val), &isAvailable)
	return isAvailable, true, err
}

func (c *availabilityCache) InvalidateExpertRedis(expertID int) error {
	pattern := fmt.Sprintf("availability:expert:%d:*", expertID)

	keys, err := c.client.Keys(context.Background(), pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.client.Del(context.Background(), keys...).Err()
	}

	return nil
}

func (c *availabilityCache) getKey(expertID int, date string) string {
	return fmt.Sprintf("availability:expert:%d:date:%s", expertID, date)
}
