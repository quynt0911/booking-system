package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"booking-service/internal/repository"
)

type ConflictCheckerInterface interface {
	CheckBookingConflict(expertID, userID uint, startTime, endTime time.Time) (bool, error)
	CheckExpertAvailability(expertID uint, startTime, endTime time.Time) (bool, error)
	CheckUserConflict(userID uint, startTime, endTime time.Time) (bool, error)
	LockTimeSlot(expertID uint, startTime, endTime time.Time) (string, error)
	ReleaseLock(lockKey string) error
}

type ConflictChecker struct {
	bookingRepo repository.BookingRepositoryInterface
	redisClient *redis.Client
}

func NewConflictChecker(
	bookingRepo repository.BookingRepositoryInterface,
	redisClient *redis.Client,
) ConflictCheckerInterface {
	return &ConflictChecker{
		bookingRepo: bookingRepo,
		redisClient: redisClient,
	}
}

func (c *ConflictChecker) CheckBookingConflict(expertID, userID uint, startTime, endTime time.Time) (bool, error) {
	// Check expert availability
	expertHasConflict, err := c.CheckExpertAvailability(expertID, startTime, endTime)
	if err != nil {
		return false, fmt.Errorf("failed to check expert availability: %v", err)
	}
	if !expertHasConflict {
		return true, nil // Expert is not available
	}

	// Check user conflict
	userHasConflict, err := c.CheckUserConflict(userID, startTime, endTime)
	if err != nil {
		return false, fmt.Errorf("failed to check user conflict: %v", err)
	}
	if userHasConflict {
		return true, nil // User has conflicting booking
	}

	return false, nil // No conflicts
}

func (c *ConflictChecker) CheckExpertAvailability(expertID uint, startTime, endTime time.Time) (bool, error) {
	// First check Redis cache for quick lookup
	cacheKey := fmt.Sprintf("expert_busy:%d:%s:%s", 
		expertID, 
		startTime.Format("2006-01-02T15:04:05"), 
		endTime.Format("2006-01-02T15:04:05"))
	
	ctx := context.Background()
	cached, err := c.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var isBusy bool
		if json.Unmarshal([]byte(cached), &isBusy) == nil {
			return !isBusy, nil // Return availability (opposite of busy)
		}
	}

	// Check database for existing bookings
	hasConflict, err := c.bookingRepo.HasExpertConflict(expertID, startTime, endTime)
	if err != nil {
		return false, err
	}

	// Cache the result for 5 minutes
	isBusy := hasConflict
	busyData, _ := json.Marshal(isBusy)
	c.redisClient.Set(ctx, cacheKey, busyData, 5*time.Minute)

	return !hasConflict, nil // Return availability
}

func (c *ConflictChecker) CheckUserConflict(userID uint, startTime, endTime time.Time) (bool, error) {
	// Check Redis cache first
	cacheKey := fmt.Sprintf("user_busy:%d:%s:%s", 
		userID, 
		startTime.Format("2006-01-02T15:04:05"), 
		endTime.Format("2006-01-02T15:04:05"))
	
	ctx := context.Background()
	cached, err := c.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var hasBusy bool
		if json.Unmarshal([]byte(cached), &hasBusy) == nil {
			return hasBusy, nil
		}
	}

	// Check database for user conflicts
	hasConflict, err := c.bookingRepo.HasUserConflict(userID, startTime, endTime)
	if err != nil {
		return false, err
	}

	// Cache the result for 5 minutes
	conflictData, _ := json.Marshal(hasConflict)
	c.redisClient.Set(ctx, cacheKey, conflictData, 5*time.Minute)

	return hasConflict, nil
}

func (c *ConflictChecker) LockTimeSlot(expertID uint, startTime, endTime time.Time) (string, error) {
	lockKey := fmt.Sprintf("booking_lock:%d:%s:%s", 
		expertID, 
		startTime.Format("2006-01-02T15:04:05"), 
		endTime.Format("2006-01-02T15:04:05"))
	
	ctx := context.Background()
	
	// Try to acquire lock with expiration (5 minutes)
	success, err := c.redisClient.SetNX(ctx, lockKey, "locked", 5*time.Minute).Result()
	if err != nil {
		return "", fmt.Errorf("failed to acquire lock: %v", err)
	}
	
	if !success {
		return "", fmt.Errorf("time slot is currently being booked by another user")
	}
	
	return lockKey, nil
}

func (c *ConflictChecker) ReleaseLock(lockKey string) error {
	ctx := context.Background()
	return c.redisClient.Del(ctx, lockKey).Err()
}

// Additional helper methods

// CheckMultipleTimeSlots checks availability for multiple time slots at once
func (c *ConflictChecker) CheckMultipleTimeSlots(expertID uint, timeSlots []TimeSlot) (map[int]bool, error) {
	availability := make(map[int]bool)
	
	for i, slot := range timeSlots {
		available, err := c.CheckExpertAvailability(expertID, slot.StartTime, slot.EndTime)
		if err != nil {
			return nil, fmt.Errorf("failed to check slot %d: %v", i, err)
		}
		availability[i] = available
	}
	
	return availability, nil
}

// TimeSlot represents a time period
type TimeSlot struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// GetExpertBusySlots returns all busy time slots for an expert on a given date
func (c *ConflictChecker) GetExpertBusySlots(expertID uint, date time.Time) ([]TimeSlot, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("expert_busy_slots:%d:%s", expertID, date.Format("2006-01-02"))
	ctx := context.Background()
	
	cached, err := c.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var slots []TimeSlot
		if json.Unmarshal([]byte(cached), &slots) == nil {
			return slots, nil
		}
	}

	// Get from database
	bookings, err := c.bookingRepo.GetExpertBookingsByDate(expertID, date)
	if err != nil {
		return nil, err
	}

	var busySlots []TimeSlot
	for _, booking := range bookings {
		busySlots = append(busySlots, TimeSlot{
			StartTime: booking.StartTime,
			EndTime:   booking.EndTime,
		})
	}

	// Cache for 10 minutes
	slotsData, _ := json.Marshal(busySlots)
	c.redisClient.Set(ctx, cacheKey, slotsData, 10*time.Minute)

	return busySlots, nil
}

// IsTimeSlotOverlapping checks if two time slots overlap
func (c *ConflictChecker) IsTimeSlotOverlapping(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && start2.Before(end1)
}

// ValidateTimeSlot validates if a time slot is valid (end time after start time, minimum duration, etc.)
func (c *ConflictChecker) ValidateTimeSlot(startTime, endTime time.Time) error {
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		return fmt.Errorf("end time must be after start time")
	}

	duration := endTime.Sub(startTime)
	if duration < 30*time.Minute {
		return fmt.Errorf("booking duration must be at least 30 minutes")
	}

	if duration > 4*time.Hour {
		return fmt.Errorf("booking duration cannot exceed 4 hours")
	}

	// Check if booking is too far in the future (max 6 months)
	if startTime.After(time.Now().AddDate(0, 6, 0)) {
		return fmt.Errorf("cannot book more than 6 months in advance")
	}

	// Check if booking is in the past
	if startTime.Before(time.Now()) {
		return fmt.Errorf("cannot book in the past")
	}

	return nil
}

// ClearExpertCache clears all cached data for an expert
func (c *ConflictChecker) ClearExpertCache(expertID uint) error {
	ctx := context.Background()
	pattern := fmt.Sprintf("expert_*:%d:*", expertID)
	
	keys, err := c.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.redisClient.Del(ctx, keys...).Err()
	}

	return nil
}