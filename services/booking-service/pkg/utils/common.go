package utils

import (
	"time"

	"github.com/google/uuid"
)

// ParseTime parses a time string in the format "2006-01-02"
func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse("2006-01-02", timeStr)
}

// FormatTime formats a time to string in the format "2006-01-02"
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02")
}

// ParseUUID parses a string to UUID
func ParseUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

// GetStartOfDay returns the start of the day for a given time
func GetStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetEndOfDay returns the end of the day for a given time
func GetEndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// IsWithinTimeRange checks if a time is within a given range
func IsWithinTimeRange(t, start, end time.Time) bool {
	return t.After(start) && t.Before(end)
}
