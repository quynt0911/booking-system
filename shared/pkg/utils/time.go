package utils

import (
	"fmt"
	"time"
)

const (
	DateFormat     = "2006-01-02"
	TimeFormat     = "15:04:05"
	DateTimeFormat = "2006-01-02 15:04:05"
	ISOFormat      = "2006-01-02T15:04:05Z07:00"
)

// TimeHelper provides utility functions for time operations
type TimeHelper struct {
	location *time.Location
}

// NewTimeHelper creates a new TimeHelper with Vietnam timezone
func NewTimeHelper() *TimeHelper {
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	return &TimeHelper{
		location: loc,
	}
}

// Now returns current time in Vietnam timezone
func (th *TimeHelper) Now() time.Time {
	return time.Now().In(th.location)
}

// ParseDate parses date string to time.Time
func (th *TimeHelper) ParseDate(dateStr string) (time.Time, error) {
	return time.ParseInLocation(DateFormat, dateStr, th.location)
}

// ParseDateTime parses datetime string to time.Time
func (th *TimeHelper) ParseDateTime(datetimeStr string) (time.Time, error) {
	return time.ParseInLocation(DateTimeFormat, datetimeStr, th.location)
}

// ParseTime parses time string to time.Time for today
func (th *TimeHelper) ParseTime(timeStr string) (time.Time, error) {
	now := th.Now()
	todayStr := now.Format(DateFormat)
	datetimeStr := fmt.Sprintf("%s %s", todayStr, timeStr)
	return time.ParseInLocation(DateTimeFormat, datetimeStr, th.location)
}

// FormatDate formats time to date string
func (th *TimeHelper) FormatDate(t time.Time) string {
	return t.In(th.location).Format(DateFormat)
}

// FormatTime formats time to time string
func (th *TimeHelper) FormatTime(t time.Time) string {
	return t.In(th.location).Format(TimeFormat)
}

// FormatDateTime formats time to datetime string
func (th *TimeHelper) FormatDateTime(t time.Time) string {
	return t.In(th.location).Format(DateTimeFormat)
}

// IsWorkingHour checks if time is within working hours (9:00 - 17:00)
func (th *TimeHelper) IsWorkingHour(t time.Time) bool {
	hour := t.In(th.location).Hour()
	return hour >= 9 && hour < 17
}

// IsWorkingDay checks if time is on working day (Monday to Friday)
func (th *TimeHelper) IsWorkingDay(t time.Time) bool {
	weekday := t.In(th.location).Weekday()
	return weekday >= time.Monday && weekday <= time.Friday
}

// GetNextWorkingDay returns next working day
func (th *TimeHelper) GetNextWorkingDay(from time.Time) time.Time {
	next := from.AddDate(0, 0, 1)
	for !th.IsWorkingDay(next) {
		next = next.AddDate(0, 0, 1)
	}
	return next
}

// GetWeekStart returns start of week (Monday)
func (th *TimeHelper) GetWeekStart(t time.Time) time.Time {
	weekday := t.In(th.location).Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	days := int(weekday) - 1
	return t.AddDate(0, 0, -days).Truncate(24 * time.Hour)
}

// GetWeekEnd returns end of week (Sunday)
func (th *TimeHelper) GetWeekEnd(t time.Time) time.Time {
	return th.GetWeekStart(t).AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
}

// GetMonthStart returns start of month
func (th *TimeHelper) GetMonthStart(t time.Time) time.Time {
	year, month, _ := t.In(th.location).Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, th.location)
}

// GetMonthEnd returns end of month
func (th *TimeHelper) GetMonthEnd(t time.Time) time.Time {
	return th.GetMonthStart(t).AddDate(0, 1, 0).Add(-time.Second)
}

// IsPast checks if time is in the past
func (th *TimeHelper) IsPast(t time.Time) bool {
	return t.Before(th.Now())
}

// IsFuture checks if time is in the future
func (th *TimeHelper) IsFuture(t time.Time) bool {
	return t.After(th.Now())
}

// IsToday checks if time is today
func (th *TimeHelper) IsToday(t time.Time) bool {
	now := th.Now()
	return th.FormatDate(t) == th.FormatDate(now)
}

// AddBusinessDays adds business days to time
func (th *TimeHelper) AddBusinessDays(t time.Time, days int) time.Time {
	result := t
	for i := 0; i < days; i++ {
		result = th.GetNextWorkingDay(result)
	}
	return result
}

// GetTimeSlots returns available time slots for a day
func (th *TimeHelper) GetTimeSlots(date time.Time, startHour, endHour, durationMinutes int) []TimeSlot {
	var slots []TimeSlot

	start := time.Date(date.Year(), date.Month(), date.Day(), startHour, 0, 0, 0, th.location)
	end := time.Date(date.Year(), date.Month(), date.Day(), endHour, 0, 0, 0, th.location)
	duration := time.Duration(durationMinutes) * time.Minute

	for current := start; current.Before(end); current = current.Add(duration) {
		slotEnd := current.Add(duration)
		if slotEnd.After(end) {
			break
		}

		slots = append(slots, TimeSlot{
			Start: current,
			End:   slotEnd,
		})
	}

	return slots
}

// TimeSlot represents a time slot
type TimeSlot struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// String returns string representation of time slot
func (ts TimeSlot) String() string {
	th := NewTimeHelper()
	return fmt.Sprintf("%s - %s",
		th.FormatTime(ts.Start),
		th.FormatTime(ts.End))
}

// Duration returns duration of time slot
func (ts TimeSlot) Duration() time.Duration {
	return ts.End.Sub(ts.Start)
}

// Contains checks if time is within the slot
func (ts TimeSlot) Contains(t time.Time) bool {
	return (t.Equal(ts.Start) || t.After(ts.Start)) && t.Before(ts.End)
}

// Overlaps checks if two time slots overlap
func (ts TimeSlot) Overlaps(other TimeSlot) bool {
	return ts.Start.Before(other.End) && ts.End.After(other.Start)
}
