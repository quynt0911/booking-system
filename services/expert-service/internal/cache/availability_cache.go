package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type AvailabilityCache interface {
	SetAvailability(expertID int, date string, isAvailable bool) error
	GetAvailability(expertID int, date string) (bool, bool, error) // value, exists, error
	InvalidateExpert(expertID int) error
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

func (c *availabilityCache) SetAvailability(expertID int, date string, isAvailable bool) error {
	key := c.getKey(expertID, date)
	value, _ := json.Marshal(isAvailable)

	return c.client.Set(context.Background(), key, value, c.ttl).Err()
}

func (c *availabilityCache) GetAvailability(expertID int, date string) (bool, bool, error) {
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

func (c *availabilityCache) InvalidateExpert(expertID int) error {
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
