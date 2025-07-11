package listener

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type DBListener struct {
	listener     *pq.Listener
	redisClient  *redis.Client
	dbConnString string
}

func NewDBListener(dbConnString string, redisClient *redis.Client) *DBListener {
	return &DBListener{
		dbConnString: dbConnString,
		redisClient:  redisClient,
	}
}

func (l *DBListener) Start() error {
	listener := pq.NewListener(l.dbConnString, 10*time.Second, time.Minute, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Printf("Error in database listener: %v\n", err)
		}
	})

	err := listener.Listen("booking_changes")
	if err != nil {
		return fmt.Errorf("error starting listener: %v", err)
	}

	l.listener = listener

	// Start processing notifications in a goroutine
	go l.processNotifications()

	return nil
}

func (l *DBListener) processNotifications() {
	for {
		notification, ok := <-l.listener.Notify
		if !ok {
			return
		}

		if notification == nil {
			continue
		}

		// Parse notification payload
		// Expected format: "DELETE:booking_id:expert_id:user_id"
		var operation, bookingID, expertID, userID string
		_, err := fmt.Sscanf(notification.Extra, "%s:%s:%s:%s", &operation, &bookingID, &expertID, &userID)
		if err != nil {
			log.Printf("Error parsing notification: %v\n", err)
			continue
		}

		if operation == "DELETE" {
			l.handleBookingDeletion(bookingID, expertID, userID)
		}
	}
}

func (l *DBListener) handleBookingDeletion(bookingID, expertID, userID string) {
	ctx := context.Background()

	// 1. Clear booking cache
	bookingKey := fmt.Sprintf("booking:%s", bookingID)
	l.redisClient.Del(ctx, bookingKey)

	// 2. Clear expert cache
	expertCachePattern := fmt.Sprintf("expert_*:%s:*", expertID)
	if keys, err := l.redisClient.Keys(ctx, expertCachePattern).Result(); err == nil && len(keys) > 0 {
		l.redisClient.Del(ctx, keys...)
	}

	// 3. Clear expert busy time cache
	busyTimePattern := fmt.Sprintf("expert_busy:%s:*", expertID)
	if keys, err := l.redisClient.Keys(ctx, busyTimePattern).Result(); err == nil && len(keys) > 0 {
		l.redisClient.Del(ctx, keys...)
	}

	// 4. Clear availability cache
	availabilityPattern := fmt.Sprintf("availability:%s:*", expertID)
	if keys, err := l.redisClient.Keys(ctx, availabilityPattern).Result(); err == nil && len(keys) > 0 {
		l.redisClient.Del(ctx, keys...)
	}

	// 5. Clear user cache
	userCachePattern := fmt.Sprintf("user_*:%s:*", userID)
	if keys, err := l.redisClient.Keys(ctx, userCachePattern).Result(); err == nil && len(keys) > 0 {
		l.redisClient.Del(ctx, keys...)
	}

	log.Printf("Cleared cache for deleted booking %s of expert %s and user %s\n", bookingID, expertID, userID)
}

func (l *DBListener) Stop() {
	if l.listener != nil {
		l.listener.Close()
	}
}
